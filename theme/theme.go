// Package theme render the tmux status theme and prompt theme
package theme

import (
	"fmt"

	"github.com/chaopeng/ph/context"
)

type ColorFormat int

const (
	ColorFormatAnsi ColorFormat = iota
	ColorFormatTmux
	Tmux   = "tmux"
	Prompt = "prompt"
)

var (
	themes = map[string]ThemeRender{}
)

type ThemeRender interface {
	Render(place string, lastStatus string, ctx *context.Context) string
	Register()
}

func RegisterTheme(name string, t ThemeRender) {
	themes[name] = t
	t.Register()
}

func Render(place string, lastStatus string, ctx *context.Context) {
	theme := ctx.Conf.Theme.Tmux
	if place == Prompt {
		theme = ctx.Conf.Theme.Prompt
		if ctx.Tmux {
			theme = ctx.Conf.Theme.PromptInTmux
		}
	}

	fmt.Printf(themes[theme].Render(place, lastStatus, ctx))
}
