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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	api_errors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/knative/client/pkg/kn/commands"
)

func NewConnectionUpdateCommand(p *commands.KnParams) *cobra.Command {
	var editFlags ConfigurationEditFlags

	connUpdateCommand := &cobra.Command{
		Use:   "update NAME",
		Short: "Update a connection.",

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("requires the connection name.")
			}

			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}

			client, err := p.NewEventClient(namespace)
			if err != nil {
				return err
			}

			var retries = 0
			for {
				trigger, err := client.GetTrigger(args[0])
				if err != nil {
					return err
				}
				trigger = trigger.DeepCopy()

				err = editFlags.Apply(trigger, cmd)
				if err != nil {
					return err
				}

				err = client.UpdateTrigger(trigger)
				if err != nil {
					// Retry to update when a resource version conflict exists
					if api_errors.IsConflict(err) && retries < MaxUpdateRetries {
						retries++
						continue
					}
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Connection '%s' updated in namespace '%s'.\n", args[0], namespace)
				return nil
			}
		},
	}

	commands.AddNamespaceFlags(connUpdateCommand.Flags(), false)
	editFlags.AddUpdateFlags(connUpdateCommand)
	return connUpdateCommand
}
