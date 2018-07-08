// +build !windows

package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfigureToken(t *testing.T) {
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

	oldOut := Out
	oldErr := Err
	Out = ioutil.Discard
	defer func() {
		Out = oldOut
		Err = oldErr
	}()

	for _, tc := range testCases {
		var buf bytes.Buffer
		Err = &buf

		flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
		v := viper.New()
		v.Set("token", tc.configured)
		setupConfigureFlags(flags, v)

		err := flags.Parse(tc.args)
		assert.NoError(t, err)

		cfg := config.Configuration{
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
			args:       []string{},
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

	oldOut := Out
	oldErr := Err
	Out = ioutil.Discard
	defer func() {
		Out = oldOut
		Err = oldErr
	}()

	for _, tc := range testCases {
		var buf bytes.Buffer
		Err = &buf

		flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
		v := viper.New()
		v.Set("apibaseurl", tc.configured)
		setupConfigureFlags(flags, v)

		err := flags.Parse(tc.args)
		assert.NoError(t, err)

		cfg := config.Configuration{
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

func TestConfigure(t *testing.T) {
	oldOut := Out
	oldErr := Err
	Out = ioutil.Discard
	Err = ioutil.Discard
	defer func() {
		Out = oldOut
		Err = oldErr
	}()

	type testCase struct {
		desc           string
		args           []string
		existingUsrCfg *config.UserConfig
		expectedUsrCfg *config.UserConfig
	}

	testCases := []testCase{
		testCase{
			desc: "It writes the flags when there is no config file.",
			args: []string{
				"fakeapp", "configure", "--no-verify",
				"--token", "abc123",
				"--workspace", "/workspace",
				"--api", "http://api.example.com",
			},
			existingUsrCfg: nil,
			expectedUsrCfg: &config.UserConfig{Token: "abc123", Workspace: "/workspace", APIBaseURL: "http://api.example.com"},
		},
		testCase{
			desc: "It overwrites the flags in the config file.",
			args: []string{
				"fakeapp", "configure", "--no-verify",
				"--token", "new-token",
				"--workspace", "/new-workspace",
				"--api", "http://new.example.com",
			},
			existingUsrCfg: &config.UserConfig{Token: "old-token", Workspace: "/old-workspace", APIBaseURL: "http://old.example.com"},
			expectedUsrCfg: &config.UserConfig{Token: "new-token", Workspace: "/new-workspace", APIBaseURL: "http://new.example.com"},
		},
		testCase{
			desc: "It overwrites the flags that are passed, without losing the ones that are not.",
			args: []string{
				"fakeapp", "configure", "--no-verify",
				"--token", "replacement-token",
			},
			existingUsrCfg: &config.UserConfig{Token: "original-token", Workspace: "/unmodified", APIBaseURL: "http://unmodified.example.com"},
			expectedUsrCfg: &config.UserConfig{Token: "replacement-token", Workspace: "/unmodified", APIBaseURL: "http://unmodified.example.com"},
		},
		testCase{
			desc:           "It gets the default API base url.",
			args:           []string{"fakeapp", "configure", "--no-verify"},
			existingUsrCfg: &config.UserConfig{Workspace: "/configured-workspace"},
			expectedUsrCfg: &config.UserConfig{Workspace: "/configured-workspace", APIBaseURL: "https://v2.exercism.io/api/v1"},
		},
	}

	for _, tc := range testCases {
		cmdTest := &CommandTest{
			Cmd:    configureCmd,
			InitFn: initConfigureCmd,
			Args:   tc.args,
		}
		cmdTest.Setup(t)
		defer cmdTest.Teardown(t)

		if tc.existingUsrCfg != nil {
			// Write a fake config.
			cfg := config.NewEmptyUserConfig()
			cfg.Token = tc.existingUsrCfg.Token
			cfg.Workspace = tc.existingUsrCfg.Workspace
			cfg.APIBaseURL = tc.existingUsrCfg.APIBaseURL
			err := cfg.Write()
			assert.NoError(t, err, tc.desc)
		}

		cmdTest.App.Execute()

		if tc.expectedUsrCfg != nil {
			if runtime.GOOS == "windows" {
				tc.expectedUsrCfg.SetDefaults()
			}

			cfg, err := config.NewUserConfig()

			assert.NoError(t, err, tc.desc)
			assert.Equal(t, tc.expectedUsrCfg.Token, cfg.Token, tc.desc)
			assert.Equal(t, tc.expectedUsrCfg.Workspace, cfg.Workspace, tc.desc)
			assert.Equal(t, tc.expectedUsrCfg.APIBaseURL, cfg.APIBaseURL, tc.desc)
		}
	}
}
