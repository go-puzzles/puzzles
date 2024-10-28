package cores

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
)

func (c *CoreService) GracefulStopPuzzle() {
	for name, puzzle := range c.opts.puzzles {
		if err := puzzle.Stop(); err != nil {
			plog.Errorc(c.ctx, "stop %v error: %v", name, err)
		}
	}
}

func (c *CoreService) gracefulKill() mountFn {
	return mountFn{
		name: "GracefulKill",
		fn: func(ctx context.Context) error {
			ch := make(chan os.Signal, 1)
			signal.Notify(
				ch,
				os.Interrupt,
				syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT,
			)

			select {
			case sg := <-ch:
				plog.Infoc(c.ctx, "Graceful stopping puzzles...")
				c.GracefulStopPuzzle()

				if c.opts.Cmux != nil {
					c.opts.Cmux.Close()
				}

				if c.listener != nil {
					c.listener.Close()
				}

				c.cancel()

				plog.Infoc(c.ctx, "Graceful stopped puzzles successfully")
				return errors.Errorf("Signal: %s", sg.String())
			case <-ctx.Done():
				return ctx.Err()
			}
		},
		daemon: true,
	}
}
