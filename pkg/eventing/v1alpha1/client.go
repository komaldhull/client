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
	"fmt"

	"github.com/knative/pkg/apis"

	"github.com/knative/client/pkg/eventing"
	servingv1alpha1 "github.com/knative/client/pkg/serving/v1alpha1"

	"github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	client_v1alpha1 "github.com/knative/eventing/pkg/client/clientset/versioned/typed/eventing/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Kn interface to eventing. All methods are relative to the
// namespace specified during construction
type KnEventClient interface {

	// List brokers
	ListBrokers(opts ...servingv1alpha1.ListConfig) (*v1alpha1.BrokerList, error)

	// Get a broker by its unique name
	GetBroker(name string) (*v1alpha1.Broker, error)

	// List triggers
	ListTriggers(opts ...servingv1alpha1.ListConfig) (*v1alpha1.TriggerList, error)

	// Get a trigger by its unique name
	GetTrigger(name string) (*v1alpha1.Trigger, error)

	// Create a new trigger
	CreateTrigger(trigger *v1alpha1.Trigger) error

	// Update the given trigger
	UpdateTrigger(trigger *v1alpha1.Trigger) error

	// Delete a trigger by name
	DeleteTrigger(name string) error
}

type knEventClient struct {
	client    client_v1alpha1.EventingV1alpha1Interface
	namespace string
}

// Create a new client facade for the provided namespace
func NewKnEventClient(client client_v1alpha1.EventingV1alpha1Interface, namespace string) KnEventClient {
	return &knEventClient{
		client:    client,
		namespace: namespace,
	}
}

// Get a broker by its unique name
func (cl *knEventClient) GetBroker(name string) (*v1alpha1.Broker, error) {
	broker, err := cl.client.Brokers(cl.namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	err = updateEventingGvk(broker)
	if err != nil {
		return nil, err
	}
	return broker, nil
}

// List brokers
func (cl *knEventClient) ListBrokers(config ...servingv1alpha1.ListConfig) (*v1alpha1.BrokerList, error) {
	brokerList, err := cl.client.Brokers(cl.namespace).List(servingv1alpha1.ListConfigs(config).ToListOptions())
	if err != nil {
		return nil, err
	}
	return updateEventingGvkForBrokerList(brokerList)

}

// Get a trigger by its unique name
func (cl *knEventClient) GetTrigger(name string) (*v1alpha1.Trigger, error) {
	trigger, err := cl.client.Triggers(cl.namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	err = updateEventingGvk(trigger)
	if err != nil {
		return nil, err
	}
	return trigger, nil
}

// List triggers
func (cl *knEventClient) ListTriggers(config ...servingv1alpha1.ListConfig) (*v1alpha1.TriggerList, error) {
	triggerList, err := cl.client.Triggers(cl.namespace).List(servingv1alpha1.ListConfigs(config).ToListOptions())
	if err != nil {
		return nil, err
	}
	return updateEventingGvkForTriggerList(triggerList)

}

// Create a new trigger
func (cl *knEventClient) CreateTrigger(trigger *v1alpha1.Trigger) error {
	_, err := cl.client.Triggers(cl.namespace).Create(trigger)
	if err != nil {
		return err
	}
	return updateEventingGvk(trigger)
}

// Update the given trigger
func (cl *knEventClient) UpdateTrigger(trigger *v1alpha1.Trigger) error {
	_, err := cl.client.Triggers(cl.namespace).Update(trigger)
	if err != nil {
		return err
	}
	return updateEventingGvk(trigger)
}

// Delete a service by name
func (cl *knEventClient) DeleteTrigger(triggerName string) error {
	return cl.client.Triggers(cl.namespace).Delete(
		triggerName,
		&v1.DeleteOptions{},
	)
}

//Private fns

// update all the list + all items contained in the list with
// the proper GroupVersionKind specific to Knative eventing
func updateEventingGvkForBrokerList(brokerList *v1alpha1.BrokerList) (*v1alpha1.BrokerList, error) {
	brokerListNew := brokerList.DeepCopy()
	err := updateEventingGvk(brokerListNew)
	if err != nil {
		return nil, err
	}

	brokerListNew.Items = make([]v1alpha1.Broker, len(brokerList.Items))
	for idx := range brokerList.Items {
		br := brokerList.Items[idx].DeepCopy()
		err := updateEventingGvk(br)
		if err != nil {
			return nil, err
		}
		brokerListNew.Items[idx] = *br
	}
	return brokerListNew, nil
}

func updateEventingGvkForTriggerList(triggerList *v1alpha1.TriggerList) (*v1alpha1.TriggerList, error) {
	triggerListNew := triggerList.DeepCopy()
	err := updateEventingGvk(triggerListNew)
	if err != nil {
		return nil, err
	}

	triggerListNew.Items = make([]v1alpha1.Trigger, len(triggerList.Items))
	for idx := range triggerList.Items {
		tri := triggerList.Items[idx].DeepCopy()
		err := updateEventingGvk(tri)
		if err != nil {
			return nil, err
		}
		triggerListNew.Items[idx] = *tri
	}
	return triggerListNew, nil
}

// update with the v1alpha1 group + version
func updateEventingGvk(obj runtime.Object) error {
	return eventing.UpdateGroupVersionKind(obj, v1alpha1.SchemeGroupVersion)
}

func brokerConditionExtractor(obj runtime.Object) (apis.Conditions, error) {
	broker, ok := obj.(*v1alpha1.Broker)
	if !ok {
		return nil, fmt.Errorf("%v is not a broker", obj)
	}
	return apis.Conditions(broker.Status.Conditions), nil
}
