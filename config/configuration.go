package config

import "github.com/spf13/viper"

// Configuration lets us inject configuration options into commands.
// Note that we are slowly working towards getting rid of the
// config.Config, config.UserConfig, and config.CLIConfig types.
// Once we do, we can rename this type to Config, and get rid of the
// User and CLI fields.
type Configuration struct {
	Home                string
	Dir                 string
	DefaultBaseURL      string
	DefaultWorkspaceDir string
	UserViperConfig     *viper.Viper
	UserConfig          *UserConfig
	CLI                 *CLIConfig
}

// NewConfiguration provides a configuration with default values.
func NewConfiguration() Configuration {
	home := userHome()

	return Configuration{
		Dir:                 Dir(),
		Home:                home,
		DefaultBaseURL:      defaultBaseURL,
		DefaultWorkspaceDir: defaultWorkspace(home),
	}
}
