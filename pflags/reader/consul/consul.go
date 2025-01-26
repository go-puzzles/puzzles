package consul

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/go-puzzles/puzzles/cores/discover/consul"
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

func listPossibleTags(tag string) []string {
	v, err := semver.NewVersion(tag)
	if err != nil {
		return []string{tag}
	}

	majorV := v.Major()
	minorV := v.Minor()
	patchV := v.Patch()

	possibleTags := []string{tag}
	for i := int(patchV - 1); i >= 0; i-- {
		possibleTags = append(possibleTags, fmt.Sprintf("v%d.%d.%d", majorV, minorV, i))
	}

	possibleTags = append(possibleTags, fmt.Sprintf("v%d.%d", majorV, minorV))

	return possibleTags
}

func (cr *ConfigReader) getServerKey(name string) string {
	return fmt.Sprintf("/etc/configs/%v", name)
}

func (cr *ConfigReader) getRemotePossiblePath(name, tag string) string {
	var possiblePath string

	kv := consul.GetConsulClient().KV()
	for _, t := range listPossibleTags(tag) {
		path := fmt.Sprintf("%s/%s.yaml", cr.getServerKey(name), t)

		pair, _, err := kv.Get(path, nil)
		if err == nil && pair != nil {
			possiblePath = path
			break
		}
	}

	return possiblePath
}

func (cr *ConfigReader) ReadConfig(v *viper.Viper, opt *reader.Option) error {
	if opt.ServiceName == "" || opt.Tag == "" {
		return errors.New("No service find. ServiceName and Tag is empty.")
	}

	path := cr.getRemotePossiblePath(opt.ServiceName, opt.Tag)
	if path == "" {
		return errors.New("No config found.")
	}

	plog.Infof("Reading consul config from %v", path)
	v.AddRemoteProvider("consul", share.GetConsulAddr(), path)
	v.SetConfigType("yaml")

	if err := v.ReadRemoteConfig(); err != nil {
		return errors.Wrap(err, "readRemoteConfig")
	}

	plog.Infof("Read remote config(%v:%v) success. Config=%v", opt.ServiceName, opt.Tag, opt.ConfigPath)

	return nil
}
