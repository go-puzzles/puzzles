package reader

import (
	"github.com/spf13/viper"
)

type ReaderOption struct {
	Servicename string
	Tag         string
	ConfigPath  string
}

type ConfigReader interface {
	ReadConfig(v *viper.Viper, opt *ReaderOption) error
}
