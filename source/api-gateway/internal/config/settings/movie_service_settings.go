package settings

import "github.com/spf13/viper"

type MovieServiceSettings struct {
	Address string `mapstructure:"address"`
}

func SetMovieServiceDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".address", "localhost:50051")
}
