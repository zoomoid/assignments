/*
Copyright 2022 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

type Configuration struct {
	Spec   *ConfigurationSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status *ConfigurationStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// ConfigurationSpec are the static configuration fields of an assignments environment
type ConfigurationSpec struct {
	// Course contains the course's name
	Course string `json:"course,omitempty" yaml:"course,omitempty"`
	// Group name of the members
	Group string `json:"group,omitempty" yaml:"group,omitempty"`
	// Includes contains additional includes into the LaTeX source file,
	// relative to the repository's root
	Includes []Include `json:"includes,omitempty" yaml:"includes,omitempty"`
	// Members are the group members
	Members []GroupMember `json:"members,omitempty" yaml:"members,omitempty"`
	// Template allows users to provide their own assignment template
	// deviating from the default LaTeX source template
	Template string `json:"template,omitempty" yaml:"template,omitempty"`
	// GenerateOptions define options configured statically for generating new assignment directories
	GenerateOptions *GenerateOptions `json:"generate,omitempty" yaml:"generate,omitempty"`
	// BuildOptions are user options for the LaTeX build process
	BuildOptions *BuildOptions `json:"build,omitempty" yaml:"build,omitempty"`
	// BundleOptions are user options for bundling
	BundleOptions *BundleOptions `json:"bundle,omitempty" yaml:"bundle,omitempty"`
}

type Include struct {
	// Path defines a relative path for additional files to include in a TeX template
	// They are included as literals in the template, thus should be relative to
	// the assignment TeX file
	Path string `json:"path" yaml:"path"`
}

// BundleOptions contains configuration for bundling
type BundleOptions struct {
	// Format contains a go template for the filename of a bundle
	Template string `json:"template,omitempty" yaml:"template,omitempty"`
	// Pass in arbitrary data for the template as a map
	Data map[string]interface{} `json:"data,omitempty" yaml:"data,omitempty"`
	// Include defines a list of files to include in the bundle, supports globs
	Include []string `json:"include,omitempty" yaml:"include,omitempty"`
}

type GenerateOptions struct {
	// Create defines a list of bare directories to create when generating a new assignment
	Create []string `json:"create" yaml:"create"`
}

type BuildOptions struct {
	// BuildRecipe is the specification of a LaTeX compiler program and its arguments
	BuildRecipe *Recipe `json:"recipe,omitempty" yaml:"recipe,omitempty"`
	// Cleanup defines the two modes of cleanup, either by running latexmk -C or by directly
	// deleting files based on glob patterns
	Cleanup *CleanupOptions `json:"cleanup,omitempty" yaml:"cleanup,omitempty"`
}

type CleanupOptions struct {
	Glob    *CleanupGlobOptions    `json:"glob,omitempty" yaml:"glob,omitempty"`
	Command *CleanupCommandOptions `json:"command,omitempty" yaml:"command,omitempty"`
}

type CleanupGlobOptions struct {
	Recursive bool     `json:"recursive,omitempty" yaml:"recursive,omitempty"`
	Patterns  []string `json:"patterns,omitempty" yaml:"patterns,omitempty"`
}

type CleanupCommandOptions struct {
	Recipe *Recipe `json:"recipe" yaml:"recipe"`
}

// Recipe is a wrapper type for a list of Tools such that we can more easily
// implement cloning on top of recipes
type Recipe []Tool

type Tool struct {
	// Program name of the LaTeX compiler (or proxy) to use
	Command string `json:"command" yaml:"command"`
	// Argument list for the compiler
	Args []string `json:"args" yaml:"args"`
}

// GroupMembers are part of an assignments group
type GroupMember struct {
	// Name is the group member's full name
	Name string `json:"name" yaml:"name"`
	// ID is the group member's student ID or else
	ID string `json:"id" yaml:"id"`
}

// ConfigurationStatus contains the fields permutated by commands other than the bootstrapping
type ConfigurationStatus struct {
	// Assignment records the current assignment number
	Assignment uint32 `json:"assignment,omitempty" yaml:"assignment,omitempty"`
}

func Minimal() *Configuration {
	return &Configuration{
		Spec:   &ConfigurationSpec{},
		Status: &ConfigurationStatus{},
	}
}
