package config

import (
	"strings"

	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/config/settings"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP         settings.HTTPSettings         `mapstructure:"http"`
	UserService  settings.UserServiceSettings  `mapstructure:"user_service"`
	MovieService settings.MovieServiceSettings `mapstructure:"movie_service"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/api-gateway")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	settings.SetHTTPDefaults(v, "http")
	settings.SetUserServiceDefaults(v, "user_service")
	settings.SetMovieServiceDefaults(v, "movie_service")
}
