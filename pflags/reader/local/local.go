package local

import (
	"os"

	"github.com/go-puzzles/puzzles/pflags/reader"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type localConfigReader struct {
}

func NewLocalConfigReader() *localConfigReader {
	return &localConfigReader{}
}

func (lr *localConfigReader) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (lr *localConfigReader) ReadConfig(v *viper.Viper, opt *reader.Option) error {
	if opt.ConfigPath == "" || !lr.fileExists(opt.ConfigPath) {
		return nil
	}

	v.SetConfigFile(opt.ConfigPath)
	if err := v.ReadInConfig(); err != nil {
		return errors.Wrap(err, "readInConfig")
	}
	plog.Infof("Read local config success. Config=%v", opt.ConfigPath)
	return nil
}
