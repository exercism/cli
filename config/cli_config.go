package config

import "github.com/spf13/viper"

// CLIConfig contains settings specific to the behavior of the CLI.
type CLIConfig struct {
	*Config
	Tracks Tracks
}

// NewCLIConfig loads the config file in the config directory.
func NewCLIConfig() (*CLIConfig, error) {
	cfg := NewEmptyCLIConfig()

	if err := cfg.Load(viper.New()); err != nil {
		return nil, err
	}
	cfg.SetDefaults()

	return cfg, nil
}

// NewEmptyCLIConfig doesn't load the config from file or set default values.
func NewEmptyCLIConfig() *CLIConfig {
	return &CLIConfig{
		Config: New(Dir(), "cli"),
		Tracks: Tracks{},
	}
}

// Write stores the config to disk.
func (cfg *CLIConfig) Write() error {
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return err
	}
	return Write(cfg)
}

// Validate ensures that the config is valid.
// This is called before writing it.
func (cfg *CLIConfig) Validate() error {
	for _, track := range cfg.Tracks {
		if err := track.CompileRegexes(); err != nil {
			return err
		}
	}
	return nil
}

// SetDefaults ensures that we have all the necessary settings for the CLI.
func (cfg *CLIConfig) SetDefaults() {
	for _, track := range cfg.Tracks {
		track.SetDefaults()
	}
}

// Load reads a viper configuration into the config.
func (cfg *CLIConfig) Load(v *viper.Viper) error {
	cfg.readIn(v)
	return v.Unmarshal(&cfg)
}
