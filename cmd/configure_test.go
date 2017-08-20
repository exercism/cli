package cmd

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {

	// Path stuff is complicated and platform-dependent.
	var root string
	if runtime.GOOS == "windows" {
		cfg := config.NewEmptyUserConfig()
		cfg.Normalize()
		root = cfg.Home
	} else {
		root = "/"
	}

	tests := []struct {
		desc           string
		args           []string
		existingUsrCfg *config.UserConfig
		expectedUsrCfg *config.UserConfig
		existingAPICfg *config.APIConfig
		expectedAPICfg *config.APIConfig
	}{
		{
			desc:           "It writes the flags when there is no config file.",
			args:           []string{"fakeapp", "configure", "--token", "a", "--workspace", "/a", "--api", "http://example.com"},
			existingUsrCfg: nil,
			expectedUsrCfg: &config.UserConfig{Token: "a", Workspace: filepath.Join(root, "a")},
			existingAPICfg: nil,
			expectedAPICfg: &config.APIConfig{BaseURL: "http://example.com"},
		},
		{
			desc:           "It overwrites the flags in the config file.",
			args:           []string{"fakeapp", "configure", "--token", "b", "--workspace", "/b", "--api", "http://example.com/v2"},
			existingUsrCfg: &config.UserConfig{Token: "token-b", Workspace: "/workspace-b"},
			expectedUsrCfg: &config.UserConfig{Token: "b", Workspace: filepath.Join(root, "b")},
			existingAPICfg: &config.APIConfig{BaseURL: "http://example.com/v1"},
			expectedAPICfg: &config.APIConfig{BaseURL: "http://example.com/v2"},
		},
		{
			desc:           "It overwrites the flags that are passed, without losing the ones that are not.",
			args:           []string{"fakeapp", "configure", "--token", "c"},
			existingUsrCfg: &config.UserConfig{Token: "token-c", Workspace: "/workspace-c"},
			expectedUsrCfg: &config.UserConfig{Token: "c", Workspace: filepath.Join(root, "workspace-c")},
		},
		{
			desc:           "It gets the default API base URL.",
			args:           []string{"fakeapp", "configure"},
			existingAPICfg: &config.APIConfig{},
			expectedAPICfg: &config.APIConfig{BaseURL: "https://api.exercism.io/v1"},
		},
	}

	for _, test := range tests {
		cmdTest := &CommandTest{
			Cmd:    configureCmd,
			InitFn: initConfigureCmd,
			Args:   test.args,
		}
		cmdTest.Setup(t)
		defer cmdTest.Teardown(t)

		if test.existingUsrCfg != nil {
			// Write a fake config.
			cfg := config.NewEmptyUserConfig()
			cfg.Token = test.existingUsrCfg.Token
			cfg.Workspace = test.existingUsrCfg.Workspace
			err := cfg.Write()
			assert.NoError(t, err, test.desc)
		}
		if test.existingAPICfg != nil {
			// Write a fake config.
			cfg := config.NewEmptyAPIConfig()
			err := cfg.Write()
			assert.NoError(t, err, test.desc)
		}

		cmdTest.App.Execute()

		if test.expectedUsrCfg != nil {
			usrCfg, err := config.NewUserConfig()
			assert.NoError(t, err, test.desc)
			assert.Equal(t, test.expectedUsrCfg.Token, usrCfg.Token, test.desc)
			assert.Equal(t, test.expectedUsrCfg.Workspace, usrCfg.Workspace, test.desc)
		}

		if test.expectedAPICfg != nil {
			apiCfg, err := config.NewAPIConfig()
			assert.NoError(t, err, test.desc)
			assert.Equal(t, test.expectedAPICfg.BaseURL, apiCfg.BaseURL, test.desc)
		}
	}
}
