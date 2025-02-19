package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Model       map[string]string `yaml:"model"`
	RegistryUrl string            `yaml:"registryUrl"`
}

var configPath string

// Initializes the config instance. If the configuration file is not found, it creates a default configuration and writes it to
// the config file. If the configuration file is found but another error occurs, it prints an error message.
func InitConfig(cfg *viper.Viper) {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath = filepath.Join(home, ".modelflux")

	cfg.AddConfigPath(configPath)
	cfg.SetConfigType("yaml")
	cfg.SetConfigName("config")

	// If a config file is found, read it in.
	if err := cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found")
			cfg.Set("model", map[string]string{
				"key":        "",
				"endpoint":   "",
				"deployment": "",
				"version":    "",
			})
			cfg.Set("registryUrl", "http://localhost:3000/api")
			if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
				fmt.Println("Failed to create config directory, ", err)
			}
			if err := cfg.WriteConfigAs(filepath.Join(configPath, "config.yaml")); err != nil {
				fmt.Println("Failed to write default config file, ", err)
			}
		} else {
			fmt.Println("Config file found, but encountered another error")
		}
	}
}

func GetConfigPath() (string, error) {
	if configPath == "" {
		return "", fmt.Errorf("configuration path was not set")
	}
	return configPath, nil
}

func GetWorkflowsPath() (string, error) {
	cfgPath, err := GetConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgPath, "workflows"), nil
}
