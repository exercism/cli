package cmd

import (
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
}

func TestConfigure(t *testing.T) {
	testCases := []testCase{
		testCase{
			desc:           "It writes the flags when there is no config file.",
			args:           []string{"fakeapp", "configure", "--token", "a", "--workspace", "/a", "--api", "http://example.com", "--skip-auth"},
			existingUsrCfg: nil,
			expectedUsrCfg: &config.UserConfig{Token: "a", Workspace: "/a", APIBaseURL: "http://example.com"},
		},
		testCase{
			desc:           "It overwrites the flags in the config file.",
			args:           []string{"fakeapp", "configure", "--token", "b", "--workspace", "/b", "--api", "http://example.com/v2", "--skip-auth"},
			existingUsrCfg: &config.UserConfig{Token: "token-b", Workspace: "/workspace-b", APIBaseURL: "http://example.com/v1"},
			expectedUsrCfg: &config.UserConfig{Token: "b", Workspace: "/b", APIBaseURL: "http://example.com/v2"},
		},
		testCase{
			desc:           "It overwrites the flags that are passed, without losing the ones that are not.",
			args:           []string{"fakeapp", "configure", "--token", "c", "--skip-auth"},
			existingUsrCfg: &config.UserConfig{Token: "token-c", Workspace: "/workspace-c", APIBaseURL: "http://example.com"},
			expectedUsrCfg: &config.UserConfig{Token: "c", Workspace: "/workspace-c", APIBaseURL: "http://example.com"},
		},
		testCase{
			desc:           "It gets the default API base url.",
			args:           []string{"fakeapp", "configure", "--skip-auth"},
			existingUsrCfg: &config.UserConfig{Workspace: "/workspace-c"},
			expectedUsrCfg: &config.UserConfig{Workspace: "/workspace-c", APIBaseURL: "https://v2.exercism.io/api/v1"},
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
