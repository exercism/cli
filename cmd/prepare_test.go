package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPrepareTrack(t *testing.T) {
	fakeEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := `
		{
			"track": {
				"id": "bogus",
				"language": "Bogus",
				"test_pattern": "_spec[.]ext$"
			}
		}
		`
		fmt.Fprintln(w, payload)
	})
	ts := httptest.NewServer(fakeEndpoint)
	defer ts.Close()

	tmpDir, err := ioutil.TempDir("", "prepare-track")
	assert.NoError(t, err)
	defer os.Remove(tmpDir)

	// Until we can decouple CLIConfig from filesystem, overwrite config dir.
	originalConfigDir := os.Getenv(cfgHomeKey)
	os.Setenv(cfgHomeKey, tmpDir)
	defer os.Setenv(cfgHomeKey, originalConfigDir)

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupPrepareFlags(flags)
	flags.Set("track", "bogus")

	v := viper.New()
	v.Set("apibaseurl", ts.URL)

	cfg := config.Configuration{
		UserViperConfig: v,
	}

	err = runPrepare(cfg, flags, []string{})
	assert.NoError(t, err)

	cliCfg, err := config.NewCLIConfig()
	os.Remove(cliCfg.File())
	assert.NoError(t, err)

	expected := []string{
		".*[.]md",
		"[.]?solution[.]json",
		"_spec[.]ext$",
	}
	track := cliCfg.Tracks["bogus"]
	if track == nil {
		t.Fatal("track missing from config")
	}
	assert.Equal(t, expected, track.IgnorePatterns)
}
