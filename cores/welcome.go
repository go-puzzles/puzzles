package cores

import (
	"embed"
	"strings"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/lukesampson/figlet/figletlib"
)

//go:embed fonts/standard.flf
var fontStandard embed.FS

func (c *CoreService) showServiceName() {
	standardFont, err := fontStandard.ReadFile("fonts/standard.flf")
	if err != nil {
		return
	}
	f, err := figletlib.ReadFontFromBytes(standardFont)
	if err != nil {
		plog.Debugc(c.ctx, "Can not show service name because of: %v", err)
		return
	}

	figletlib.PrintMsg(c.opts.ServiceName, f, 80, f.Settings(), "left")
}

func (c *CoreService) welcome() {
	if c.listener != nil {
		plog.Infoc(c.ctx, "Listening... Addr=%v", c.listener.Addr().String())
	}

	if c.opts.ServiceName != "" {
		c.showServiceName()
	}

	if len(c.opts.Tags) != 0 {
		plog.Infoc(c.ctx, "Service Tag=%v", strings.Join(c.opts.Tags, ","))
	}

	plog.Infoc(c.ctx, "Go-Puzzles Service Started. Version=%v", version)
}
