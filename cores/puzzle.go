package cores

import (
	"context"

	"github.com/go-puzzles/puzzles/plog"
)

type Puzzle interface {
	Name() string
	Before(*Options) error
	StartPuzzle(context.Context, *Options) error
	Stop() error
}

func (c *CoreService) wrapPuzzles() {
	for _, puzzle := range c.opts.puzzles {
		err := puzzle.Before(c.opts)
		if err != nil {
			plog.Fatalf("Before puzzle %s error: %v", puzzle.Name(), err)
		}

		mf := mountFn{
			name:   puzzle.Name(),
			daemon: true,
			fn: func(ctx context.Context) error {
				return puzzle.StartPuzzle(ctx, c.options())
			},
		}

		c.mountFns = append(c.mountFns, mf)
	}
}
