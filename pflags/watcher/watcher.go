package watcher

import "github.com/spf13/viper"

type ConfigWatcher interface {
	WatchConfig(v *viper.Viper, path string)
}
