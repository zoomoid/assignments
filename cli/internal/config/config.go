package config

type Configuration struct {
	Status *ConfigurationStatus `json:"status,omitempty"`
	Spec   *ConfigurationSpec   `json:"spec,omitempty"`
}

// ConfigurationSpec are the static configuration fields of an assignments environment
type ConfigurationSpec struct {
	// Course contains the course's name
	Course string `json:"course,omitempty"`
	// Group name of the members
	Group string `json:"group,omitempty"`
	// Includes contains additional includes into the LaTeX source file,
	// relative to the repository's root
	Includes []Include `json:"includes,omitempty"`
	// Members are the group members
	Members []GroupMember `json:"members,omitempty"`
	// Template allows users to provide their own assignment template
	// deviating from the default LaTeX source template
	Template string `json:"template,omitempty"`
	// GenerateOptions define options configured statically for generating new assignment directories
	GenerateOptions *GenerateOptions `json:"generate,omitempty"`
	// BuildOptions are user options for the LaTeX build process
	BuildOptions *BuildOptions `json:"build,omitempty"`
	// BundleOptions are user options for bundling
	BundleOptions *BundleOptions `json:"bundle,omitempty"`
}

type Include struct {
	// Path defines a relative path for additional files to include in a TeX template
	// They are included as literals in the template, thus should be relative to
	// the assignment TeX file
	Path string `json:"path"`
}

// BundleOptions contains configuration for bundling
type BundleOptions struct {
	// Format contains a go template for the filename of a bundle
	Template string `json:"template"`
	// Pass in arbitrary data for the template as a map
	Data map[string]interface{} `json:"data"`
	// Include defines a list of files to include in the bundle, supports globs
	Include []string `json:"include,omitempty"`
}

type GenerateOptions struct {
	// Create defines a list of bare directories to create when generating a new assignment
	Create []string `json:"create,omitempty"`
}

func (g *GenerateOptions) Clone() *GenerateOptions {
	o := []string{}
	o = append(o, g.Create...)

	return &GenerateOptions{
		Create: o,
	}
}

type BuildOptions struct {
	// Recipe is the specification of a LaTeX compiler program and its arguments
	Recipe []Recipe `json:"recipe,omitempty"`
}

type Recipe struct {
	// Program name of the LaTeX compiler (or proxy) to use
	Command string `json:"command"`
	// Argument list for the compiler
	Args []string `json:"args"`
}

// GroupMembers are part of an assignments group
type GroupMember struct {
	// Name is the group member's full name
	Name string `json:"name"`
	// ID is the group member's student ID or else
	ID string `json:"id"`
}

// ConfigurationStatus contains the fields permutated by commands other than the bootstrapping
type ConfigurationStatus struct {
	// Assignment records the current assignment number
	Assignment uint32 `json:"assignment"`
}
