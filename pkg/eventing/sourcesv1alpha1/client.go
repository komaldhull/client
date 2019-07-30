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
	"github.com/knative/client/pkg/eventing"
	servingv1alpha1 "github.com/knative/client/pkg/serving/v1alpha1"

	"github.com/knative/eventing/pkg/apis/sources/v1alpha1"
	client_v1alpha1 "github.com/knative/eventing/pkg/client/clientset/versioned/typed/sources/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Kn interface to eventing sources. All methods are relative to the
// namespace specified during construction
type KnSourceClient interface {

	// List container sources
	ListContainerSources(opts ...servingv1alpha1.ListConfig) (*v1alpha1.ContainerSourceList, error)

	// Get a container source by its unique name
	GetContainerSource(name string) (*v1alpha1.ContainerSource, error)

	// Create a new container source
	CreateContainerSource(trigger *v1alpha1.ContainerSource) error

	// Update the given container source
	UpdateContainerSource(trigger *v1alpha1.ContainerSource) error

	// Delete a container source by name
	DeleteContainerSource(name string) error
}

type knSourceClient struct {
	client    client_v1alpha1.SourcesV1alpha1Interface
	namespace string
}

// Create a new client facade for the provided namespace
func NewKnSourceClient(client client_v1alpha1.SourcesV1alpha1Interface, namespace string) KnSourceClient {
	return &knSourceClient{
		client:    client,
		namespace: namespace,
	}
}

// Get a container source by its unique name
func (cl *knSourceClient) GetContainerSource(name string) (*v1alpha1.ContainerSource, error) {
	src, err := cl.client.ContainerSources(cl.namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	err = updateEventingGvk(src)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// List Container Sources
func (cl *knSourceClient) ListContainerSources(config ...servingv1alpha1.ListConfig) (*v1alpha1.ContainerSourceList, error) {
	srcList, err := cl.client.ContainerSources(cl.namespace).List(servingv1alpha1.ListConfigs(config).ToListOptions())
	if err != nil {
		return nil, err
	}
	return updateEventingGvkForContainerSourceList(srcList)

}

// Create a new Container Source
func (cl *knSourceClient) CreateContainerSource(src *v1alpha1.ContainerSource) error {
	_, err := cl.client.ContainerSources(cl.namespace).Create(src)
	if err != nil {
		return err
	}
	return updateEventingGvk(src)
}

// Update the given Container Source
func (cl *knSourceClient) UpdateContainerSource(src *v1alpha1.ContainerSource) error {
	_, err := cl.client.ContainerSources(cl.namespace).Update(src)
	if err != nil {
		return err
	}
	return updateEventingGvk(src)
}

// Delete a service by name
func (cl *knSourceClient) DeleteContainerSource(srcName string) error {
	return cl.client.ContainerSources(cl.namespace).Delete(
		srcName,
		&v1.DeleteOptions{},
	)
}

//Private fns

// update all the list + all items contained in the list with
// the proper GroupVersionKind specific to Knative eventing
func updateEventingGvkForContainerSourceList(srcList *v1alpha1.ContainerSourceList) (*v1alpha1.ContainerSourceList, error) {
	srcListNew := srcList.DeepCopy()
	err := updateEventingGvk(srcListNew)
	if err != nil {
		return nil, err
	}

	srcListNew.Items = make([]v1alpha1.ContainerSource, len(srcList.Items))
	for idx := range srcList.Items {
		src := srcList.Items[idx].DeepCopy()
		err := updateEventingGvk(src)
		if err != nil {
			return nil, err
		}
		srcListNew.Items[idx] = *src
	}
	return srcListNew, nil
}

// update with the v1alpha1 group + version
func updateEventingGvk(obj runtime.Object) error {
	return eventing.UpdateGroupVersionKind(obj, v1alpha1.SchemeGroupVersion)
}
