package settings

import "github.com/spf13/viper"

type UserServiceSettings struct {
	Address string `mapstructure:"address"`
}

func SetUserServiceDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".address", "localhost:50051")
}
