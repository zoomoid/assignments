package config

type Configuration struct {
	Status *ConfigurationStatus `json:"status,omitempty"`
	Spec   *ConfigurationSpec   `json:"spec,omitempty"`
}

// ConfigurationSpec are the static configuration fields of an assignments environment
type ConfigurationSpec struct {
	Course        string         `json:"course,omitempty"`
	Group         string         `json:"group,omitempty"`
	Members       []GroupMember  `json:"members,omitempty"`
	BundleOptions *BundleOptions `json:"bundling,omitempty"`
}

// BundleOptions contains configuration for bundling
type BundleOptions struct {
	// Format contains a go template for the filename of a bundle
	Format string `json:"format"`
}

// GroupMembers are part of an assignments group
type GroupMember struct {
	// Name is the group member's full name
	Name string `json:"name"`
	// ID is the group member's student ID or else
	ID string `json:"id"`
}

type InstallationType string

const (
	// InstallationTypeLocal hints to documentclass being a relative path
	InstallationTypeLocal InstallationType = "local"
	// InstallationTypeCTAN indicates the class not being available locally yet, but
	// using the MikTeX mpm subsystem for getting the class from CTAN
	InstallationTypeCTAN InstallationType = "ctan"
	// InstallationTypeSystemWide indicates the csassignment class being available in
	// a local TeXMF tree
	InstallationTypeSystemWide InstallationType = "systemwide"
)

// ConfigurationStatus contains the fields permutated by commands other than the bootstrapping
type ConfigurationStatus struct {
	// Assignment records the current assignment number
	Assignment       uint32           `json:"assignment"`
	InstallationType InstallationType `json:"installationType"`
}
