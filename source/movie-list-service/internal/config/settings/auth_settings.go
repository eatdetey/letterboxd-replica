package settings

import "github.com/spf13/viper"

type AuthSettings struct {
	AccessSecret string `mapstructure:"access_secret"`
}

func SetAuthDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".access_secret", "access_secret")
}
