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
	"github.com/knative/client/pkg/kn/commands"
	hprinters "github.com/knative/client/pkg/printers"
	sourcesv1alpha1 "github.com/knative/eventing/pkg/apis/sources/v1alpha1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func SourceListHandlers(h hprinters.PrintHandler) {
	kSourceColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Name", Type: "string", Description: "Name of the knative source."},
		{Name: "Type", Type: "string", Description: "Type of the source."},
		{Name: "Sink", Type: "string", Description: "Sink name of the source."},
		{Name: "Age", Type: "string", Description: "Age of the source."},
		{Name: "Conditions", Type: "string", Description: "Conditions describing statuses of source components."},
		{Name: "Ready", Type: "string", Description: "Ready condition status of the source."},
		{Name: "Reason", Type: "string", Description: "Reason for non-ready condition of the source."},
	}
	h.TableHandler(kSourceColumnDefinitions, printKSource)
	h.TableHandler(kSourceColumnDefinitions, printKSourceList)
}

// printKSourceList populates the knative Source list table rows
func printKSourceList(kSourceList *sourcesv1alpha1.ContainerSourceList, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(kSourceList.Items))
	for _, ks := range kSourceList.Items {
		r, err := printKSource(&ks, options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

// printKSource populates the knative source table rows
func printKSource(kSource *sourcesv1alpha1.ContainerSource, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	name := kSource.Name
	age := commands.TranslateTimestampSince(kSource.CreationTimestamp)
	conditions := commands.ConditionsValue(kSource.Status.Status.Conditions)
	ready := commands.ReadyCondition(kSource.Status.Conditions)
	reason := commands.NonReadyConditionReason(kSource.Status.Conditions)
	sink := GetSink(kSource)

	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: kSource},
	}
	row.Cells = append(row.Cells,
		name,
		"ContainerSource",
		sink,
		age,
		conditions,
		ready,
		reason)

	return []metav1beta1.TableRow{row}, nil
}
