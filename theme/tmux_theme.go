package theme

import (
	"fmt"
	"github.com/chaopeng/ph/context"
)

func TmuxThemeRender(ctx *context.Context) {
	switch ctx.Conf.Theme {
	case "powerline":
		tmuxPowerlineThemeRender(ctx)
	}
}

func tmuxPowerlineThemeRender(ctx *context.Context) {
	render := PowerlineThemeRender{
		LastStatus:  "",
		Ctx:         ctx,
		ColorFormat: ColorFormatTmux,
		NeedStatus:  false,
		NeedSsh:     true,
		Spliter:     nfBackSlash,
		End:         nfBackSlash,
	}
	fmt.Printf(render.Render())
}
