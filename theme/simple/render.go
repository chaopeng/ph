// Package simple contains a simple theme for prompt
package simple

import (
	"log"
	"strings"
	"time"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/context"
	"github.com/chaopeng/ph/theme"

	"github.com/mgutz/ansi"
)

const (
	ThemeName = "simple"
)

var (
	defaultScheme = map[string]config.Color{
		"time": config.Color{
			Fg: "15",
		},
		"good": config.Color{
			Fg: "2",
		},
		"bad": config.Color{
			Fg: "1",
		},
	}
)

type Theme struct{}

type scheme struct {
	time config.Color
	good config.Color
	bad  config.Color
}

func (s *Theme) Register() {
	config.RegisterDefaultScheme(ThemeName, defaultScheme)
}

// Render simple theme: timestamp >
func (s *Theme) Render(place string, lastStatus string, ctx *context.Context) string {
	if place == theme.Tmux {
		log.Println("simple theme does not support tmux")
		return ""
	}

	scheme := scheme{
		time: ctx.Conf.GetColor(ThemeName, "time"),
		good: ctx.Conf.GetColor(ThemeName, "good"),
		bad:  ctx.Conf.GetColor(ThemeName, "bad"),
	}

	sb := &strings.Builder{}
	sb.WriteString(ansi.Color(time.Now().Format("2006-01-02 15:04:05 "), scheme.time.Fg+"+b"))
	color := ""
	if lastStatus == "0" {
		color = ansi.ColorCode(scheme.good.Fg + "+b")
	} else {
		color = ansi.ColorCode(scheme.bad.Fg + "+b")
	}
	sb.WriteString(color)
	sb.WriteString(">")
	sb.WriteString(ansi.Color(" ", "reset"))

	return sb.String()
}
