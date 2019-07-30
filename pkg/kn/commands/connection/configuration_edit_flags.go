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

package connection

import (
	eventinglib "github.com/knative/client/pkg/eventing/v1alpha1"
	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"github.com/spf13/cobra"
)

type ConfigurationEditFlags struct {
	Service     string
	Sequence    string
	Broker      string
	Source      string
	Type        string
	ForceCreate bool
}

func (p *ConfigurationEditFlags) AddUpdateFlags(command *cobra.Command) {
	command.Flags().StringVar(&p.Service, "service", "", "Service that is the subscriber")
	command.Flags().StringVar(&p.Sequence, "sequence", "", "Sequence that is the subscriber")
	command.Flags().StringVar(&p.Source, "source", "", "Event source filter")
	command.Flags().StringVar(&p.Type, "type", "", "Event type filter")

}

func (p *ConfigurationEditFlags) AddCreateFlags(command *cobra.Command) {
	p.AddUpdateFlags(command)
	command.Flags().BoolVar(&p.ForceCreate, "force", false, "Create trigger forcefully, replaces existing trigger if any.")
	command.Flags().StringVar(&p.Broker, "broker", "", "Broker to subscribe to")
	command.MarkFlagRequired("broker")

}

func (p *ConfigurationEditFlags) Apply(trigger *eventingv1alpha1.Trigger, cmd *cobra.Command) error {

	if cmd.Flags().Changed("service") {
		err := eventinglib.UpdateSvcSubscriber(trigger, p.Service)
		if err != nil {
			return err
		}
	} else if cmd.Flags().Changed("sequence") {
		err := eventinglib.UpdateSeqSubscriber(trigger, p.Sequence)
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("source") {
		err := eventinglib.UpdateSourceFilter(trigger, p.Source)
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("type") {
		err := eventinglib.UpdateTypeFilter(trigger, p.Type)
		if err != nil {
			return err
		}
	}

	return nil
}
