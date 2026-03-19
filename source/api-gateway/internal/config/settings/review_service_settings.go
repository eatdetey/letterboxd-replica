package settings

import "github.com/spf13/viper"

type ReviewServiceSettings struct {
	Address string `mapstructure:"address"`
}

func SetReviewServiceDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".address", "localhost:50054")
}
