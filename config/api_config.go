package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// APIConfig provides API-specific configuration values.
type APIConfig struct {
	*Config
	BaseURL   string
	Endpoints map[string]string
}

// NewAPIConfig loads the config file in the config directory.
func NewAPIConfig() (*APIConfig, error) {
	cfg := NewEmptyAPIConfig()

	// Set defaults.
	cfg.BaseURL = "https://api.exercism.com/v1"
	cfg.Endpoints = map[string]string{
		"download": "/solutions/%s",
		"submit":   "/solutions/%s",
	}

	if err := cfg.Load(viper.New()); err != nil {
		return nil, err
	}

	return cfg, nil
}

// URL provides the API URL for a given endpoint key.
func (cfg *APIConfig) URL(key string) string {
	return fmt.Sprintf("%s%s", cfg.BaseURL, cfg.Endpoints[key])
}

// NewEmptyAPIConfig doesn't load the config from file or set default values.
func NewEmptyAPIConfig() *APIConfig {
	return &APIConfig{
		Config: New(Dir(), "api"),
	}
}

// Write stores the config to disk.
func (cfg *APIConfig) Write() error {
	return Write(cfg)
}

// Load reads a viper configuration into the config.
func (cfg *APIConfig) Load(v *viper.Viper) error {
	cfg.readIn(v)
	return v.Unmarshal(&cfg)
}
