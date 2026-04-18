package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for vaultpull.
type Config struct {
	VaultAddr  string   `mapstructure:"vault_addr"`
	VaultToken string   `mapstructure:"vault_token"`
	Namespace  string   `mapstructure:"namespace"`
	SecretPath string   `mapstructure:"secret_path"`
	OutputFile string   `mapstructure:"output_file"`
	FilterKeys []string `mapstructure:"filter_keys"`
}

// Load reads configuration from environment variables and an optional config file.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	v.SetDefault("vault_addr", "http://127.0.0.1:8200")
	v.SetDefault("output_file", ".env")

	v.SetEnvPrefix("VAULTPULL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Also respect native Vault env vars.
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		v.SetDefault("vault_addr", addr)
	}
	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		v.SetDefault("vault_token", token)
	}

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return errors.New("vault_addr must not be empty")
	}
	if c.VaultToken == "" {
		return errors.New("vault_token is required (set VAULT_TOKEN or VAULTPULL_VAULT_TOKEN)")
	}
	if c.SecretPath == "" {
		return errors.New("secret_path is required")
	}
	return nil
}
