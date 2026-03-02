package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

const (
	AppName         = "plane-cli"
	ConfigFileName  = "config"
	KeyringService  = "plane-cli"
	KeyringUser     = "api-key"
	DefaultAPIHost  = "https://api.plane.so"
)

type Config struct {
	Version          string `mapstructure:"version"`
	DefaultWorkspace string `mapstructure:"default_workspace"`
	DefaultProject   string `mapstructure:"default_project"`
	OutputFormat     string `mapstructure:"output_format"`
	APIHost          string `mapstructure:"api_host"`
}

var (
	cfgFile string
	Cfg     Config
)

func InitConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		configDir := filepath.Join(home, ".config", AppName)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigName(ConfigFileName)
		viper.SetConfigType("yaml")
	}

	viper.SetDefault("version", "1.0")
	viper.SetDefault("output_format", "table")
	viper.SetDefault("api_host", DefaultAPIHost)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func SaveConfig() error {
	viper.Set("version", Cfg.Version)
	viper.Set("default_workspace", Cfg.DefaultWorkspace)
	viper.Set("default_project", Cfg.DefaultProject)
	viper.Set("output_format", Cfg.OutputFormat)
	viper.Set("api_host", Cfg.APIHost)

	return viper.WriteConfig()
}

func GetAPIKey() (string, error) {
	return keyring.Get(KeyringService, KeyringUser)
}

func SetAPIKey(apiKey string) error {
	return keyring.Set(KeyringService, KeyringUser, apiKey)
}

func DeleteAPIKey() error {
	return keyring.Delete(KeyringService, KeyringUser)
}

func SetConfigFile(file string) {
	cfgFile = file
}
