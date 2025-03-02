package env

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	envPrefix = "MORU"
)

func init() {
	viper.SetDefault("application.name", "moru")
	viper.SetDefault("application.role", "host")
	viper.SetDefault("application.stage", "dev")
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("discovery.port", 32516)
	viper.SetDefault("http.port", 5623)

	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func LoadViper() (*Config, error) {
	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
