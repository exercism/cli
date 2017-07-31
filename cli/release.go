package cli

import "strings"

type Release struct {
	Location string  `json:"html_url"`
	TagName  string  `json:"tag_name"`
	Assets   []Asset `json:"assets"`
}

func (r *Release) Version() string {
	return strings.TrimPrefix(r.TagName, "v")
}
