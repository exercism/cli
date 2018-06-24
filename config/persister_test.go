package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFilePersisterSave(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "fake-config")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	fp := FilePersister{
		// Make sure we don't bomb if the dir doesn't exist.
		Dir: filepath.Join(tmpDir, "subdir"),
	}

	v := viper.New()
	v.Set("hello", "world")

	if err = fp.Save(v, TypeAPI); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(fp.Dir, "api.json")
	b, err := ioutil.ReadFile(path)
	assert.NoError(t, err)

	type apiConfig struct {
		Hello string `json:"hello"`
	}
	var cfg apiConfig
	err = json.Unmarshal(b, &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "world", cfg.Hello)
}

func TestFilePersisterLoad(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "fake-config")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Write a JSON config.
	body := `{"hello": "world"}`
	if err := ioutil.WriteFile(filepath.Join(tmpDir, "api.json"), []byte(body), os.FileMode(0600)); err != nil {
		t.Fatal(err)
	}

	// Load it into a viper config.
	fp := FilePersister{
		Dir: tmpDir,
	}
	v, err := fp.Load(TypeAPI)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, v.GetString("hello"), "world")
}

func TestInMemoryPersister(t *testing.T) {
	v1 := viper.New()
	v1.Set("hello", "world")

	imp := NewInMemoryPersister()

	if err := imp.Save(v1, TypeAPI); err != nil {
		t.Fatal(err)
	}

	v2, err := imp.Load(TypeAPI)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "world", v2.GetString("hello"))
}
