package api

import "fmt"

// Track is a collection of exercises in a given language.
type Track struct {
	ID        string   `json:"id"`
	Language  string   `json:"language"`
	Active    bool     `json:"active"`
	Exercises []string `json:"problems"`
}

// Len lists the number of exercises a track has.
func (t *Track) Len() int {
	return len(t.Exercises)
}

func (t *Track) String() string {
	return fmt.Sprintf("%s (%s)", t.Language, t.ID)
}
