package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server  string `mapstructure:"server"`
	Email   string `mapstructure:"email"`
	Token   string `mapstructure:"token"`
	Project string `mapstructure:"project"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Use OS-standard config directory:
	//   macOS:   ~/Library/Application Support/jt
	//   Linux:   ~/.config/jt  (XDG_CONFIG_HOME)
	//   Windows: %APPDATA%\jt
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(filepath.Join(configDir, "jt"))

	// Also support legacy ~/.config/jt path (backwards compat on macOS)
	home, _ := os.UserHomeDir()
	if home != "" {
		viper.AddConfigPath(filepath.Join(home, ".config", "jt"))
	}

	// Also allow config in current directory (for development)
	viper.AddConfigPath(".")

	// Environment variable overrides: JT_SERVER, JT_EMAIL, JT_TOKEN, JT_PROJECT
	viper.SetEnvPrefix("JT")
	viper.AutomaticEnv()

	// Also support legacy JIRA_ env vars
	if v := os.Getenv("JIRA_SERVER"); v != "" {
		viper.SetDefault("server", v)
	}
	if v := os.Getenv("JIRA_EMAIL"); v != "" {
		viper.SetDefault("email", v)
	}
	if v := os.Getenv("JIRA_API_TOKEN"); v != "" {
		viper.SetDefault("token", v)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	if cfg.Server == "" || cfg.Email == "" || cfg.Token == "" {
		return nil, fmt.Errorf("not configured — run `jt init` to get started")
	}

	return cfg, nil
}

func ConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", "jt", "config.yml")
	}
	return filepath.Join(configDir, "jt", "config.yml")
}
