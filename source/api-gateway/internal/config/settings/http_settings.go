package settings

import "github.com/spf13/viper"

type HTTPSettings struct {
	Addr string `mapstructure:"addr"`
}

func SetHTTPDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".addr", ":8080")
}
