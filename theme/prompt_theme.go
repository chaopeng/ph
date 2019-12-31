package theme

import (
	"fmt"
	"time"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/context"

	"github.com/mgutz/ansi"
)

func PromptThemeRender(lastStatus string, ctx *context.Context) {
	if ctx.Tmux {
		simpleThemeRender(lastStatus, ctx)
		return
	}
	switch ctx.Conf.Theme.Prompt {
	case "powerline":
		ctx.ReadCompleteInfo()
		promptPowerlineThemeRender(lastStatus, ctx)
	}
}

func promptPowerlineThemeRender(lastStatus string, ctx *context.Context) {
	render := PowerlineThemeRender{
		LastStatus:  lastStatus,
		Ctx:         ctx,
		ColorFormat: ColorFormatAnsi,
		NeedStatus:  true,
		NeedSsh:     false,
		Spliter:     nfSlash,
		End:         nfRightArrow,
	}
	fmt.Printf(render.Render())
}

type simpleScheme struct {
	time config.Color
	good config.Color
	bad  config.Color
}

// simple theme: last status >
func simpleThemeRender(lastStatus string, ctx *context.Context) {
	scheme := simpleScheme{
		time: ctx.Conf.Scheme["simple/time"],
		good: ctx.Conf.Scheme["simple/good"],
		bad:  ctx.Conf.Scheme["simple/bad"],
	}

	fmt.Printf(ansi.Color(time.Now().Format("2006-01-02 15:04:05 "), scheme.time.Fg+"+b"))
	color := ""
	if lastStatus == "0" {
		color = ansi.ColorCode(scheme.good.Fg + "+b")
	} else {
		color = ansi.ColorCode(scheme.bad.Fg + "+b")
	}
	fmt.Printf(color)
	fmt.Printf(">")
	space := ansi.Color(" ", "reset")
	fmt.Printf(space)
}
