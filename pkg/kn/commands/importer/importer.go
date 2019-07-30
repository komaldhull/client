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
	"github.com/knative/client/pkg/kn/commands"
	sourcesv1alpha1 "github.com/knative/eventing/pkg/apis/sources/v1alpha1"
	"github.com/spf13/cobra"
)

const (
	// How often to retry in case of an optimistic lock error when replacing an importer (--force)
	MaxUpdateRetries = 3
)

func NewImporterCommand(p *commands.KnParams) *cobra.Command {
	importerCmd := &cobra.Command{
		Use:   "importer",
		Short: "Importer command group",
	}
	importerCmd.AddCommand(NewImporterListCommand(p))
	importerCmd.AddCommand(NewImporterDescribeCommand(p))
	importerCmd.AddCommand(NewImporterCreateCommand(p))
	importerCmd.AddCommand(NewImporterUpdateCommand(p))
	importerCmd.AddCommand(NewImporterDeleteCommand(p))
	return importerCmd
}

func GetSink(src *sourcesv1alpha1.ContainerSource) string {
	if src.Spec.Sink == nil {
		return ""
	}
	return src.Spec.Sink.Name
}
