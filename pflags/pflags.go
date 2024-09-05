package pflags

import (
	"os"
	"strings"
	"time"

	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/cores/share"
	"github.com/go-puzzles/puzzles/pflags/reader"
	"github.com/go-puzzles/puzzles/pflags/watcher"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/plog/level"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	localReader "github.com/go-puzzles/puzzles/pflags/reader/local"
	localWatcher "github.com/go-puzzles/puzzles/pflags/watcher/local"
)

var (
	v                 = viper.New()
	requiredFlags     []string
	nestedKey         = map[string]any{}
	defaultConfigFile = "config.yaml"
	serviceName       StringGetter
	debug             BoolGetter
	config            StringGetter
)

type Option struct {
	useConsul     bool
	configReader  reader.ConfigReader
	configWatcher watcher.ConfigWatcher
}

type OptionFunc func(*Option)

func WithConsulEnable() OptionFunc {
	return func(o *Option) {
		o.useConsul = true
	}
}

func WithConfigReader(reader reader.ConfigReader) OptionFunc {
	return func(o *Option) {
		o.configReader = reader
	}
}

func WithConfigWatcher(watchers ...watcher.ConfigWatcher) OptionFunc {
	return func(o *Option) {
		var watcher watcher.ConfigWatcher
		if len(watchers) == 0 {
			watcher = localWatcher.NewLocalWatcher()
		} else {
			watcher = watchers[0]
		}

		lw, ok := watcher.(*localWatcher.LocalWatcher)
		if ok {
			lw.SetCallbacks(structConfReload)
		}

		o.configWatcher = watcher
	}
}

func OverrideDefaultConfigFile(configFile string) {
	defaultConfigFile = configFile
}

func Viper() *viper.Viper {
	return v
}

func initConsulFlag() {
	share.UseConsul = Bool("useConsul", true, "Whether to use the consul service center.")
	share.ConsulAddr = String("consulAddr", discover.GetConsulAddress(), "Set the conusl addr.")
}

func initViper(opt *Option) {
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.SetConfigFile("config.yaml")

	serviceName = StringP("service", "s", os.Getenv("GO_PUZZLE_SERVICE"), "Set the service name")
	config = StringP("config", "f", defaultConfigFile, "Specify config file. Support json, yaml.")
	debug = Bool("debug", false, "Whether to enable debug mode.")

	if opt.useConsul {
		initConsulFlag()
	}

	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		plog.Fatalf("BindPflags error: %v", err)
	}
}

func initOption(opts ...OptionFunc) *Option {
	opt := &Option{}

	for _, o := range opts {
		o(opt)
	}

	if opt.configReader == nil {
		opt.configReader = localReader.NewLocalConfigReader()
	}

	return opt
}

func checkFlagKey() {
	for _, rk := range requiredFlags {
		if isZero(v.Get(rk)) {
			plog.Fatalf("Missing key: %v", rk)
		}
	}
}

func isZero(i interface{}) bool {
	switch i.(type) {
	case bool:
		// It's trivial to check a bool, since it makes the flag no sense(always true).
		return !i.(bool)
	case string:
		return i.(string) == ""
	case time.Duration:
		return i.(time.Duration) == 0
	case float64:
		return i.(float64) == 0
	case int:
		return i.(int) == 0
	case []string:
		return len(i.([]string)) == 0
	case []interface{}:
		return len(i.([]interface{})) == 0
	default:
		return true
	}
}

func readConfig(opt *Option) {
	sp := strings.Split(serviceName(), ":")
	tag := ""
	if len(sp) > 1 {
		tag = sp[len(sp)-1]
	}

	readOpt := &reader.ReaderOption{
		Servicename: sp[0],
		Tag:         tag,
		ConfigPath:  config(),
	}

	// set use consul in flags parse but cancel it in command line
	if opt.useConsul && !BoolGetter(share.UseConsul).Value() {
		opt.configReader = localReader.NewLocalConfigReader()
	}

	if err := opt.configReader.ReadConfig(v, readOpt); err != nil {
		plog.Fatalf(err.Error())
	}
}

func Parse(opts ...OptionFunc) {
	opt := initOption(opts...)
	initViper(opt)
	pflag.Parse()

	readConfig(opt)
	checkFlagKey()

	if debug.Value() {
		plog.Enable(level.LevelDebug)
	}
	if BoolGetter(share.UseConsul).Value() {
		discover.SetConsulFinderToDefault()
	}

	if opt.configWatcher != nil {
		opt.configWatcher.WatchConfig(v, config())
	}
}

func BindPFlag(key string, flag *pflag.Flag) {
	if err := v.BindPFlag(key, flag); err != nil {
		plog.Fatalf("BindPFlag key: %v, err: %v", key, err)
	}
}

func GetServiceName() string {
	return serviceName()
}
