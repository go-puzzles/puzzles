package consul

import (
	"fmt"
	
	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/pflags/reader"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type consulConfigReader struct{}

func ConsulReader() *consulConfigReader {
	return &consulConfigReader{}
}

func (cr *consulConfigReader) ReadConfig(v *viper.Viper, opt *reader.ReaderOption) error {
	if opt.Servicename == "" && opt.Tag == "" {
		return errors.New("No service find. ServiceName and Tag is empty.")
	}
	path := fmt.Sprintf("/etc/configs/%v/%v.yaml", opt.Servicename, opt.Tag)
	v.AddRemoteProvider("consul", discover.GetConsulAddress(), path)
	v.SetConfigType("yaml")
	
	if err := v.ReadRemoteConfig(); err != nil {
		return errors.Wrap(err, "readRemoteConfig")
	}
	
	plog.Infof("Read remote config(%v:%v) success. Config=%v", opt.Servicename, opt.Tag, opt.ConfigPath)
	
	return nil
}
