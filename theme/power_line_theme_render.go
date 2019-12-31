// Package theme Powerline Theme Render
package theme

import (
	"strings"

	"github.com/chaopeng/ph/config"
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

type powerlineScheme struct {
	status         config.Color
	ssh            config.Color
	os             config.Color
	pwd            config.Color
	prePWD         config.Color
	dangerZone     config.Color
	vcsStatusNone  config.Color
	vcsStatusClean config.Color
	vcsStatusDirty config.Color
}

// last status > ssh > os > short_pwd > vcs > branch/client > clean >
func (r *PowerlineThemeRender) Render() string {
	scheme := powerlineScheme{
		status:         r.Ctx.Conf.Scheme["power/status"],
		ssh:            r.Ctx.Conf.Scheme["power/ssh"],
		os:             r.Ctx.Conf.Scheme["power/os"],
		pwd:            r.Ctx.Conf.Scheme["power/pwd"],
		prePWD:         r.Ctx.Conf.Scheme["power/pre_pwd"],
		dangerZone:     r.Ctx.Conf.Scheme["power/danger_zone"],
		vcsStatusNone:  r.Ctx.Conf.Scheme["power/vcs_status_none"],
		vcsStatusClean: r.Ctx.Conf.Scheme["power/vcs_status_clean"],
		vcsStatusDirty: r.Ctx.Conf.Scheme["power/vcs_status_dirty"],
	}

	sb := &strings.Builder{}
	bg := ""
	// last status
	if r.NeedStatus {
		if r.LastStatus != "0" {
			bg = scheme.status.Bg
			r.powerlineRender(bg, scheme.status.Fg, " ! ", sb)
		}
	}

	if r.NeedSsh && r.Ctx.Ssh {
		r.powerlineArrowRender(bg, scheme.ssh.Bg, r.Spliter, sb)
		bg = scheme.ssh.Bg
		r.powerlineRender(bg, scheme.ssh.Fg, " "+nfLock+" ", sb)
	}

	// OS
	r.powerlineArrowRender(bg, scheme.os.Bg, r.Spliter, sb)
	bg = scheme.os.Bg

	if r.Ctx.Os == "darwin" {
		r.powerlineRender(bg, scheme.os.Fg, " Mac ", sb)
	} else {
		r.powerlineRender(bg, scheme.os.Fg, " "+r.Ctx.Os+" ", sb)
	}

	// short pwd
	preBg := bg
	bg = scheme.pwd.Bg
	if r.Ctx.PathInfo.DangerZone {
		bg = scheme.dangerZone.Bg
	}
	r.powerlineArrowRender(preBg, bg, r.Spliter, sb)
	r.powerlineRender(bg, scheme.prePWD.Fg, " "+r.Ctx.PathInfo.ShorternPrefix, sb)
	r.powerlineRender(bg, scheme.pwd.Fg, r.Ctx.PathInfo.BaseDir+" ", sb)

	// vcs
	if r.Ctx.VCSInfo.RepoType != "" {
		preBg = bg
		fg := ""
		if r.Ctx.VCSInfo.Status == vcs.StatusNone {
			bg = scheme.vcsStatusNone.Bg
			fg = scheme.vcsStatusNone.Fg
		} else if r.Ctx.VCSInfo.Status == vcs.StatusClean {
			bg = scheme.vcsStatusClean.Bg
			fg = scheme.vcsStatusClean.Fg
		} else {
			bg = scheme.vcsStatusDirty.Bg
			fg = scheme.vcsStatusDirty.Fg
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
