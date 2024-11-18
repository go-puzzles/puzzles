package consulpuzzle

import (
	"context"
	"fmt"
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
)

type consulPuzzle struct {
}

var (
	consulAddr = pflags.String("consulAddr", share.GetConsulAddr(), "Set the conusl addr.")
)

func init() {
	share.ConsulAddr = consulAddr
	snail.RegisterObject("setConsulClient", func() error {
		discover.SetFinder(consul.GetConsulClient())
		return nil
	})
}

func WithConsulRegister() cores.ServiceOption {
	return func(o *cores.Options) {
		o.RegisterPuzzle(&consulPuzzle{})
	}
}

func (cp *consulPuzzle) Name() string {
	return "ConsulRegisterHandler"
}

func (cp *consulPuzzle) StartPuzzle(ctx context.Context, opt *cores.Options) error {
	if opt.ListenerAddr == "" {
		return errors.New("consul register handler can only be used when the service is listening on a port")
	}

	if opt.ServiceName == "" {
		return errors.New("consul register handler need a ServiceName to register")
	}

	if len(opt.Tags) == 0 {
		opt.Tags = append(opt.Tags, "dev")
	}

	tags := opt.Tags[:]
	sort.Strings(tags)
	if err := discover.GetServiceFinder().RegisterServiceWithTags(opt.ServiceName, opt.ListenerAddr, tags); err != nil {
		return errors.Wrap(err, "registerConsul")
	}

	var logArgs []any
	logText := "Registered into consul(%s) success. Service=%v"
	logArgs = append(logArgs, consul.GetConsulAddress(), opt.ServiceName)
	if len(tags) > 0 {
		logText = fmt.Sprintf("%v %v", logText, "Tag=%v")
		logArgs = append(logArgs, strings.Join(tags, ","))
	}

	plog.Infoc(ctx, logText, logArgs...)

	<-ctx.Done()

	discover.GetServiceFinder().Close()
	return nil
}

func (cp *consulPuzzle) Stop() error {
	return nil
}
