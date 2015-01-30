package api

import "fmt"

// Problem represents a specific problem in a given language track.
type Problem struct {
	ID        string            `json:"id"`
	TrackID   string            `json:"track_id"`
	Language  string            `json:"language"`
	Slug      string            `json:"slug"`
	Name      string            `json:"name"`
	Files     map[string]string `json:"files"`
	Submitted bool
}

func (p *Problem) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.Language)
}
