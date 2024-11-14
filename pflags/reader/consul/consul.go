package consul

import (
	"fmt"
	
	"github.com/go-puzzles/puzzles/cores/share"
	"github.com/go-puzzles/puzzles/pflags/reader"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type ConfigReader struct{}

func ConsulReader() *ConfigReader {
	return &ConfigReader{}
}

func (cr *ConfigReader) ReadConfig(v *viper.Viper, opt *reader.Option) error {
	if opt.ServiceName == "" && opt.Tag == "" {
		return errors.New("No service find. ServiceName and Tag is empty.")
	}
	path := fmt.Sprintf("/etc/configs/%v/%v.yaml", opt.ServiceName, opt.Tag)
	v.AddRemoteProvider("consul", share.GetConsulAddr(), path)
	v.SetConfigType("yaml")
	
	if err := v.ReadRemoteConfig(); err != nil {
		return errors.Wrap(err, "readRemoteConfig")
	}
	
	plog.Infof("Read remote config(%v:%v) success. Config=%v", opt.ServiceName, opt.Tag, opt.ConfigPath)
	
	return nil
}
