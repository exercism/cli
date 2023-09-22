//go:build !windows

package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestBareConfigure(t *testing.T) {
	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupConfigureFlags(flags)

	v := viper.New()
	err := flags.Parse([]string{})
	assert.NoError(t, err)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	err = runConfigure(cfg, flags)
	if assert.Error(t, err) {
		assert.Regexp(t, "no token configured", err.Error())
	}
}

func TestConfigureShow(t *testing.T) {
	co := newCapturedOutput()
	co.newErr = &bytes.Buffer{}
	co.override()
	defer co.reset()

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupConfigureFlags(flags)

	v := viper.New()
	v.Set("token", "configured-token")
	v.Set("workspace", "configured-workspace")
	v.Set("apibaseurl", "http://configured.example.com")

	// it will ignore any flags
	args := []string{
		"--show",
		"--api", "http://override.example.com",
		"--token", "token-override",
		"--workspace", "workspace-override",
	}
	err := flags.Parse(args)
	assert.NoError(t, err)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
	}

	err = runConfigure(cfg, flags)
	assert.NoError(t, err)

	assert.Regexp(t, "configured.example", Err)
	assert.NotRegexp(t, "override.example", Err)

	assert.Regexp(t, "configured-token", Err)
	assert.NotRegexp(t, "token-override", Err)

	assert.Regexp(t, "configured-workspace", Err)
	assert.NotRegexp(t, "workspace-override", Err)
}

func TestConfigureToken(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	testCases := []struct {
		desc       string
		configured string
		args       []string
		expected   string
		message    string
		err        bool
	}{
		{
			desc:       "It doesn't lose a configured value",
			configured: "existing-token",
			args:       []string{"--no-verify"},
			expected:   "existing-token",
		},
		{
			desc:       "It writes a token when passed as a flag",
			configured: "",
			args:       []string{"--no-verify", "--token", "a-token"},
			expected:   "a-token",
		},
		{
			desc:       "It overwrites the token",
			configured: "old-token",
			args:       []string{"--no-verify", "--token", "replacement-token"},
			expected:   "replacement-token",
		},
		{
			desc:       "It complains when token is neither configured nor passed",
			configured: "",
			args:       []string{"--no-verify"},
			expected:   "",
			err:        true,
			message:    "no token configured",
		},
		{
			desc:       "It validates the existing token if we're not skipping validations",
			configured: "configured-token",
			args:       []string{},
			expected:   "configured-token",
			err:        true,
			message:    "token.*invalid",
		},
		{
			desc:       "It validates the replacement token if we're not skipping validations",
			configured: "",
			args:       []string{"--token", "invalid-token"},
			expected:   "",
			err:        true,
			message:    "token.*invalid",
		},
	}

	endpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/validate_token" {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
	ts := httptest.NewServer(endpoint)
	defer ts.Close()

	for _, tc := range testCases {
		flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
		setupConfigureFlags(flags)

		v := viper.New()
		v.Set("token", tc.configured)

		err := flags.Parse(tc.args)
		assert.NoError(t, err)

		cfg := config.Config{
			Persister:       config.InMemoryPersister{},
			UserViperConfig: v,
			DefaultBaseURL:  ts.URL,
		}

		err = runConfigure(cfg, flags)
		if err != nil || tc.err {
			assert.Regexp(t, tc.message, err.Error(), tc.desc)
		}
		assert.Equal(t, tc.expected, cfg.UserViperConfig.GetString("token"), tc.desc)
	}
}

func TestConfigureAPIBaseURL(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	endpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	ts := httptest.NewServer(endpoint)
	defer ts.Close()

	testCases := []struct {
		desc       string
		configured string
		args       []string
		expected   string
		message    string
		err        bool
	}{
		{
			desc:       "It doesn't lose a configured value",
			configured: "http://example.com",
			args:       []string{"--no-verify"},
			expected:   "http://example.com",
		},
		{
			desc:       "It writes a base url when passed as a flag",
			configured: "",
			args:       []string{"--no-verify", "--api", "http://api.example.com"},
			expected:   "http://api.example.com",
		},
		{
			desc:       "It overwrites the base url",
			configured: "http://old.example.com",
			args:       []string{"--no-verify", "--api", "http://replacement.example.com"},
			expected:   "http://replacement.example.com",
		},
		{
			desc:       "It validates the existing base url if we're not skipping validations",
			configured: ts.URL,
			args:       []string{"--token", "some-token"}, // need to bypass the error message on "bare configure"
			expected:   ts.URL,
			err:        true,
			message:    "API.*cannot be reached",
		},
		{
			desc:       "It validates the replacement base URL if we're not skipping validations",
			configured: "",
			args:       []string{"--api", ts.URL},
			expected:   "",
			err:        true,
			message:    "API.*cannot be reached",
		},
	}

	for _, tc := range testCases {
		flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
		setupConfigureFlags(flags)

		v := viper.New()
		v.Set("apibaseurl", tc.configured)

		err := flags.Parse(tc.args)
		assert.NoError(t, err)

		cfg := config.Config{
			Persister:       config.InMemoryPersister{},
			UserViperConfig: v,
			DefaultBaseURL:  ts.URL,
		}

		err = runConfigure(cfg, flags)
		if err != nil || tc.err {
			assert.Regexp(t, tc.message, err.Error(), tc.desc)
		}
		assert.Equal(t, tc.expected, cfg.UserViperConfig.GetString("apibaseurl"), tc.desc)
	}
}

