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
	hprinters "github.com/knative/client/pkg/printers"
	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ConnectionListHandlers adds print handlers for connection list command
func ConnectionListHandlers(h hprinters.PrintHandler) {
	ConnColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Name", Type: "string", Description: "Name of the connection."},
		{Name: "Subscriber", Type: "string", Description: "Name of the resource recieving events."},
		{Name: "Broker", Type: "string", Description: "Name of the connection."},
		{Name: "Age", Type: "string", Description: "Age of the connection."},
		{Name: "Conditions", Type: "string", Description: "Conditions describing statuses of the connection."},
		{Name: "Ready", Type: "string", Description: "Ready condition status of the connection."},
		{Name: "Reason", Type: "string", Description: "Reason for non-ready condition of the connection."},
	}
	h.TableHandler(ConnColumnDefinitions, printConn)
	h.TableHandler(ConnColumnDefinitions, printConnList)
}

// Private functions

// printConnList populates the Knative connection list table rows
func printConnList(triggerList *eventingv1alpha1.TriggerList, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(triggerList.Items))
	for _, c := range triggerList.Items {
		r, err := printConn(&c, options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

// printConn populates the Knative connection table rows
func printConn(trigger *eventingv1alpha1.Trigger, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	name := trigger.Name
	br := GetBrokerName(trigger)
	sub := GetSubscriberName(trigger)
	age := commands.TranslateTimestampSince(trigger.CreationTimestamp)
	conditions := commands.ConditionsValue(trigger.Status.Conditions)
	ready := commands.ReadyCondition(trigger.Status.Conditions)
	reason := commands.NonReadyConditionReason(trigger.Status.Conditions)
	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: trigger},
	}
	row.Cells = append(row.Cells,
		name,
		sub,
		br,
		age,
		conditions,
		ready,
		reason)
	return []metav1beta1.TableRow{row}, nil
}
