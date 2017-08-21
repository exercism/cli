package config

import (
	"regexp"
	"sort"
)

var defaultIgnorePatterns = []string{
	".*[.]md",
	"[.]solution[.]json",
}

// Track holds the CLI-related settings for a track.
type Track struct {
	ID             string
	IgnorePatterns []string
	ignoreRegexes  []*regexp.Regexp
}

// NewTrack provides a track configured with default values.
func NewTrack(id string) *Track {
	t := &Track{
		ID: id,
	}
	t.SetDefaults()
	return t
}

// SetDefaults configures a track with default values.
func (t *Track) SetDefaults() {
	m := map[string]bool{}
	for _, pattern := range t.IgnorePatterns {
		m[pattern] = true
	}
	for _, pattern := range defaultIgnorePatterns {
		if !m[pattern] {
			t.IgnorePatterns = append(t.IgnorePatterns, pattern)
		}
	}
	sort.Strings(t.IgnorePatterns)
}

// AcceptFilename judges a files admissability based on the name.
func (t *Track) AcceptFilename(f string) (bool, error) {
	if err := t.CompileRegexes(); err != nil {
		return false, err
	}

	for _, re := range t.ignoreRegexes {
		if re.MatchString(f) {
			return false, nil
		}
	}
	return true, nil
}

// CompileRegexes precompiles the ignore patterns.
func (t *Track) CompileRegexes() error {
	if len(t.ignoreRegexes) == len(t.IgnorePatterns) {
		return nil
	}

	for _, pattern := range t.IgnorePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		t.ignoreRegexes = append(t.ignoreRegexes, re)
	}
	return nil
}
