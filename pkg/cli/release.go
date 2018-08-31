package cli

import "strings"

// Release is a specific build of the CLI, released on GitHub.
type Release struct {
	Location string  `json:"html_url"`
	TagName  string  `json:"tag_name"`
	Assets   []Asset `json:"assets"`
}

// Version is the CLI version that is built for the release.
func (r *Release) Version() string {
	return strings.TrimPrefix(r.TagName, "v")
}
