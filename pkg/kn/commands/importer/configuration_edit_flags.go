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

package importer

import (
	"errors"
	"strings"

	sourceslib "github.com/knative/client/pkg/eventing/sourcesv1alpha1"
	sourcesv1alpha1 "github.com/knative/eventing/pkg/apis/sources/v1alpha1"
	"github.com/spf13/cobra"
)

type ConfigurationEditFlags struct {
	Type        string
	Image       string
	Sink        string
	ForceCreate bool
}

func (p *ConfigurationEditFlags) AddUpdateFlags(command *cobra.Command) {
	command.Flags().StringVar(&p.Image, "image", "", "If you want to create a containersource, image to run.")
	command.Flags().StringVar(&p.Type, "type", "", "Type of the source. Currently supported option: 'container'")
	command.Flags().StringVar(&p.Sink, "sink", "", "Name and type of the sink, ex. 'broker:default'. Defaults to namespace broker when not specified ")

}

func (p *ConfigurationEditFlags) AddCreateFlags(command *cobra.Command) {
	p.AddUpdateFlags(command)
	command.MarkFlagRequired("type")
	command.Flags().BoolVar(&p.ForceCreate, "force", false, "Create source forcefully, replaces existing source if any.")
}

func (p *ConfigurationEditFlags) Apply(src *sourcesv1alpha1.ContainerSource, cmd *cobra.Command) error {

	if cmd.Flags().Changed("image") {
		err := sourceslib.UpdateImage(src, p.Image)
		if err != nil {
			return err
		}
	}
	if p.Sink != "" {
		sink := strings.Split(p.Sink, ":")
		if len(sink) < 2 || (sink[0] != "broker" && sink[0] != "service" && sink[0] != "sequence") {
			return errors.New("incorrect format for sink, specify 'type:name', ex. 'broker:default'")
		}
		if sink[0] == "broker" {
			err := sourceslib.UpdateBrokerSink(src, sink[1])
			if err != nil {
				return err
			}
		} else if sink[0] == "service" {
			err := sourceslib.UpdateSvcSink(src, sink[1])
			if err != nil {
				return err
			}
		} else {
			err := sourceslib.UpdateSeqSink(src, sink[1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