func TestConfigureWorkspace(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	testCases := []struct {
		desc       string
		configured string
		args       []string
		expected   string
		message    string
		err        bool
	}{
		{
			desc:       "It doesn't lose a configured value",
			configured: "/the-workspace",
			args:       []string{"--no-verify"},
			expected:   "/the-workspace",
		},
		{
			desc:       "It writes a workspace when passed as a flag",
			configured: "",
			args:       []string{"--no-verify", "--workspace", "/new-workspace"},
			expected:   "/new-workspace",
		},
		{
			desc:       "It overwrites the configured workspace",
			configured: "/configured-workspace",
			args:       []string{"--no-verify", "--workspace", "/replacement-workspace"},
			expected:   "/replacement-workspace",
		},
		{
			desc:       "It gets the default workspace when neither configured nor passed as a flag",
			configured: "",
			args:       []string{"--token", "some-token"}, // need to bypass the error message on "bare configure"
			expected:   "/home/default-workspace",
		},
		{
			desc:       "It resolves the passed workspace to expand ~",
			configured: "",
			args:       []string{"--workspace", "~/workspace-dir"},
			expected:   "/home/workspace-dir",
		},

		{
			desc:       "It resolves the configured workspace to expand ~",
			configured: "~/configured-dir",
			args:       []string{"--token", "some-token"}, // need to bypass the error message on "bare configure"
			expected:   "/home/configured-dir",            // The configuration object hard-codes the home directory below
		},
	}

	endpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 200 OK by default. Ping and TokenAuth will both pass.
	})
	ts := httptest.NewServer(endpoint)
	defer ts.Close()

	for _, tc := range testCases {
		flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
		setupConfigureFlags(flags)

		v := viper.New()
		v.Set("token", "abc123") // set a token so we get past the no token configured logic
		v.Set("workspace", tc.configured)

		err := flags.Parse(tc.args)
		assert.NoError(t, err)

		cfg := config.Config{
			Persister:       config.InMemoryPersister{},
			UserViperConfig: v,
			DefaultBaseURL:  ts.URL,
			DefaultDirName:  "default-workspace",
			Home:            "/home",
			OS:              "linux",
		}

		err = runConfigure(cfg, flags)
		assert.NoError(t, err, tc.desc)
		assert.Equal(t, tc.expected, cfg.UserViperConfig.GetString("workspace"), tc.desc)
	}
}

func TestConfigureDefaultWorkspaceWithoutClobbering(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	// Stub server to always be 200 OK
	endpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	ts := httptest.NewServer(endpoint)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "no-clobber")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	cfg := config.Config{
		OS:              "linux",
		DefaultDirName:  "workspace",
		Home:            tmpDir,
		Dir:             tmpDir,
		DefaultBaseURL:  ts.URL,
		UserViperConfig: viper.New(),
		Persister:       config.InMemoryPersister{},
	}

	// Create a directory at the workspace directory's location
	// so that it's already present.
	err = os.MkdirAll(config.DefaultWorkspaceDir(cfg), os.FileMode(0755))
	assert.NoError(t, err)

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupConfigureFlags(flags)
	err = flags.Parse([]string{"--token", "abc123"})
	assert.NoError(t, err)

	err = runConfigure(cfg, flags)
	if assert.Error(t, err) {
		assert.Regexp(t, "already something", err.Error())
	}
}

func TestConfigureExplicitWorkspaceWithoutClobberingNonDirectory(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	tmpDir, err := os.MkdirTemp("", "no-clobber")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")

	cfg := config.Config{
		OS:              "linux",
		DefaultDirName:  "workspace",
		Home:            tmpDir,
		Dir:             tmpDir,
		UserViperConfig: v,
		Persister:       config.InMemoryPersister{},
	}

	// Create a file at the workspace directory's location
	err = os.WriteFile(filepath.Join(tmpDir, "workspace"), []byte("This is not a directory"), os.FileMode(0755))
	assert.NoError(t, err)

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupConfigureFlags(flags)
	err = flags.Parse([]string{"--no-verify", "--workspace", config.DefaultWorkspaceDir(cfg)})
	assert.NoError(t, err)

	err = runConfigure(cfg, flags)
	if assert.Error(t, err) {
		assert.Regexp(t, "set a different workspace", err.Error())
	}
}

func TestCommandifyFlagSet(t *testing.T) {
	flags := pflag.NewFlagSet("primitives", pflag.PanicOnError)
	flags.StringP("word", "w", "", "a word")
	flags.BoolP("yes", "y", false, "just do it")
	flags.IntP("number", "n", 1, "count to one")

	err := flags.Parse([]string{"--word", "banana", "--yes"})
	assert.NoError(t, err)
	assert.Equal(t, commandify(flags), "--word=banana --yes=true")
}
