package cmd

import (
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

type testDefinition struct {
	desc           string
	args           []string
	existingUsrCfg *config.UserConfig
	expectedUsrCfg *config.UserConfig
	existingAPICfg *config.APIConfig
	expectedAPICfg *config.APIConfig
}

func TestConfigure(t *testing.T) {
	tests := []testDefinition{
		testDefinition{
			desc:           "It writes the flags when there is no config file.",
			args:           []string{"fakeapp", "configure", "--token", "a", "--workspace", "/a", "--api", "http://example.com"},
			existingUsrCfg: nil,
			expectedUsrCfg: &config.UserConfig{Token: "a", Workspace: "/a"},
			existingAPICfg: nil,
			expectedAPICfg: &config.APIConfig{BaseURL: "http://example.com"},
		},
		testDefinition{
			desc:           "It overwrites the flags in the config file.",
			args:           []string{"fakeapp", "configure", "--token", "b", "--workspace", "/b", "--api", "http://example.com/v2"},
			existingUsrCfg: &config.UserConfig{Token: "token-b", Workspace: "/workspace-b"},
			expectedUsrCfg: &config.UserConfig{Token: "b", Workspace: "/b"},
			existingAPICfg: &config.APIConfig{BaseURL: "http://example.com/v1"},
			expectedAPICfg: &config.APIConfig{BaseURL: "http://example.com/v2"},
		},
		testDefinition{
			desc:           "It overwrites the flags that are passed, without losing the ones that are not.",
			args:           []string{"fakeapp", "configure", "--token", "c"},
			existingUsrCfg: &config.UserConfig{Token: "token-c", Workspace: "/workspace-c"},
			expectedUsrCfg: &config.UserConfig{Token: "c", Workspace: "/workspace-c"},
		},
		testDefinition{
			desc:           "It gets the default API base URL.",
			args:           []string{"fakeapp", "configure"},
			existingAPICfg: &config.APIConfig{},
			expectedAPICfg: &config.APIConfig{BaseURL: "https://v2.exercism.io/api/v1"},
		},
	}

	for _, definition := range tests {
		t.Run(definition.desc, makeTest(definition))
	}
}

func makeTest(definition testDefinition) func(*testing.T) {

	return func(t *testing.T) {
		var cmdTest *CommandTest
		cmdTest = &CommandTest{
			Cmd:    configureCmd,
			InitFn: initConfigureCmd,
			Args:   definition.args,
		}
		cmdTest.Setup(t)
		defer cmdTest.Teardown(t)

		if definition.existingUsrCfg != nil {
			// Write a fake config.
			cfg := config.NewEmptyUserConfig()
			cfg.Token = definition.existingUsrCfg.Token
			cfg.Workspace = definition.existingUsrCfg.Workspace
			err := cfg.Write()
			assert.NoError(t, err, definition.desc)
		}

		cmdTest.App.Execute()

		if definition.expectedUsrCfg != nil {
			usrCfg, err := config.NewUserConfig()
			assert.NoError(t, err, definition.desc)
			assert.Equal(t, definition.expectedUsrCfg.Token, usrCfg.Token, definition.desc)
			assert.Equal(t, definition.expectedUsrCfg.Workspace, usrCfg.Workspace, definition.desc)
		}

		if definition.expectedAPICfg != nil {
			apiCfg, err := config.NewAPIConfig()
			assert.NoError(t, err, definition.desc)
			assert.Equal(t, definition.expectedAPICfg.BaseURL, apiCfg.BaseURL, definition.desc)
		}
	}
}
