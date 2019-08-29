package theme

import (
	"fmt"
	"time"

	"github.com/chaopeng/ph/context"

	"github.com/mgutz/ansi"
)

func PromptThemeRender(lastStatus string, ctx *context.Context) {
	if ctx.Tmux {
		simpleThemeRender(lastStatus, ctx)
		return
	}
	switch ctx.Conf.Theme {
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

// simple theme: last status >
func simpleThemeRender(lastStatus string, ctx *context.Context) {
	fmt.Printf(ansi.Color(time.Now().Format("2006-01-02 15:04:05 "), ctx.Conf.ColorScheme.Time.Fg+"+b"))
	color := ""
	if lastStatus == "0" {
		color = ansi.ColorCode(ctx.Conf.ColorScheme.StatusGood.Fg + "+b")
	} else {
		color = ansi.ColorCode(ctx.Conf.ColorScheme.StatusBad.Fg + "+b")
	}
	fmt.Printf(color)
	fmt.Printf(">")
	space := ansi.Color(" ", "reset")
	fmt.Printf(space)
}
