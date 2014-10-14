package api

import "fmt"

// Track is a collection of problems in a given language.
type Track struct {
	ID       string   `json:"id"`
	Language string   `json:"language"`
	Active   bool     `json:"active"`
	Problems []string `json:"problems"`
}

// Len lists the number of problems a track has.
func (t *Track) Len() int {
	return len(t.Problems)
}

func (t *Track) String() string {
	return fmt.Sprintf("%s (%s)", t.Language, t.ID)
}
