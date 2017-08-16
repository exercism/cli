package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

func TestPrepareTrack(t *testing.T) {
	cmdTest := &CommandTest{
		Cmd:    prepareCmd,
		InitFn: initPrepareCmd,
		Args:   []string{"fakeapp", "prepare", "--track", "bogus"},
	}
	cmdTest.Setup(t)
	defer cmdTest.Teardown(t)

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

	apiCfg := config.NewEmptyAPIConfig()
	apiCfg.BaseURL = ts.URL
	apiCfg.Endpoints = map[string]string{"prepare-track": "?%s"}
	err := apiCfg.Write()
	assert.NoError(t, err)

	cmdTest.App.Execute()

	cliCfg, err := config.NewCLIConfig()
	assert.NoError(t, err)

	expected := []string{
		".solution.json",
		"README.md",
		"_spec[.]ext$",
	}
	track := cliCfg.Tracks["bogus"]
	if track == nil {
		t.Fatal("track missing from config")
	}
	assert.Equal(t, expected, track.IgnorePatterns)
}
