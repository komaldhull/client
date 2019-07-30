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
	"encoding/json"
	"testing"

	"gotest.tools/assert"

	"github.com/knative/client/pkg/kn/commands"
	"github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	client_testing "k8s.io/client-go/testing"
	"sigs.k8s.io/yaml"
)

func fakeBrokerDescribe(args []string, response *v1alpha1.Broker) (action client_testing.Action, output string, err error) {
	knParams := &commands.KnParams{}
	cmd, fakeEventing, buf := commands.CreateTestKnEventCommand(NewBrokerCommand(knParams), knParams)
	fakeEventing.AddReactor("*", "*",
		func(a client_testing.Action) (bool, runtime.Object, error) {
			action = a
			return true, response, nil
		})
	cmd.SetArgs(args)
	err = cmd.Execute()
	if err != nil {
		return
	}
	output = buf.String()
	return
}

func TestDescribeBrokerWithNoName(t *testing.T) {
	_, _, err := fakeBrokerDescribe([]string{"broker", "describe"}, &v1alpha1.Broker{})
	expectedError := "requires the broker name"
	assert.ErrorContains(t, err, expectedError)
}

func TestDescribeBrokerYaml(t *testing.T) {
	expectedBroker := v1alpha1.Broker{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Broker",
			APIVersion: "knative.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
	}

	action, data, err := fakeBrokerDescribe([]string{"broker", "describe", "foo"}, &expectedBroker)
	if err != nil {
		t.Fatal(err)
	}

	if action == nil {
		t.Fatal("No action")
	} else if !action.Matches("get", "brokers") {
		t.Fatalf("Bad action %v", action)
	}

	jsonData, err := yaml.YAMLToJSON([]byte(data))
	assert.NilError(t, err)

	var returnedBroker v1alpha1.Broker
	err = json.Unmarshal(jsonData, &returnedBroker)
	assert.NilError(t, err)
	assert.DeepEqual(t, expectedBroker, returnedBroker)

	if !equality.Semantic.DeepEqual(expectedBroker, returnedBroker) {
		t.Fatal("mismatched objects")
	}
}
