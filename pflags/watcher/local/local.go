package local

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/spf13/viper"
)

type LocalWatcher struct {
	callbacks []watchCallback
}

type watchCallback func()

func NewLocalWatcher() *LocalWatcher {
	return &LocalWatcher{
		callbacks: make([]watchCallback, 0),
	}
}

func (lw *LocalWatcher) SetCallbacks(callbacks ...watchCallback) {
	lw.callbacks = append(lw.callbacks, callbacks...)
}

func (lw *LocalWatcher) WatchConfig(v *viper.Viper, path string) {
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		plog.Debugf("local config change")
		for _, callback := range lw.callbacks {
			callback()
		}
	})
}
