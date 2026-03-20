package settings

import "github.com/spf13/viper"

type CORSSettings struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

func SetCORSDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".allow_origins", []string{"*"})
	v.SetDefault(prefix+".allow_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	v.SetDefault(prefix+".allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	v.SetDefault(prefix+".expose_headers", []string{})
	v.SetDefault(prefix+".allow_credentials", false)
	v.SetDefault(prefix+".max_age", 600)
}
