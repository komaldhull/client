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
	"github.com/knative/client/pkg/kn/commands"
	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"github.com/spf13/cobra"
)

const (
	// How often to retry in case of an optimistic lock error when replacing a trigger (--force)
	MaxUpdateRetries = 3
)

func NewConnectionCommand(p *commands.KnParams) *cobra.Command {
	connectionCmd := &cobra.Command{
		Use:   "connection",
		Short: "Connection command group",
	}
	connectionCmd.AddCommand(NewConnectionListCommand(p))
	connectionCmd.AddCommand(NewConnectionDescribeCommand(p))
	connectionCmd.AddCommand(NewConnectionCreateCommand(p))
	connectionCmd.AddCommand(NewConnectionUpdateCommand(p))
	connectionCmd.AddCommand(NewConnectionDeleteCommand(p))
	return connectionCmd
}

func GetBrokerName(kTrigger *eventingv1alpha1.Trigger) string {
	return kTrigger.Spec.Broker

}
func GetSubscriberName(kTrigger *eventingv1alpha1.Trigger) string {
	return kTrigger.Spec.Subscriber.Ref.Name

}
