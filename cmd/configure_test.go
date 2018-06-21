package cmd

import (
	"os"
	"runtime"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	desc           string
	args           []string
	existingUsrCfg *config.UserConfig
	expectedUsrCfg *config.UserConfig
	existingAPICfg *config.APIConfig
	expectedAPICfg *config.APIConfig
}

func TestConfigure(t *testing.T) {
	testCases := []testCase{
		testCase{
			desc:           "It writes the flags when there is no config file.",
			args:           []string{"fakeapp", "configure", "--token", "a", "--workspace", "/a", "--api", "http://example.com"},
			existingUsrCfg: nil,
			expectedUsrCfg: &config.UserConfig{Token: "a", Workspace: "/a"},
			existingAPICfg: nil,
			expectedAPICfg: &config.APIConfig{BaseURL: "http://example.com"},
		},
		testCase{
			desc:           "It overwrites the flags in the config file.",
			args:           []string{"fakeapp", "configure", "--token", "b", "--workspace", "/b", "--api", "http://example.com/v2"},
			existingUsrCfg: &config.UserConfig{Token: "token-b", Workspace: "/workspace-b"},
			expectedUsrCfg: &config.UserConfig{Token: "b", Workspace: "/b"},
			existingAPICfg: &config.APIConfig{BaseURL: "http://example.com/v1"},
			expectedAPICfg: &config.APIConfig{BaseURL: "http://example.com/v2"},
		},
		testCase{
			desc:           "It overwrites the flags that are passed, without losing the ones that are not.",
			args:           []string{"fakeapp", "configure", "--token", "c"},
			existingUsrCfg: &config.UserConfig{Token: "token-c", Workspace: "/workspace-c"},
			expectedUsrCfg: &config.UserConfig{Token: "c", Workspace: "/workspace-c"},
		},
		testCase{
			desc:           "It gets the default API base URL.",
			args:           []string{"fakeapp", "configure"},
			existingAPICfg: &config.APIConfig{},
			expectedAPICfg: &config.APIConfig{BaseURL: "https://mentors-beta.exercism.io/api/v1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, makeTest(tc))
	}
}

func makeTest(tc testCase) func(*testing.T) {

	return func(t *testing.T) {
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
			err := cfg.Write()
			assert.NoError(t, err, tc.desc)
		}

		cmdTest.App.Execute()

		if tc.expectedUsrCfg != nil {
			if runtime.GOOS == "windows" {
				tc.expectedUsrCfg.Normalize()
			}

			usrCfg, err := config.NewUserConfig()

			assert.NoError(t, err, tc.desc)
			assert.Equal(t, tc.expectedUsrCfg.Token, usrCfg.Token, tc.desc)
			assert.Equal(t, tc.expectedUsrCfg.Workspace, usrCfg.Workspace, tc.desc)
		}

		if tc.expectedAPICfg != nil {
			apiCfg, err := config.NewAPIConfig()
			assert.NoError(t, err, tc.desc)
			assert.Equal(t, tc.expectedAPICfg.BaseURL, apiCfg.BaseURL, tc.desc)
			os.Remove(apiCfg.File())
		}
	}
}
