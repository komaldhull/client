// Copyright © 2018 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/knative/client/pkg/kn/commands"
	"github.com/knative/client/pkg/kn/commands/plugin"
	"github.com/knative/client/pkg/kn/commands/revision"
	"github.com/knative/client/pkg/kn/commands/route"
	"github.com/knative/client/pkg/kn/commands/service"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// NewDefaultKnCommand creates the default `kn` command with a default plugin handler
func NewDefaultKnCommand() *cobra.Command {
	rootCmd := NewKnCommand()

	// Needed since otherwise --plugins-dir and --lookup-plugins-in-path
	// will not be accounted for since the plugin is not a Cobra command
	// and will not be parsed
	pluginsDir, lookupPluginsInPath, err := extractKnPluginFlags(os.Args)
	if err != nil {
		panic("Invalid plugin flag value")
	}

	pluginHandler := plugin.NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes,
		pluginsDir, lookupPluginsInPath)

	return NewDefaultKnCommandWithArgs(rootCmd, pluginHandler,
		os.Args, os.Stdin,
		os.Stdout, os.Stderr)
}

// NewDefaultKnCommandWithArgs creates the `kn` command with arguments
func NewDefaultKnCommandWithArgs(rootCmd *cobra.Command,
	pluginHandler plugin.PluginHandler,
	args []string,
	in io.Reader,
	out,
	errOut io.Writer) *cobra.Command {
	if pluginHandler == nil {
		return rootCmd
	}

	if len(args) > 1 {
		cmdPathPieces := args[1:]
		cmdPathPieces = removeKnPluginFlags(cmdPathPieces) // Plugin does not need these flags

		// only look for suitable extension executables if
		// the specified command does not already exist
		if _, _, err := rootCmd.Find(cmdPathPieces); err != nil {
			err := plugin.HandlePluginCommand(pluginHandler, cmdPathPieces)
			if err != nil {
				fmt.Fprintf(errOut, "%v\n", err)
				os.Exit(1)
			}
		}
	}

	return rootCmd
}

// NewKnCommand creates the rootCmd which is the base command when called without any subcommands
func NewKnCommand(params ...commands.KnParams) *cobra.Command {
	var p *commands.KnParams
	if len(params) == 0 {
		p = &commands.KnParams{}
	} else if len(params) == 1 {
		p = &params[0]
	} else {
		panic("Too many params objects to NewKnCommand")
	}
	p.Initialize()

	rootCmd := &cobra.Command{
		Use:   "kn",
		Short: "Knative client",
		Long: `Manage your Knative building blocks:

* Serving: Manage your services and release new software to them.
* Eventing: Manage event subscriptions and channels. Connect up event sources.`,

		// Disable docs header
		DisableAutoGenTag: true,

		// Affects children as well
		SilenceUsage: true,

		// Prevents Cobra from dealing with errors as we deal with them in main.go
		SilenceErrors: true,
	}
	if p.Output != nil {
		rootCmd.SetOutput(p.Output)
	}

	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&commands.CfgFile, "config", "", "kn config file (default is $HOME/.kn/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&p.KubeCfgPath, "kubeconfig", "", "kubectl config file (default is $HOME/.kube/config)")

	plugin.AddPluginFlags(rootCmd)
	plugin.BindPluginsFlagToViper(rootCmd)

	// root child commands
	rootCmd.AddCommand(service.NewServiceCommand(p))
	rootCmd.AddCommand(revision.NewRevisionCommand(p))
	rootCmd.AddCommand(plugin.NewPluginCommand(p))
	rootCmd.AddCommand(route.NewRouteCommand(p))
	rootCmd.AddCommand(commands.NewCompletionCommand(p))
	rootCmd.AddCommand(commands.NewVersionCommand(p))

	// Deal with empty and unknown sub command groups
	EmptyAndUnknownSubCommands(rootCmd)

	// For glog parse error.
	flag.CommandLine.Parse([]string{})

	return rootCmd
}

// InitializeConfig initializes the kubeconfig used by all commands
func InitializeConfig() {
	cobra.OnInitialize(initConfig)
}

// EmptyAndUnknownSubCommands adds a RunE to all commands that are groups to
// deal with errors when called with empty or unknown sub command
func EmptyAndUnknownSubCommands(cmd *cobra.Command) {
	for _, childCmd := range cmd.Commands() {
		if childCmd.HasSubCommands() && childCmd.RunE == nil {
			childCmd.RunE = func(aCmd *cobra.Command, args []string) error {
				aCmd.Help()
				fmt.Println()

				if len(args) == 0 {
					return errors.New(fmt.Sprintf("please provide a valid sub-command for \"kn %s\"", aCmd.Name()))
				} else {
					return errors.New(fmt.Sprintf("unknown sub-command \"%s\" for \"kn %s\"", args[0], aCmd.Name()))
				}
			}
		}

		// recurse to deal with child commands that are themselves command groups
		EmptyAndUnknownSubCommands(childCmd)
	}
}

// Private

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if commands.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(commands.CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kn" (without extension)
		viper.AddConfigPath(path.Join(home, ".kn"))
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Fprintln(os.Stderr, "Using kn config file:", viper.ConfigFileUsed())
	}
}

func extractKnPluginFlags(args []string) (string, bool, error) {
	pluginsDir := "~/.kn/plugins"
	lookupPluginsInPath := false
	for _, arg := range args {
		if strings.Contains(arg, "--plugins-dir") {
			values := strings.Split(arg, "=")
			if len(values) < 1 {
				return "", false, errors.New("Invalid --plugins-dir flag value")
			}
			pluginsDir = values[1]
		}

		if strings.Contains(arg, "--lookup-plugins-in-path") {
			values := strings.Split(arg, "=")
			if len(values) < 1 {
				return "", false, errors.New("Invalid --lookup-plugins-in-path flag value")
			}

			boolValue, err := strconv.ParseBool(values[1])
			if err != nil {
				return "", false, err
			}

			lookupPluginsInPath = boolValue
		}
	}
	return pluginsDir, lookupPluginsInPath, nil
}

func removeKnPluginFlags(args []string) []string {
	var remainingArgs []string
	for _, arg := range args {
		if strings.Contains(arg, "--plugins-dir") ||
			strings.Contains(arg, "--lookup-plugins-in-path") {
			continue
		} else {
			remainingArgs = append(remainingArgs, arg)
		}
	}

	return remainingArgs
}
