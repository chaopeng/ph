// Package theme Powerline Theme Render
package theme

import (
	"strings"

	"github.com/chaopeng/ph/context"
	"github.com/chaopeng/ph/vcs"

	"github.com/mgutz/ansi"
)

type ColorFormat int

const (
	ColorFormatAnsi ColorFormat = iota
	ColorFormatTmux
)

type PowerlineThemeRender struct {
	LastStatus  string
	Ctx         *context.Context
	ColorFormat ColorFormat
	NeedStatus  bool
	NeedSsh     bool
	Spliter     string
	End         string
}

func (r *PowerlineThemeRender) powerlineRender(bg string, fg string, text string, sb *strings.Builder) {
	if r.ColorFormat == ColorFormatAnsi {
		if fg == "" {
			sb.WriteString(ansi.Reset)
		} else {
			sb.WriteString(ansi.ColorCode(fg + "+b:" + bg))
		}
	} else if r.ColorFormat == ColorFormatTmux {
		if fg == "" {
			return
		} else {
			sb.WriteString("#[fg=colour")
			sb.WriteString(fg)
			sb.WriteString("]#[bg=colour")
			sb.WriteString(bg)
			sb.WriteString("]")
		}
	}

	sb.WriteString(text)
}

func (r *PowerlineThemeRender) powerlineArrowRender(preBg string, bg string, symbol string, sb *strings.Builder) {
	if preBg == "" {
		return
	}

	if r.ColorFormat == ColorFormatAnsi {
		if bg == "" {
			// We don't want background now. Need reset first.
			sb.WriteString(ansi.Reset)
			sb.WriteString(ansi.ColorCode(preBg))
		} else {
			sb.WriteString(ansi.ColorCode(preBg + ":" + bg))
		}
	} else if r.ColorFormat == ColorFormatTmux {
		sb.WriteString("#[fg=colour")
		sb.WriteString(preBg)
		if bg == "" {
			sb.WriteString("]#[bg=default]")
		} else {
			sb.WriteString("]#[bg=colour")
			sb.WriteString(bg)
			sb.WriteString("]")
		}
	}

	sb.WriteString(symbol)
}

// last status > ssh > os > short_pwd > vcs > branch/client > clean >
func (r *PowerlineThemeRender) Render() string {
	sb := &strings.Builder{}
	bg := ""
	// last status
	if r.NeedStatus {
		if r.LastStatus != "0" {
			bg = r.Ctx.Conf.ColorScheme.Status.Bg
			r.powerlineRender(bg, r.Ctx.Conf.ColorScheme.Status.Fg, " ! ", sb)
		}
	}

	if r.NeedSsh && r.Ctx.Ssh {
		r.powerlineArrowRender(bg, r.Ctx.Conf.ColorScheme.Ssh.Bg, r.Spliter, sb)
		bg = r.Ctx.Conf.ColorScheme.Ssh.Bg
		r.powerlineRender(bg, r.Ctx.Conf.ColorScheme.Ssh.Fg, " "+nfLock+" ", sb)
	}

	// OS
	r.powerlineArrowRender(bg, r.Ctx.Conf.ColorScheme.Os.Bg, r.Spliter, sb)
	bg = r.Ctx.Conf.ColorScheme.Os.Bg

	if r.Ctx.Os == "darwin" {
		r.powerlineRender(bg, r.Ctx.Conf.ColorScheme.Os.Fg, " Mac ", sb)
	} else {
		r.powerlineRender(bg, r.Ctx.Conf.ColorScheme.Os.Fg, " "+r.Ctx.Os+" ", sb)
	}

	// short pwd
	preBg := bg
	bg = r.Ctx.Conf.ColorScheme.Pwd.Bg
	if r.Ctx.PathInfo.DangerZone {
		bg = r.Ctx.Conf.ColorScheme.DangerZone.Bg
	}
	r.powerlineArrowRender(preBg, bg, r.Spliter, sb)
	r.powerlineRender(bg, r.Ctx.Conf.ColorScheme.PrePwd.Fg, " "+r.Ctx.PathInfo.ShorternPrefix, sb)
	r.powerlineRender(bg, r.Ctx.Conf.ColorScheme.Pwd.Fg, r.Ctx.PathInfo.BaseDir+" ", sb)

	// vcs
	if r.Ctx.VCSInfo.RepoType != "" {
		preBg = bg
		fg := ""
		if r.Ctx.VCSInfo.Status == vcs.StatusNone {
			bg = r.Ctx.Conf.ColorScheme.StatusNone.Bg
			fg = r.Ctx.Conf.ColorScheme.StatusNone.Fg
		} else if r.Ctx.VCSInfo.Status == vcs.StatusClean {
			bg = r.Ctx.Conf.ColorScheme.StatusClean.Bg
			fg = r.Ctx.Conf.ColorScheme.StatusClean.Fg
		} else {
			bg = r.Ctx.Conf.ColorScheme.StatusDirty.Bg
			fg = r.Ctx.Conf.ColorScheme.StatusDirty.Fg
		}
		r.powerlineArrowRender(preBg, bg, r.Spliter, sb)
		r.powerlineRender(bg, fg, " "+r.Ctx.VCSInfo.RepoType+":"+r.Ctx.VCSInfo.Name, sb)

		if r.Ctx.VCSInfo.StatusDirty() {
			r.powerlineRender(bg, fg, "*", sb)
		}
		r.powerlineRender(bg, bg, " ", sb)
	}
	r.powerlineArrowRender(bg, "", r.End, sb)

	r.powerlineRender("", "", " ", sb)
	return sb.String()
}
