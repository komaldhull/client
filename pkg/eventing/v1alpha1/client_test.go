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
	"testing"

	"gotest.tools/assert"

	"github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"github.com/knative/eventing/pkg/client/clientset/versioned/typed/eventing/v1alpha1/fake"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	client_testing "k8s.io/client-go/testing"

	"github.com/knative/client/pkg/eventing"
)

var testNamespace = "test-ns"

func setup() (eventing fake.FakeEventingV1alpha1, client KnEventClient) {
	eventing = fake.FakeEventingV1alpha1{Fake: &client_testing.Fake{}}
	client = NewKnEventClient(&eventing, testNamespace)
	return
}

func TestGetBroker(t *testing.T) {
	eventing, client := setup()
	brokerName := "test-broker"

	eventing.AddReactor("get", "brokers",
		func(a client_testing.Action) (bool, runtime.Object, error) {
			broker := newBroker(brokerName)
			name := a.(client_testing.GetAction).GetName()
			// Sanity check
			assert.Assert(t, name != "")
			assert.Equal(t, testNamespace, a.GetNamespace())
			if name == brokerName {
				return true, broker, nil
			}
			return true, nil, errors.NewNotFound(v1alpha1.Resource("broker"), name)
		})

	t.Run("get known broker by name returns broker", func(t *testing.T) {
		broker, err := client.GetBroker(brokerName)
		assert.NilError(t, err)
		assert.Equal(t, brokerName, broker.Name, "broker name should be equal")
		validateGroupVersionKind(t, broker)
	})

	t.Run("get unknown broker name returns error", func(t *testing.T) {
		nonExistingBrokerName := "nonexistent-broker"
		broker, err := client.GetBroker(nonExistingBrokerName)
		assert.Assert(t, broker == nil, "no broker should be returned")
		assert.ErrorContains(t, err, "not found")
		assert.ErrorContains(t, err, nonExistingBrokerName)
	})
}

func TestListBroker(t *testing.T) {
	eventing, client := setup()

	t.Run("list broker returns a list of brokers", func(t *testing.T) {
		broker1 := newBroker("broker-1")
		broker2 := newBroker("broker-2")

		eventing.AddReactor("list", "brokers",
			func(a client_testing.Action) (bool, runtime.Object, error) {
				assert.Equal(t, testNamespace, a.GetNamespace())
				return true, &v1alpha1.BrokerList{Items: []v1alpha1.Broker{*broker1, *broker2}}, nil
			})

		listBrokers, err := client.ListBrokers()
		assert.NilError(t, err)
		assert.Assert(t, len(listBrokers.Items) == 2)
		assert.Equal(t, listBrokers.Items[0].Name, "broker-1")
		assert.Equal(t, listBrokers.Items[1].Name, "broker-2")
		validateGroupVersionKind(t, listBrokers)
		validateGroupVersionKind(t, &listBrokers.Items[0])
		validateGroupVersionKind(t, &listBrokers.Items[1])
	})
}

func validateGroupVersionKind(t *testing.T, obj runtime.Object) {
	gvkExpected, err := eventing.GetGroupVersionKind(obj, v1alpha1.SchemeGroupVersion)
	assert.NilError(t, err)
	gvkGiven := obj.GetObjectKind().GroupVersionKind()
	assert.Equal(t, *gvkExpected, gvkGiven, "GVK should be the same")
}

func newBroker(name string) *v1alpha1.Broker {
	return &v1alpha1.Broker{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace}}
}
