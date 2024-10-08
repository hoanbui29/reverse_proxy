package config

import (
	"github.com/spf13/viper"
)

type Strategy int

const (
	StrategyRoundRobin Strategy = iota
)

type Config struct {
	Server    ConfigServer `mapstructure:"server"`
	LogLevel  string       `mapstructure:"log_level"`
	Resources []Resource   `mapstructure:"resources"`
}

type ConfigServer struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Resource struct {
	Timeout      int      `mapstructure:"timeout"`
	Prefix       string   `mapstructure:"prefix"`
	Strategy     Strategy `mapstructure:"strategy"`
	Destinations []string `mapstructure:"destinations"`
	Methods      []string `mapstructure:"methods"`
}

func Load() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	// viper.AutomaticEnv()
	// viper.SetEnvPrefix("REVERSE_PROXY")
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("log_level", "DEBUG")
	err := viper.ReadInConfig()

	cfg := Config{}

	if err != nil {
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)

	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
