package consulpuzzle

import (
	"context"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/cores/discover/consul"
	"github.com/go-puzzles/puzzles/cores/share"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/snail"
	"github.com/pkg/errors"

	basepuzzle "github.com/go-puzzles/puzzles/cores/puzzles/base"
	consulReader "github.com/go-puzzles/puzzles/pflags/reader/consul"
)

const (
	consulRegisterHostEnvKey = "GO_PUZZLES_CONSUL_REGISTER_HOST"
)

type consulPuzzle struct {
	*basepuzzle.BasePuzzle
}

var (
	consulAddr      = pflags.String("consulAddr", share.GetConsulAddr(), "Set the conusl addr.")
	useRemoteConfig = pflags.Bool("useRemoteConfig", false, "Whether to use remote configuration")
)

func init() {
	share.ConsulAddr = consulAddr

	snail.RegisterObject("setConsulConfigReader", func() error {
		if !useRemoteConfig() {
			return nil
		}

		pflags.SetDefaultConfigReader(consulReader.ConsulReader())
		return nil
	})
	snail.RegisterObject("setConsulClient", func() error {
		discover.SetFinder(consul.GetConsulClient())
		return nil
	})
}

func WithConsulRegister() cores.ServiceOption {
	return func(o *cores.Options) {
		o.RegisterPuzzle(&consulPuzzle{
			BasePuzzle: &basepuzzle.BasePuzzle{
				PuzzleName: "ConsulRegisterHandler",
			},
		})
	}
}

func (cp *consulPuzzle) StartPuzzle(ctx context.Context, opt *cores.Options) error {
	if opt.ListenerAddr == "" {
		return errors.New("consul register handler can only be used when the service is listening on a port")
	}

	registerAddr := opt.ListenerAddr
	customRegisterHost := os.Getenv(consulRegisterHostEnvKey)
	if customRegisterHost != "" {
		_, port, _ := net.SplitHostPort(opt.ListenerAddr)
		registerAddr = fmt.Sprintf("%s:%s", customRegisterHost, port)
	}

	if opt.ServiceName == "" {
		return errors.New("consul register handler need a ServiceName to register")
	}

	if len(opt.Tags) == 0 {
		opt.Tags = append(opt.Tags, "dev")
	}

	tags := opt.Tags[:]
	sort.Strings(tags)
	if err := discover.GetServiceFinder().RegisterServiceWithTags(opt.ServiceName, registerAddr, tags); err != nil {
		return errors.Wrap(err, "registerConsul")
	}

	var logArgs []any
	logText := "Registered into consul(%s) success. Service=%v Addr=%v"
	logArgs = append(logArgs, consul.GetConsulAddress(), opt.ServiceName, registerAddr)
	if len(tags) > 0 {
		logText = fmt.Sprintf("%v %v", logText, "Tag=%v")
		logArgs = append(logArgs, strings.Join(tags, ","))
	}

	plog.Infoc(ctx, logText, logArgs...)

	<-ctx.Done()

	return nil
}

func (cp *consulPuzzle) Stop() error {
	discover.GetServiceFinder().Close()
	return nil
}
