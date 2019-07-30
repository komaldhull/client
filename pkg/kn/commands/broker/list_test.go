// Copyright © 2018 The Knative Authors
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
	"strings"
	"testing"

	"github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"gotest.tools/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	client_testing "k8s.io/client-go/testing"

	"github.com/knative/client/pkg/kn/commands"
	"github.com/knative/client/pkg/util"
)

func fakeBrokerList(args []string, response *v1alpha1.BrokerList) (action client_testing.Action, output []string, err error) {
	knParams := &commands.KnParams{}
	cmd, fakeEventing, buf := commands.CreateTestKnEventCommand(NewBrokerCommand(knParams), knParams)
	fakeEventing.AddReactor("list", "*",
		func(a client_testing.Action) (bool, runtime.Object, error) {
			action = a
			return true, response, nil
		})
	cmd.SetArgs(args)
	err = cmd.Execute()
	if err != nil {
		return
	}
	output = strings.Split(buf.String(), "\n")
	return
}

func TestBrokerListEmpty(t *testing.T) {
	action, output, err := fakeBrokerList([]string{"broker", "list"}, &v1alpha1.BrokerList{})
	if err != nil {
		t.Error(err)
		return
	}
	if action == nil {
		t.Errorf("No action")
	} else if !action.Matches("list", "brokers") {
		t.Errorf("Bad action %v", action)
	} else if output[0] != "No resources found." {
		t.Errorf("Bad output %s", output[0])
	}
}

func TestBrokerListEmptyByName(t *testing.T) {
	action, _, err := fakeBrokerList([]string{"broker", "list", "name"}, &v1alpha1.BrokerList{})
	assert.NilError(t, err)
	if action == nil {
		t.Errorf("No action")
	} else if !action.Matches("list", "brokers") {
		t.Errorf("Bad action %v", action)
	}
}

func TestBrokerListDefaultOutput(t *testing.T) {
	broker1 := createMockBrokerWithParams("foo")
	broker2 := createMockBrokerWithParams("bar")
	BrokerList := &v1alpha1.BrokerList{Items: []v1alpha1.Broker{*broker1, *broker2}}
	action, output, err := fakeBrokerList([]string{"broker", "list"}, BrokerList)
	assert.NilError(t, err)
	if action == nil {
		t.Errorf("No action")
	} else if !action.Matches("list", "brokers") {
		t.Errorf("Bad action %v", action)
	}
	assert.Check(t, util.ContainsAll(output[0], "NAME", "AGE", "CONDITIONS", "READY", "REASON"))
	assert.Check(t, util.ContainsAll(output[1], "foo"))
	assert.Check(t, util.ContainsAll(output[2], "bar"))
}
func TestBrokerListOneOutput(t *testing.T) {
	broker := createMockBrokerWithParams("foo")
	BrokerList := &v1alpha1.BrokerList{Items: []v1alpha1.Broker{*broker}}
	action, output, err := fakeBrokerList([]string{"broker", "list", "foo"}, BrokerList)
	assert.NilError(t, err)
	if action == nil {
		t.Errorf("No action")
	} else if !action.Matches("list", "brokers") {
		t.Errorf("Bad action %v", action)
	}

	assert.Check(t, util.ContainsAll(output[0], "NAME", "AGE", "CONDITIONS", "READY", "REASON"))
	assert.Check(t, util.ContainsAll(output[1], "foo"))
}

func createMockBrokerWithParams(name string) *v1alpha1.Broker {
	broker := &v1alpha1.Broker{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Broker",
			APIVersion: "knative.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
	}
	return broker
}
