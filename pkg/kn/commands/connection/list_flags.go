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
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or im
// See the License for the specific language governing permissions and
// limitations under the License.

package connection

import (
	"github.com/knative/client/pkg/kn/commands"
	hprinters "github.com/knative/client/pkg/printers"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// ListFilterFlags consists of flags used to filter the list of connections
type ListFilterFlags struct {
	Subscriber string
	Broker     string
}

// AddFlags receives a *cobra.Command reference and
// adds the ListFilter flags as optional flags
func (f *ListFilterFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Subscriber, "subscriber", "s", "", "Subscriber name")
	cmd.Flags().StringVarP(&f.Broker, "broker", "b", "", "Broker name")
}

// ConnListFlags composes common printer flag structs
// used in the List command.
type ConnListFlags struct {
	GenericPrintFlags  *genericclioptions.PrintFlags
	HumanReadableFlags *commands.HumanPrintFlags
	FilterFlags        ListFilterFlags
}

// AllowedFormats is the list of formats in which data can be displayed
func (f *ConnListFlags) AllowedFormats() []string {
	formats := f.GenericPrintFlags.AllowedFormats()
	formats = append(formats, f.HumanReadableFlags.AllowedFormats()...)
	return formats
}

// ToPrinter attempts to find a composed set of BrokerListFlags suitable for
// returning a printer based on current flag values.
func (f *ConnListFlags) ToPrinter() (hprinters.ResourcePrinter, error) {
	// if there are flags specified for generic printing
	if f.GenericPrintFlags.OutputFlagSpecified() {
		p, err := f.GenericPrintFlags.ToPrinter()
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	// if no flags specified, use the table printing
	p, err := f.HumanReadableFlags.ToPrinter(ConnectionListHandlers)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// AddFlags receives a *cobra.Command reference and binds
// flags related to humanreadable and template printing
// as well as to reference a service
func (f *ConnListFlags) AddFlags(cmd *cobra.Command) {
	f.GenericPrintFlags.AddFlags(cmd)
	f.HumanReadableFlags.AddFlags(cmd)
	f.FilterFlags.AddFlags(cmd)
}

// NewConnListFlags returns flags associated with humanreadable,
// template, and "name" printing, with default values set.
func NewConnListFlags() *ConnListFlags {
	return &ConnListFlags{
		GenericPrintFlags:  genericclioptions.NewPrintFlags(""),
		HumanReadableFlags: commands.NewHumanPrintFlags(),
	}
}
