// Copyright © 2019 The Knative Authors
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

package importer

import (
	"fmt"

	"github.com/knative/eventing/pkg/apis/sources/v1alpha1"
	"github.com/spf13/cobra"

	v1alpha12 "github.com/knative/client/pkg/eventing/sourcesv1alpha1"
	servingv1alpha1 "github.com/knative/client/pkg/serving/v1alpha1"

	"github.com/knative/client/pkg/kn/commands"
)

func NewImporterListCommand(p *commands.KnParams) *cobra.Command {
	importerListFlags := NewImporterListFlags()

	importerListCommand := &cobra.Command{
		Use:   "list [name]",
		Short: "List available importers",
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}
			client, err := p.NewSourceClient(namespace)
			if err != nil {
				return err
			}
			importerList, err := getImporterInfo(args, client, cmd)
			if err != nil {
				return err
			}
			if len(importerList.Items) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No resources found.\n")
				return nil
			}
			printer, err := importerListFlags.ToPrinter()
			if err != nil {
				return err
			}

			err = printer.PrintObj(importerList, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return nil
		},
	}
	commands.AddNamespaceFlags(importerListCommand.Flags(), true)
	importerListFlags.AddFlags(importerListCommand)
	return importerListCommand
}

func getImporterInfo(args []string, client v1alpha12.KnSourceClient, cmd *cobra.Command) (*v1alpha1.ContainerSourceList, error) {
	var (
		srcList *v1alpha1.ContainerSourceList
		err     error
	)
	switch len(args) {
	case 0:
		srcList, err = client.ListContainerSources()
	case 1:
		srcList, err = client.ListContainerSources(servingv1alpha1.WithName(args[0]))
	default:
		return nil, fmt.Errorf("'kn importer list' accepts maximum 1 argument")
	}
	return srcList, err
}
