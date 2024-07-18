package config

import (
	"github.com/spf13/viper"
	"os"
)

// Config defines the structure of the configuration.
type Config struct {
	Nsq struct {
		Address    string `mapstructure:"address" yaml:"address"`
		Nsqd       int    `mapstructure:"nsqd" yaml:"nsqd"`
		Topic      string `mapstructure:"topic" yaml:"topic"`
		Channel    string `mapstructure:"channel" yaml:"channel"`
		Nsqlookupd int    `mapstructure:"nsqlookupd" yaml:"nsqlookupd"`
	}

	Judge struct {
		Address string `mapstructure:"address" yaml:"address"`
	}

	Server struct {
		Port int `mapstructure:"port" yaml:"port"`
	}
}

// CoreConfig is the global configuration variable.
var CoreConfig *Config

// LoadConfig loads the configuration from file.
func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add multiple configuration file paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	// Assign the configuration to the global variable
	CoreConfig = &config

	// Determine whether running in Docker or host environment
	if _, err := os.Stat("/.dockerenv"); err == nil {
		CoreConfig.Nsq.Address = "host.docker.internal"
	}
	return nil
}
