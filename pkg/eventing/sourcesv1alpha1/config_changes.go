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

package sourcesv1alpha1

import (
	sourcesv1alpha1 "github.com/knative/eventing/pkg/apis/sources/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

//TODO: shouldn't hardcode apiversion
func UpdateImage(src *sourcesv1alpha1.ContainerSource, img string) error {
	src.Spec.Template.Spec.Containers[0].Image = img
	return nil
}

func UpdateSvcSink(src *sourcesv1alpha1.ContainerSource, name string) error {
	src.Spec.Sink = &corev1.ObjectReference{
		Kind:       "Service",
		APIVersion: "serving.knative.dev/v1alpha1",
		Name:       name,
	}
	return nil
}

func UpdateSeqSink(src *sourcesv1alpha1.ContainerSource, name string) error {
	src.Spec.Sink = &corev1.ObjectReference{
		Kind:       "Sequence",
		APIVersion: "messaging.knative.dev/v1alpha1",
		Name:       name,
	}
	return nil
}

func DeleteSink(src *sourcesv1alpha1.ContainerSource) error {
	src.Spec.Sink = nil
	return nil
}

func UpdateBrokerSink(src *sourcesv1alpha1.ContainerSource, name string) error {
	src.Spec.Sink = &corev1.ObjectReference{
		Kind:       "Broker",
		APIVersion: "eventing.knative.dev/v1alpha1",
		Name:       name,
	}
	return nil
}
