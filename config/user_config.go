package config

import "github.com/spf13/viper"

var (
	defaultBaseURL = "https://v2.exercism.io/api/v1"
)

// UserConfig contains user-specific settings.
type UserConfig struct {
	*Config
	Workspace  string
	Token      string
	Home       string
	APIBaseURL string
}

// NewUserConfig loads a user configuration if it exists.
func NewUserConfig() (*UserConfig, error) {
	cfg := NewEmptyUserConfig()

	if err := cfg.Load(viper.New()); err != nil {
		return nil, err
	}

	return cfg, nil
}

// NewEmptyUserConfig creates a user configuration without loading it.
func NewEmptyUserConfig() *UserConfig {
	return &UserConfig{
		Config: New(Dir(), "user"),
	}
}

// SetDefaults ensures that we have proper values where possible.
func (cfg *UserConfig) SetDefaults() {
	if cfg.Home == "" {
		cfg.Home = userHome()
	}
	if cfg.APIBaseURL == "" {
		cfg.APIBaseURL = defaultBaseURL
	}
	if cfg.Workspace == "" {
		cfg.Workspace = defaultWorkspace(cfg.Home)
	}
}

// Write stores the config to disk.
func (cfg *UserConfig) Write() error {
	cfg.SetDefaults()
	return Write(cfg)
}

// Load reads a viper configuration into the config.
func (cfg *UserConfig) Load(v *viper.Viper) error {
	cfg.readIn(v)
	return v.Unmarshal(&cfg)
}
