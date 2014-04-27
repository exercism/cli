package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchCurrentEndpoint(t *testing.T) {
	expected := "/api/v1/user/assignments/current"
	actual := FetchEndpoint([]string{})
	assert.Equal(t, expected, actual)
}

func TestFetchExerciseEndpoint(t *testing.T) {
	expected := "/api/v1/assignments/language/slug"
	actual := FetchEndpoint([]string{"language", "slug"})
	assert.Equal(t, expected, actual)
}

func TestFetchExerciseEndpointByLanguage(t *testing.T) {
	expected := "/api/v1/assignments/language"
	actual := FetchEndpoint([]string{"language"})
	assert.Equal(t, expected, actual)
}
