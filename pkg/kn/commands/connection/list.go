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

package connection

import (
	"fmt"

	"github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"github.com/spf13/cobra"

	v1alpha12 "github.com/knative/client/pkg/eventing/v1alpha1"
	servingv1alpha1 "github.com/knative/client/pkg/serving/v1alpha1"

	"github.com/knative/client/pkg/kn/commands"
)

// NewConnectionListCommand represents 'kn connection list' command
func NewConnectionListCommand(p *commands.KnParams) *cobra.Command {
	connListFlags := NewConnListFlags()

	connListCommand := &cobra.Command{
		Use:   "list [name]",
		Short: "List available connections",
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}
			client, err := p.NewEventClient(namespace)
			if err != nil {
				return err
			}
			connList, err := getConnInfo(args, client, cmd)
			if err != nil {
				return err
			}
			if len(connList.Items) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No resources found.\n")
				return nil
			}
			printer, err := connListFlags.ToPrinter()
			if err != nil {
				return err
			}

			err = printer.PrintObj(connList, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return nil
		},
	}
	commands.AddNamespaceFlags(connListCommand.Flags(), true)
	connListFlags.AddFlags(connListCommand)
	return connListCommand
}

func getConnInfo(args []string, client v1alpha12.KnEventClient, cmd *cobra.Command) (*v1alpha1.TriggerList, error) {
	var (
		triggerList *v1alpha1.TriggerList
		err         error
	)
	switch len(args) {
	case 0:
		triggerList, err = client.ListTriggers()
		//NOTE: manual filtering is jank, but neither the broker or subscriber field is a supported field selector
		//solution would be to add labels subscriber=name and broker=name to triggers upon creation
		//this would need to be done on the eventing side-- "kn connection list" also needs to show manually created triggers, so kn adding a label doesn't work
		if cmd.Flags().Changed("subscriber") {
			sub := cmd.Flag("subscriber").Value.String()
			var t []eventingv1alpha1.Trigger
			for _, tri := range triggerList.Items {
				if GetSubscriberName(&tri) == sub {
					t = append(t, tri)
				}
			}
			triggerList.Items = t
		}

		if cmd.Flags().Changed("broker") {
			br := cmd.Flag("broker").Value.String()
			var t []eventingv1alpha1.Trigger
			for _, tri := range triggerList.Items {
				if GetBrokerName(&tri) == br {
					t = append(t, tri)
				}
			}
			triggerList.Items = t
		}
	case 1:
		triggerList, err = client.ListTriggers(servingv1alpha1.WithName(args[0]))
	default:
		return nil, fmt.Errorf("'kn connection list' accepts maximum 1 argument")
	}
	return triggerList, err
}
