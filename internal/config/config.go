package config

import (
	"errors"
	"strings"

	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/postgres"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Load(args []string) (*Config, error) {
	if len(args) != 2 {
		return nil, errors.New("expected a config path as the only command line argument")
	}
	return readFrom(args[1])
}

func readFrom(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)

	// Set up automatic configuration loading from environment variables of the same name
	// Build tag viper_bind_struct is required to properly unmarshal into a struct
	// TODO: https://github.com/spf13/viper/issues/1797
	v.SetEnvPrefix("muzik")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Apply defaults
	var cfg Config
	cfg.Log.Level = log.LevelInfo
	cfg.Server.Addr = "localhost:80"

	if err := v.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc())); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Config struct {
	Server struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"server"`

	Storage postgres.Config `mapstructure:"storage"`

	Log struct {
		Level log.Level `mapstructure:"level"`
	} `mapstructure:"log"`
}
