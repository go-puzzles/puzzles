package cores

import "context"

type Puzzle interface {
	Name() string
	StartPuzzle(context.Context, *Options) error
	Stop() error
}

func (c *CoreService) wrapPuzzles() {
	for _, puzzle := range c.opts.puzzles {
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
