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
	"github.com/knative/client/pkg/kn/commands"
	hprinters "github.com/knative/client/pkg/printers"
	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// BrokerListHandlers adds print handlers for broker list command
func BrokerListHandlers(h hprinters.PrintHandler) {
	BrokerColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Name", Type: "string", Description: "Name of the broker."},
		{Name: "Age", Type: "string", Description: "Age of the broker."},
		{Name: "Conditions", Type: "string", Description: "Conditions describing statuses of the broker."},
		{Name: "Ready", Type: "string", Description: "Ready condition status of the broker."},
		{Name: "Reason", Type: "string", Description: "Reason for non-ready condition of the broker."},
	}
	h.TableHandler(BrokerColumnDefinitions, printBroker)
	h.TableHandler(BrokerColumnDefinitions, printBrokerList)
}

// Private functions

// printBrokerList populates the Knative broker list table rows
func printBrokerList(brokerList *eventingv1alpha1.BrokerList, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(brokerList.Items))
	for _, br := range brokerList.Items {
		r, err := printBroker(&br, options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

// printBroker populates the Knative broker table rows
func printBroker(broker *eventingv1alpha1.Broker, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	name := broker.Name
	age := commands.TranslateTimestampSince(broker.CreationTimestamp)
	conditions := commands.ConditionsValue(broker.Status.Conditions)
	ready := commands.ReadyCondition(broker.Status.Conditions)
	reason := commands.NonReadyConditionReason(broker.Status.Conditions)
	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: broker},
	}
	row.Cells = append(row.Cells,
		name,
		age,
		conditions,
		ready,
		reason)
	return []metav1beta1.TableRow{row}, nil
}
