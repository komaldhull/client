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
	"fmt"

	"github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"github.com/spf13/cobra"

	v1alpha12 "github.com/knative/client/pkg/eventing/v1alpha1"
	servingv1alpha1 "github.com/knative/client/pkg/serving/v1alpha1"

	"github.com/knative/client/pkg/kn/commands"
)

// NewBrokerListCommand represents 'kn broker list' command
func NewBrokerListCommand(p *commands.KnParams) *cobra.Command {
	brokerListFlags := NewBrokerListFlags()

	brokerListCommand := &cobra.Command{
		Use:   "list [name]",
		Short: "List available brokers",
		Example: `
  # List all brokers
  kn broker list
  
  # List broker 'default'
  kn broker list default`,
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}
			client, err := p.NewEventClient(namespace)
			if err != nil {
				return err
			}
			brokerList, err := getBrokerInfo(args, client)
			if err != nil {
				return err
			}
			if len(brokerList.Items) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No resources found.\n")
				return nil
			}
			printer, err := brokerListFlags.ToPrinter()
			if err != nil {
				return err
			}

			err = printer.PrintObj(brokerList, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return nil
		},
	}
	commands.AddNamespaceFlags(brokerListCommand.Flags(), true)
	brokerListFlags.AddFlags(brokerListCommand)
	return brokerListCommand
}

func getBrokerInfo(args []string, client v1alpha12.KnEventClient) (*v1alpha1.BrokerList, error) {
	var (
		brokerList *v1alpha1.BrokerList
		err        error
	)
	switch len(args) {
	case 0:
		brokerList, err = client.ListBrokers()
	case 1:
		brokerList, err = client.ListBrokers(servingv1alpha1.WithName(args[0]))
	default:
		return nil, fmt.Errorf("'kn broker list' accepts maximum 1 argument")
	}
	return brokerList, err
}
