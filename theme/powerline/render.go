// Package powerline a powerline theme for prompt
package powerline

import (
	"strings"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/context"
	"github.com/chaopeng/ph/theme"
	"github.com/chaopeng/ph/vcs"

	"github.com/mgutz/ansi"
)

const (
	ThemeName = "powerline"

	// Nerf Fonts
	nfLock       = "\uf840"
	nfSlash      = "\ue0bc" // /
	nfBackSlash  = "\ue0b8" // \
	nfRightArrow = "\ue0b0"
)

var (
	defaultScheme = map[string]config.Color{
		"status": config.Color{
			Fg: "1",
			Bg: "15",
		},
		"ssh": config.Color{
			Fg: "252",
			Bg: "240",
		},
		"os": config.Color{
			Fg: "15",
			Bg: "33",
		},
		"pwd": config.Color{
			Fg: "15",
			Bg: "240",
		},
		"pre_pwd": config.Color{
			Fg: "252",
		},
		"danger_zone": config.Color{
			Bg: "124",
		},
		"vcs_status_none": config.Color{
			Fg: "238",
			Bg: "3",
		},
		"vcs_status_clean": config.Color{
			Fg: "238",
			Bg: "2",
		},
		"vcs_status_dirty": config.Color{
			Fg: "15",
			Bg: "1",
		},
	}
)

type Theme struct{}

func (s *Theme) Register() {
	config.RegisterDefaultScheme(ThemeName, defaultScheme)
}

func (s *Theme) Render(place string, lastStatus string, ctx *context.Context) string {
	ctx.ReadCompleteInfo()
	var in input
	if place == theme.Tmux {
		in = input{
			lastStatus: "",
			ctx:        ctx,
			format:     theme.ColorFormatTmux,
			needStatus: false,
			needSSH:    true,
			spliter:    nfBackSlash,
			end:        nfBackSlash,
		}
	}

	if place == theme.Prompt {
		in = input{
			lastStatus: lastStatus,
			ctx:        ctx,
			format:     theme.ColorFormatAnsi,
			needStatus: true,
			needSSH:    false,
			spliter:    nfSlash,
			end:        nfRightArrow,
		}
	}

	return render(in)
}

type input struct {
	lastStatus string
	ctx        *context.Context
	format     theme.ColorFormat
	needStatus bool
	needSSH    bool
	spliter    string
	end        string
}

func (in *input) powerlineRender(bg string, fg string, text string, sb *strings.Builder) {
	if in.format == theme.ColorFormatAnsi {
		if fg == "" {
			sb.WriteString(ansi.Reset)
		} else {
			sb.WriteString(ansi.ColorCode(fg + "+b:" + bg))
		}
	} else if in.format == theme.ColorFormatTmux {
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

func (in *input) powerlineArrowRender(preBg string, bg string, symbol string, sb *strings.Builder) {
	if preBg == "" {
		return
	}

	if in.format == theme.ColorFormatAnsi {
		if bg == "" {
			// We don't want background now. Need reset first.
			sb.WriteString(ansi.Reset)
			sb.WriteString(ansi.ColorCode(preBg))
		} else {
			sb.WriteString(ansi.ColorCode(preBg + ":" + bg))
		}
	} else if in.format == theme.ColorFormatTmux {
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

type scheme struct {
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
func render(in input) string {
	scheme := scheme{
		status:         in.ctx.Conf.GetColor(ThemeName, "status"),
		ssh:            in.ctx.Conf.GetColor(ThemeName, "ssh"),
		os:             in.ctx.Conf.GetColor(ThemeName, "os"),
		pwd:            in.ctx.Conf.GetColor(ThemeName, "pwd"),
		prePWD:         in.ctx.Conf.GetColor(ThemeName, "pre_pwd"),
		dangerZone:     in.ctx.Conf.GetColor(ThemeName, "danger_zone"),
		vcsStatusNone:  in.ctx.Conf.GetColor(ThemeName, "vcs_status_none"),
		vcsStatusClean: in.ctx.Conf.GetColor(ThemeName, "vcs_status_clean"),
		vcsStatusDirty: in.ctx.Conf.GetColor(ThemeName, "vcs_status_dirty"),
	}

	sb := &strings.Builder{}
	bg := ""
	// last status
	if in.needStatus {
		if in.lastStatus != "0" {
			bg = scheme.status.Bg
			in.powerlineRender(bg, scheme.status.Fg, " ! ", sb)
		}
	}

	if in.needSSH && in.ctx.SSH {
		in.powerlineArrowRender(bg, scheme.ssh.Bg, in.spliter, sb)
		bg = scheme.ssh.Bg
		in.powerlineRender(bg, scheme.ssh.Fg, " "+nfLock+" ", sb)
	}

	// OS
	in.powerlineArrowRender(bg, scheme.os.Bg, in.spliter, sb)
	bg = scheme.os.Bg

	if in.ctx.OS == "darwin" {
		in.powerlineRender(bg, scheme.os.Fg, " mac ", sb)
	} else {
		in.powerlineRender(bg, scheme.os.Fg, " "+in.ctx.OS+" ", sb)
	}

	// short pwd
	preBg := bg
	bg = scheme.pwd.Bg
	if in.ctx.PathInfo.DangerZone {
		bg = scheme.dangerZone.Bg
	}
	in.powerlineArrowRender(preBg, bg, in.spliter, sb)
	in.powerlineRender(bg, scheme.prePWD.Fg, " "+in.ctx.PathInfo.ShorternPrefix, sb)
	in.powerlineRender(bg, scheme.pwd.Fg, in.ctx.PathInfo.BaseDir+" ", sb)

	// vcs
	if len(in.ctx.VCSInfo.RepoType) > 0 {
		preBg = bg
		fg := ""
		if in.ctx.VCSInfo.Status == vcs.StatusNone {
			bg = scheme.vcsStatusNone.Bg
			fg = scheme.vcsStatusNone.Fg
		} else if in.ctx.VCSInfo.Status == vcs.StatusClean {
			bg = scheme.vcsStatusClean.Bg
			fg = scheme.vcsStatusClean.Fg
		} else {
			bg = scheme.vcsStatusDirty.Bg
			fg = scheme.vcsStatusDirty.Fg
		}
		in.powerlineArrowRender(preBg, bg, in.spliter, sb)
		in.powerlineRender(bg, fg, " "+in.ctx.VCSInfo.RepoType+":"+in.ctx.VCSInfo.Name, sb)

		if in.ctx.VCSInfo.StatusDirty() {
			in.powerlineRender(bg, fg, "*", sb)
		}
		in.powerlineRender(bg, bg, " ", sb)
	}
	in.powerlineArrowRender(bg, "", in.end, sb)

	in.powerlineRender("", "", " ", sb)
	return sb.String()
}
