package settings

import "github.com/spf13/viper"

type AuthSettings struct {
	AccessSecret       string `mapstructure:"access_secret"`
	RefreshSecret      string `mapstructure:"refresh_secret"`
	AccessTokenTTLMin  uint   `mapstructure:"access_token_ttl_min"`
	RefreshTokenTTLMin uint   `mapstructure:"refresh_token_ttl_min"`
}

func SetAuthDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".access_secret", "access_secret")
	v.SetDefault(prefix+".refresh_secret", "refresh_secret")
	v.SetDefault(prefix+".access_token_ttl_min", 15)
	v.SetDefault(prefix+".refresh_token_ttl_min", 60*24*7)
}
