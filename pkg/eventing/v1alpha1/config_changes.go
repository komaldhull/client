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

package v1alpha1

import (
	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

//TODO: shouldn't hardcode apiversion
func UpdateSvcSubscriber(trigger *eventingv1alpha1.Trigger, sub string) error {
	trigger.Spec.Subscriber = &eventingv1alpha1.SubscriberSpec{
		Ref: &corev1.ObjectReference{
			Kind:       "Service",
			APIVersion: "serving.knative.dev/v1alpha1",
			Name:       sub,
		},
	}
	return nil
}

//TODO: shouldn't hardcode apiversion
func UpdateSeqSubscriber(trigger *eventingv1alpha1.Trigger, sub string) error {
	trigger.Spec.Subscriber = &eventingv1alpha1.SubscriberSpec{
		Ref: &corev1.ObjectReference{
			Kind:       "Sequence",
			APIVersion: "eventing.knative.dev/v1alpha1",
			Name:       sub,
		},
	}
	return nil
}

func UpdateSourceFilter(trigger *eventingv1alpha1.Trigger, src string) error {
	if trigger.Spec.Filter == nil {
		trigger.Spec.Filter = &eventingv1alpha1.TriggerFilter{
			SourceAndType: &eventingv1alpha1.TriggerFilterSourceAndType{
				Source: src,
			},
		}
	} else {
		trigger.Spec.Filter.SourceAndType.Source = src
	}
	return nil
}

func UpdateTypeFilter(trigger *eventingv1alpha1.Trigger, typ string) error {
	if trigger.Spec.Filter == nil {
		trigger.Spec.Filter = &eventingv1alpha1.TriggerFilter{
			SourceAndType: &eventingv1alpha1.TriggerFilterSourceAndType{
				Type: typ,
			},
		}
	} else {
		trigger.Spec.Filter.SourceAndType.Type = typ
	}

	return nil
}
