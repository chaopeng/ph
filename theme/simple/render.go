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

type SimpleTheme struct{}

type simpleScheme struct {
	time config.Color
	good config.Color
	bad  config.Color
}

// Render simple theme: last status >
func (s *SimpleTheme) Render(place string, lastStatus string, ctx *context.Context) string {
	if place == theme.Tmux {
		log.Println("simple theme does not support tmux")
		return ""
	}

	scheme := simpleScheme{
		time: ctx.Conf.Scheme["simple/time"],
		good: ctx.Conf.Scheme["simple/good"],
		bad:  ctx.Conf.Scheme["simple/bad"],
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
