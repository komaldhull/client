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
	"errors"
	"fmt"
	"io"

	v1alpha1 "github.com/knative/client/pkg/eventing/sourcesv1alpha1"
	"github.com/knative/client/pkg/kn/commands"

	sources_v1alpha1_api "github.com/knative/eventing/pkg/apis/sources/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewImporterCreateCommand(p *commands.KnParams) *cobra.Command {
	var editFlags ConfigurationEditFlags

	importerCreateCommand := &cobra.Command{
		Use:   "create NAME",
		Short: "Create an importer.",

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("'importer create' requires the importer name given as single argument")
			}

			if editFlags.Type == "" {
				return errors.New("'importer create' requires the source type")
			}

			if editFlags.Type != "container" {
				return errors.New("unsupported importer type")
			}

			if editFlags.Image == "" {
				return errors.New("Container source requires the image")
			}

			name := args[0]
			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}

			src, err := constructSrc(cmd, editFlags, args[0], namespace)
			if err != nil {
				return err
			}

			client, err := p.NewSourceClient(namespace)
			if err != nil {
				return err
			}

			srcExists, err := srcExists(client, name, namespace)
			if err != nil {
				return err
			}

			if srcExists {
				if !editFlags.ForceCreate {
					return fmt.Errorf(
						"cannot create source '%s' in namespace '%s' "+
							"because the source already exists and no --force option was given", name, namespace)
				}
				err = replaceSrc(client, src, namespace, cmd.OutOrStdout())
			} else {
				err = createSrc(client, src, namespace, cmd.OutOrStdout())
			}
			if err != nil {
				return err
			}

			return nil
		},
	}
	commands.AddNamespaceFlags(importerCreateCommand.Flags(), false)
	editFlags.AddCreateFlags(importerCreateCommand)
	return importerCreateCommand
}

// Duck type for writers having a flush
type flusher interface {
	Flush() error
}

func flush(out io.Writer) {
	if flusher, ok := out.(flusher); ok {
		flusher.Flush()
	}
}

func createSrc(client v1alpha1.KnSourceClient, src *sources_v1alpha1_api.ContainerSource, namespace string, out io.Writer) error {
	err := client.CreateContainerSource(src)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Source '%s' successfully created in namespace '%s'.\n", src.Name, namespace)
	return nil
}

func replaceSrc(client v1alpha1.KnSourceClient, src *sources_v1alpha1_api.ContainerSource, namespace string, out io.Writer) error {
	var retries = 0
	for {
		existingSrc, err := client.GetContainerSource(src.Name)
		if err != nil {
			return err
		}
		src.ResourceVersion = existingSrc.ResourceVersion
		err = client.UpdateContainerSource(src)
		if err != nil {
			// Retry to update when a resource version conflict exists
			if api_errors.IsConflict(err) && retries < MaxUpdateRetries {
				retries++
				continue
			}
			return err
		}
		fmt.Fprintf(out, "Source '%s' successfully replaced in namespace '%s'.\n", src.Name, namespace)
		return nil
	}
}

func srcExists(client v1alpha1.KnSourceClient, name string, namespace string) (bool, error) {
	_, err := client.GetContainerSource(name)
	if api_errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Create container source struct from provided options
func constructSrc(cmd *cobra.Command, editFlags ConfigurationEditFlags, name string, namespace string) (*sources_v1alpha1_api.ContainerSource,
	error) {

	src := sources_v1alpha1_api.ContainerSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sources_v1alpha1_api.ContainerSourceSpec{
			Template: &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: editFlags.Image,
						Name:  name,
					}},
				},
			},
			Sink: &corev1.ObjectReference{
				Kind:       "Broker",
				APIVersion: "eventing.knative.dev/v1alpha1",
				Name:       "default",
			},
		},
	}

	err := editFlags.Apply(&src, cmd)
	if err != nil {
		return nil, err
	}
	return &src, nil
}
