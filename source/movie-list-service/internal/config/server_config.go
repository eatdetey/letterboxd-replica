package config

import (
	"strings"

	"github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/config/settings"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	GRPCServer settings.GRPCServerSettings `mapstructure:"grpc_server"`
	DB         settings.PostgresSettings   `mapstructure:"db"`
	Migrate    settings.MigrateSettings    `mapstructure:"migrate"`
	Shutdown   settings.ShutdownSettings   `mapstructure:"shutdown"`
	Auth       settings.AuthSettings       `mapstructure:"auth"`
}

func LoadServerConfig() (*ServerConfig, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/movie-list-service")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setServerDefaults(v)

	var cfg ServerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setServerDefaults(v *viper.Viper) {
	settings.SetGRPCServerDefaults(v, "grpc_server", ":50051")
	settings.SetPostgresDefaults(v, "db")
	settings.SetMigrateDefaults(v, "migrate")
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetAuthDefaults(v, "auth")
}
