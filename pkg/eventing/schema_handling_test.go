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

package eventing

import (
	"testing"

	"gotest.tools/assert"

	"github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestGVKUpdate(t *testing.T) {
	broker := v1alpha1.Broker{}
	err := UpdateGroupVersionKind(&broker, v1alpha1.SchemeGroupVersion)
	assert.NilError(t, err)
	assert.Equal(t, broker.Kind, "Broker")
	assert.Equal(t, broker.APIVersion, v1alpha1.SchemeGroupVersion.Group+"/"+v1alpha1.SchemeGroupVersion.Version)
}

func TestGVKUpdateNegative(t *testing.T) {
	broker := v1alpha1.Broker{}
	err := UpdateGroupVersionKind(&broker, schema.GroupVersion{Group: "bla", Version: "blub"})
	assert.ErrorContains(t, err, "group version")
}
