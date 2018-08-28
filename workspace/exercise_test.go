package workspace

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasMetadata(t *testing.T) {
	ws, err := ioutil.TempDir("", "fake-workspace")
	defer os.RemoveAll(ws)
	assert.NoError(t, err)

	exerciseA := Exercise{Root: ws, Track: "bogus-track", Slug: "apple"}
	exerciseB := Exercise{Root: ws, Track: "bogus-track", Slug: "banana"}

	err = os.MkdirAll(filepath.Dir(exerciseA.MetadataFilepath()), os.FileMode(0755))
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Dir(exerciseB.MetadataFilepath()), os.FileMode(0755))
	assert.NoError(t, err)

	err = ioutil.WriteFile(exerciseA.MetadataFilepath(), []byte{}, os.FileMode(0600))
	assert.NoError(t, err)

	ok, err := exerciseA.HasMetadata()
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = exerciseB.HasMetadata()
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestNewFromDir(t *testing.T) {
	dir := filepath.Join("something", "another", "whatever", "the-track", "the-exercise")

	exercise := NewExerciseFromDir(dir)
	assert.Equal(t, filepath.Join("something", "another", "whatever"), exercise.Root)
	assert.Equal(t, "the-track", exercise.Track)
	assert.Equal(t, "the-exercise", exercise.Slug)
}

func TestMigrateLegacyMetadataFile(t *testing.T) {
	ws, err := ioutil.TempDir("", "fake-workspace")
	defer os.RemoveAll(ws)
	assert.NoError(t, err)

	exerciseLegacy := Exercise{Root: ws, Track: "bogus-track", Slug: "legacy"}
	exerciseModern := Exercise{Root: ws, Track: "bogus-track", Slug: "modern"}

	err = os.MkdirAll(filepath.Dir(exerciseLegacy.LegacyMetadataFilepath()), os.FileMode(0755))
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Dir(exerciseModern.MetadataFilepath()), os.FileMode(0755))
	assert.NoError(t, err)

	// TODO
}
