package config

import (
	"regexp"
	"sort"
)

var defaultIgnorePatterns = []string{
	".solution.json",
	"README.md",
}

// Track holds the CLI-related settings for a track.
type Track struct {
	ID             string
	IgnorePatterns []string
	ignoreRegexes  []*regexp.Regexp
}

func NewTrack(id string) *Track {
	t := &Track{
		ID: id,
	}
	t.SetDefaults()
	return t
}

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
