package config

// BundlingContext contains all fields passed down to the template renderer for the filename of a bundle
type BundlingContext struct {
	// ID is the assignment id, e.g. "01"
	ID string `json:"id"`
	// Members contains group member information, possibly included in the filename
	Members []GroupMember `json:"members"`
}
