package reader

import (
	"github.com/spf13/viper"
)

type Option struct {
	ServiceName string
	Tag         string
	ConfigPath  string
}

type ConfigReader interface {
	ReadConfig(v *viper.Viper, opt *Option) error
}
