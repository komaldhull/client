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

package broker

import (
	"errors"

	"github.com/knative/client/pkg/kn/commands"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewBrokerDescribeCommand(p *commands.KnParams) *cobra.Command {
	brokerDescribePrintFlags := genericclioptions.NewPrintFlags("").WithDefaultOutput("yaml")
	brokerDescribeCmd := &cobra.Command{
		Use:   "describe NAME",
		Short: "Describe brokers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires the broker name.")
			}

			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}
			client, err := p.NewEventClient(namespace)
			if err != nil {
				return err
			}

			broker, err := client.GetBroker(args[0])
			if err != nil {
				return err
			}

			printer, err := brokerDescribePrintFlags.ToPrinter()
			if err != nil {
				return err
			}

			err = printer.PrintObj(broker, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return nil
		},
	}
	commands.AddNamespaceFlags(brokerDescribeCmd.Flags(), false)
	brokerDescribePrintFlags.AddFlags(brokerDescribeCmd)
	return brokerDescribeCmd
}
