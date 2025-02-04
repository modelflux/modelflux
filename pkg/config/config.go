package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fs = afero.NewOsFs()

type Config struct {
	Model map[string]string `yaml:"model"`
}

// Initializes the config instance. If the configuration file is not found, it creates a default configuration and writes it to
// the config file. If the configuration file is found but another error occurs, it prints an error message.
func InitConfig(cfg *viper.Viper) {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	var configPath = path.Join(home, ".tbd")

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
			if err := fs.MkdirAll(configPath, os.ModePerm); err != nil {
				fmt.Println("Failed to create config directory, ", err)
			}
			if err := cfg.WriteConfigAs(path.Join(configPath, "config.yaml")); err != nil {
				fmt.Println("Failed to write default config file, ", err)
			}
		} else {
			fmt.Println("Config file found, but encountered another error")
		}
	}
}
