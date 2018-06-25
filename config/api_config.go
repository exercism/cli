package config

import (
	"strings"

	"github.com/spf13/viper"
)

var (
	defaultBaseURL = "https://v2.exercism.io/api/v1"
)

// APIConfig provides API-specific configuration values.
type APIConfig struct {
	*Config
	BaseURL string
}

// NewAPIConfig loads the config file in the config directory.
func NewAPIConfig() (*APIConfig, error) {
	cfg := NewEmptyAPIConfig()

	if err := cfg.Load(viper.New()); err != nil {
		return nil, err
	}

	cfg.SetDefaults()

	return cfg, nil
}

// SetDefaults ensures that we have all the necessary settings for the API.
func (cfg *APIConfig) SetDefaults() {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
}

// NewEmptyAPIConfig doesn't load the config from file or set default values.
func NewEmptyAPIConfig() *APIConfig {
	return &APIConfig{
		Config: New(Dir(), "api"),
	}
}

// Write stores the config to disk.
func (cfg *APIConfig) Write() error {
	cfg.BaseURL = strings.Trim(cfg.BaseURL, "/")
	cfg.SetDefaults()

	return Write(cfg)
}

// Load reads a viper configuration into the config.
func (cfg *APIConfig) Load(v *viper.Viper) error {
	cfg.readIn(v)
	return v.Unmarshal(&cfg)
}
