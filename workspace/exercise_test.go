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

func TestHasLegacyMetadata(t *testing.T) {
	ws, err := ioutil.TempDir("", "fake-workspace")
	defer os.RemoveAll(ws)
	assert.NoError(t, err)

	exerciseA := Exercise{Root: ws, Track: "bogus-track", Slug: "apple"}
	exerciseB := Exercise{Root: ws, Track: "bogus-track", Slug: "banana"}

	err = os.MkdirAll(filepath.Dir(exerciseA.LegacyMetadataFilepath()), os.FileMode(0755))
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Dir(exerciseB.LegacyMetadataFilepath()), os.FileMode(0755))
	assert.NoError(t, err)

	err = ioutil.WriteFile(exerciseA.LegacyMetadataFilepath(), []byte{}, os.FileMode(0600))
	assert.NoError(t, err)

	ok, err := exerciseA.HasLegacyMetadata()
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = exerciseB.HasLegacyMetadata()
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

func TestMigrationStatusString(t *testing.T) {
	exercise := Exercise{Root: "", Track: "bogus-track", Slug: "banana"}

	assert.Equal(t, fmt.Sprintf("\nMigrated metadata to %s\n", exercise.MetadataFilepath()),
		MigrationStatusMigrated.String(exercise))
	assert.Equal(t, fmt.Sprintf("\nRemoved legacy metadata at %s\n", exercise.LegacyMetadataFilepath()),
		MigrationStatusRemoved.String(exercise))
	assert.Equal(t, "", MigrationStatusNoop.String(exercise))
	assert.Equal(t, "", MigrationStatus(-1).String(exercise))
}

func TestMigrateLegacyMetadataFileWithoutLegacy(t *testing.T) {
	ws, err := ioutil.TempDir("", "fake-workspace")
	defer os.RemoveAll(ws)
	assert.NoError(t, err)

	exercise := Exercise{Root: ws, Track: "bogus-track", Slug: "no-legacy"}
	metadataFilepath := exercise.MetadataFilepath()
	err = os.MkdirAll(filepath.Dir(metadataFilepath), os.FileMode(0755))
	assert.NoError(t, err)

	err = ioutil.WriteFile(metadataFilepath, []byte{}, os.FileMode(0600))
	assert.NoError(t, err)

	ok, _ := exercise.HasLegacyMetadata()
	assert.False(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)

	status, err := exercise.MigrateLegacyMetadataFile()
	assert.Equal(t, MigrationStatusNoop, status)
	assert.NoError(t, err)

	ok, _ = exercise.HasLegacyMetadata()
	assert.False(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)
}

func TestMigrateLegacyMetadataFileWithLegacy(t *testing.T) {
	ws, err := ioutil.TempDir("", "fake-workspace")
	defer os.RemoveAll(ws)
	assert.NoError(t, err)

	exercise := Exercise{Root: ws, Track: "bogus-track", Slug: "legacy"}
	legacyMetadataFilepath := exercise.LegacyMetadataFilepath()
	err = os.MkdirAll(filepath.Dir(legacyMetadataFilepath), os.FileMode(0755))
	assert.NoError(t, err)

	err = ioutil.WriteFile(legacyMetadataFilepath, []byte{}, os.FileMode(0600))
	assert.NoError(t, err)

	ok, _ := exercise.HasLegacyMetadata()
	assert.True(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.False(t, ok)

	status, err := exercise.MigrateLegacyMetadataFile()
	assert.Equal(t, MigrationStatusMigrated, status)
	assert.NoError(t, err)

	ok, _ = exercise.HasLegacyMetadata()
	assert.False(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)
}

func TestMigrateLegacyMetadataFileWithLegacyAndModern(t *testing.T) {
	ws, err := ioutil.TempDir("", "fake-workspace")
	defer os.RemoveAll(ws)
	assert.NoError(t, err)

	exercise := Exercise{Root: ws, Track: "bogus-track", Slug: "both-legacy-and-modern"}
	metadataFilepath := exercise.MetadataFilepath()
	legacyMetadataFilepath := exercise.LegacyMetadataFilepath()
	err = os.MkdirAll(filepath.Dir(legacyMetadataFilepath), os.FileMode(0755))
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Dir(metadataFilepath), os.FileMode(0755))
	assert.NoError(t, err)

	err = ioutil.WriteFile(legacyMetadataFilepath, []byte{}, os.FileMode(0600))
	assert.NoError(t, err)
	err = ioutil.WriteFile(metadataFilepath, []byte{}, os.FileMode(0600))
	assert.NoError(t, err)

	ok, _ := exercise.HasLegacyMetadata()
	assert.True(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)

	status, err := exercise.MigrateLegacyMetadataFile()
	assert.Equal(t, MigrationStatusRemoved, status)
	assert.NoError(t, err)

	ok, _ = exercise.HasLegacyMetadata()
	assert.False(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)
}
