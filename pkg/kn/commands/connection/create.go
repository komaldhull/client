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
	"io"

	"github.com/knative/client/pkg/eventing/v1alpha1"
	"github.com/knative/client/pkg/kn/commands"

	eventing_v1alpha1_api "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"github.com/spf13/cobra"

	api_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewConnectionCreateCommand(p *commands.KnParams) *cobra.Command {
	var editFlags ConfigurationEditFlags
	var waitFlags commands.WaitFlags

	connCreateCommand := &cobra.Command{
		Use:   "create NAME --image IMAGE",
		Short: "Create a connection.",

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("'connection create' requires the connection name given as single argument")
			}

			if editFlags.Sequence == "" && editFlags.Service == "" {
				return errors.New("'connection create' requires that a subscriber (either service or sequence) is specified")
			}

			if editFlags.Sequence != "" && editFlags.Service != "" {
				return errors.New("'connection create' requires that there is only one subscriber (either service or sequence)")
			}

			name := args[0]
			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}

			trigger, err := constructTrigger(cmd, editFlags, args[0], namespace)
			if err != nil {
				return err
			}

			client, err := p.NewEventClient(namespace)
			if err != nil {
				return err
			}

			triggerExists, err := triggerExists(client, name, namespace)
			if err != nil {
				return err
			}

			if triggerExists {
				if !editFlags.ForceCreate {
					return fmt.Errorf(
						"cannot create connection '%s' in namespace '%s' "+
							"because the connection already exists and no --force option was given", name, namespace)
				}
				err = replaceTrigger(client, trigger, namespace, cmd.OutOrStdout())
			} else {
				err = createTrigger(client, trigger, namespace, cmd.OutOrStdout())
			}
			if err != nil {
				return err
			}

			return nil
		},
	}
	commands.AddNamespaceFlags(connCreateCommand.Flags(), false)
	editFlags.AddCreateFlags(connCreateCommand)
	waitFlags.AddConditionWaitFlags(connCreateCommand, 60, "connection")
	return connCreateCommand
}

// Duck type for writers having a flush
type flusher interface {
	Flush() error
}

func flush(out io.Writer) {
	if flusher, ok := out.(flusher); ok {
		flusher.Flush()
	}
}

func createTrigger(client v1alpha1.KnEventClient, trigger *eventing_v1alpha1_api.Trigger, namespace string, out io.Writer) error {
	err := client.CreateTrigger(trigger)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Connection '%s' successfully created in namespace '%s'.\n", trigger.Name, namespace)
	return nil
}

func replaceTrigger(client v1alpha1.KnEventClient, trigger *eventing_v1alpha1_api.Trigger, namespace string, out io.Writer) error {
	var retries = 0
	for {
		existingTrigger, err := client.GetTrigger(trigger.Name)
		if err != nil {
			return err
		}
		trigger.ResourceVersion = existingTrigger.ResourceVersion
		err = client.UpdateTrigger(trigger)
		if err != nil {
			// Retry to update when a resource version conflict exists
			if api_errors.IsConflict(err) && retries < MaxUpdateRetries {
				retries++
				continue
			}
			return err
		}
		fmt.Fprintf(out, "Connection '%s' successfully replaced in namespace '%s'.\n", trigger.Name, namespace)
		return nil
	}
}

func triggerExists(client v1alpha1.KnEventClient, name string, namespace string) (bool, error) {
	_, err := client.GetTrigger(name)
	if api_errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Create trigger struct from provided options
func constructTrigger(cmd *cobra.Command, editFlags ConfigurationEditFlags, name string, namespace string) (*eventing_v1alpha1_api.Trigger,
	error) {

	trigger := eventing_v1alpha1_api.Trigger{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: eventing_v1alpha1_api.TriggerSpec{},
	}

	err := editFlags.Apply(&trigger, cmd)
	if err != nil {
		return nil, err
	}
	return &trigger, nil
}
