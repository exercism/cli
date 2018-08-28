package workspace

import (
	"fmt"
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
	var str string
	ws, err := ioutil.TempDir("", "fake-workspace")
	defer os.RemoveAll(ws)
	assert.NoError(t, err)

	exercise := Exercise{Root: ws, Track: "bogus-track", Slug: "migration"}
	metadataFilepath := exercise.MetadataFilepath()
	legacyMetadataFilepath := exercise.LegacyMetadataFilepath()

	err = os.MkdirAll(filepath.Dir(legacyMetadataFilepath), os.FileMode(0755))
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Dir(metadataFilepath), os.FileMode(0755))
	assert.NoError(t, err)

	// returns nil if not legacy
	err = ioutil.WriteFile(metadataFilepath, []byte{}, os.FileMode(0600))
	assert.NoError(t, err)
	ok, _ := exercise.HasMetadata()
	assert.True(t, ok)
	_, err = exercise.MigrateLegacyMetadataFile()
	assert.Nil(t, err)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)

	// legacy metadata only => gets renamed
	os.Remove(metadataFilepath)
	err = ioutil.WriteFile(legacyMetadataFilepath, []byte{}, os.FileMode(0600))
	assert.NoError(t, err)
	ok, _ = exercise.HasLegacyMetadata()
	assert.True(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.False(t, ok)
	str, err = exercise.MigrateLegacyMetadataFile()
	assert.Equal(t, fmt.Sprintf("\nMigrated metadata to %s\n", metadataFilepath), str)
	assert.NoError(t, err)
	ok, _ = exercise.HasLegacyMetadata()
	assert.False(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)

	// both legacy and modern metadata files exist => legacy gets deleted
	err = ioutil.WriteFile(legacyMetadataFilepath, []byte{}, os.FileMode(0600))
	assert.NoError(t, err)
	err = ioutil.WriteFile(metadataFilepath, []byte{}, os.FileMode(0600))
	assert.NoError(t, err)
	ok, _ = exercise.HasLegacyMetadata()
	assert.True(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)
	str, err = exercise.MigrateLegacyMetadataFile()
	assert.Equal(t, fmt.Sprintf("\nRemoved legacy metadata: %s\n", legacyMetadataFilepath), str)
	assert.NoError(t, err)
	ok, _ = exercise.HasLegacyMetadata()
	assert.False(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)

}
