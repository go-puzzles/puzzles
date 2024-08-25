package consulpuzzle

import (
	"context"
	"fmt"
	"sort"
	"strings"
	
	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/cores/share"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
)

type consulPuzzle struct {
}

func WithConsulRegsiter() cores.ServiceOption {
	return func(o *cores.Options) {
		o.RegisterPuzzle(&consulPuzzle{})
	}
}

func (cp *consulPuzzle) Name() string {
	return "ConsulRegsiterHandler"
}

func (cp *consulPuzzle) StartPuzzle(ctx context.Context, opt *cores.Options) error {
	if !share.GetConsulEnable() {
		return errors.New("consul register handler need enable consul first in pflags")
	}
	
	if opt.ListenerAddr == "" {
		return errors.New("consul register handler can only be used when the service is listening on a port")
	}
	
	if opt.ServiceName == "" {
		return errors.New("consul register handler need a ServiceName to registe")
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
	logText := "Registered into consul success. Service=%v"
	logArgs = append(logArgs, opt.ServiceName)
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
