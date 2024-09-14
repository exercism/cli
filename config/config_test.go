package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferSiteURL(t *testing.T) {
	testCases := []struct {
		api, url string
	}{
		{"https://api.exercism.io/v1", "https://exercism.org"},
		{"https://v2.exercism.io/api/v1", "https://v2.exercism.io"},
		{"https://mentors-beta.exercism.io/api/v1", "https://mentors-beta.exercism.io"},
		{"http://localhost:3000/api/v1", "http://localhost:3000"},
		{"", "https://exercism.org"},           // use the default
		{"http://whatever", "http://whatever"}, // you're on your own, pal
	}

	for _, tc := range testCases {
		assert.Equal(t, InferSiteURL(tc.api), tc.url)
	}
}
