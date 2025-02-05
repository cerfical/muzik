package config

import (
	"errors"
	"strings"

	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/postgres"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Load(args []string) (*Config, error) {
	v := viper.New()
	if len(args) > 1 {
		if len(args) != 2 {
			return nil, errors.New("expected a config path as the only command line argument")
		}
		v.SetConfigFile(args[1])
	}
	return load(v)
}

func load(v *viper.Viper) (*Config, error) {
	// Set up automatic configuration loading from environment variables of the same name
	// Build tag viper_bind_struct is required to properly unmarshal into a struct
	// TODO: https://github.com/spf13/viper/issues/1797
	v.SetEnvPrefix("muzik")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// Make the configuration file optional
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	v.SetDefault("log.level", log.LevelInfo)
	v.SetDefault("server.addr", "localhost:8080")

	v.SetDefault("db.addr", "localhost:5432")
	v.SetDefault("db.name", "postgres")
	v.SetDefault("db.user", "postgres")

	var cfg Config
	if err := v.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc())); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Config struct {
	Server httpserv.Config
	DB     postgres.Config
	Log    struct {
		Level log.Level
	}
}
