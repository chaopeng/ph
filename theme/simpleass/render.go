// Package simpleass contains a prompt theme similar with https://github.com/oh-my-fish/oh-my-fish/blob/master/docs/Themes.md#simple-ass-prompt
package simpleass

import (
	"log"
	"strings"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/context"
	"github.com/chaopeng/ph/theme"
	"github.com/chaopeng/ph/vcs"
	"github.com/mgutz/ansi"
)

const (
	ThemeName = "simpleass"

	charStatus      = "↪"
	charVCSClean    = "✓"
	charVCSNone     = "*"
	charVCSUntrack  = "?"
	charVCSUnstage  = "!"
	charVCSUncommit = "+"
)

var (
	defaultScheme = map[string]config.Color{
		"text": config.Color{
			Fg: "15",
		},
		"good": config.Color{
			Fg: "15",
		},
		"bad": config.Color{
			Fg: "1",
		},
		"ssh": config.Color{
			Fg: "11",
		},
		"os": config.Color{
			Fg: "130",
		},
		"pre_pwd": config.Color{
			Fg: "2",
		},
		"pwd": config.Color{
			Fg: "10",
		},
		"danger_pre_pwd": config.Color{
			Fg: "161",
		},
		"danger_pwd": config.Color{
			Fg: "196",
		},
		"vcs_type": config.Color{
			Fg: "6",
		},
		"vcs_name": config.Color{
			Fg: "5",
		},
		"vcs_status": config.Color{
			Fg: "12",
		},
	}
)

type scheme struct {
	text         config.Color
	good         config.Color
	bad          config.Color
	ssh          config.Color
	os           config.Color
	prePWD       config.Color
	pwd          config.Color
	dangerPrePWD config.Color
	dangerPWD    config.Color
	vcsType      config.Color
	vcsName      config.Color
	vcsStatus    config.Color
}

type Theme struct{}

func (s *Theme) Register() {
	config.RegisterDefaultScheme(ThemeName, defaultScheme)
}

// Render simpleass theme: SSH at OS in PATH on VCS_NAME [VCS_STATUS] \n↪
func (s *Theme) Render(place string, lastStatus string, ctx *context.Context) string {
	if place == theme.Tmux {
		log.Println("simpleass theme does not support tmux")
		return ""
	}

	ctx.ReadCompleteInfo()

	scheme := scheme{
		text:         ctx.Conf.GetColor(ThemeName, "text"),
		good:         ctx.Conf.GetColor(ThemeName, "good"),
		bad:          ctx.Conf.GetColor(ThemeName, "bad"),
		ssh:          ctx.Conf.GetColor(ThemeName, "ssh"),
		os:           ctx.Conf.GetColor(ThemeName, "os"),
		prePWD:       ctx.Conf.GetColor(ThemeName, "pre_pwd"),
		pwd:          ctx.Conf.GetColor(ThemeName, "pwd"),
		dangerPrePWD: ctx.Conf.GetColor(ThemeName, "danger_pre_pwd"),
		dangerPWD:    ctx.Conf.GetColor(ThemeName, "danger_pwd"),
		vcsType:      ctx.Conf.GetColor(ThemeName, "vcs_type"),
		vcsName:      ctx.Conf.GetColor(ThemeName, "vcs_name"),
		vcsStatus:    ctx.Conf.GetColor(ThemeName, "vcs_status"),
	}

	sb := &strings.Builder{}

	// $SSH
	if ctx.SSH {
		// via
		sb.WriteString(ansi.Color("via ", scheme.text.Fg))
		// ssh
		sb.WriteString(ansi.Color("ssh ", scheme.ssh.Fg))
	}

	// at
	sb.WriteString(ansi.Color("at ", scheme.text.Fg))

	// $OS
	os := ctx.OS
	if ctx.OS == "darwin" {
		os = "mac"
	}
	os = strings.ToLower(os)
	sb.WriteString(ansi.Color(os, scheme.os.Fg+"+b"))

	// in
	sb.WriteString(ansi.Color(" in ", scheme.text.Fg))

	// $PWD
	if ctx.PathInfo.DangerZone {
		sb.WriteString(ansi.Color(ctx.PathInfo.ShorternPrefix, scheme.dangerPrePWD.Fg+"+b"))
		sb.WriteString(ansi.Color(ctx.PathInfo.BaseDir, scheme.dangerPWD.Fg+"+b"))
	} else {
		sb.WriteString(ansi.Color(ctx.PathInfo.ShorternPrefix, scheme.prePWD.Fg+"+b"))
		sb.WriteString(ansi.Color(ctx.PathInfo.BaseDir, scheme.pwd.Fg+"+b"))
	}

	// VCS
	if len(ctx.VCSInfo.RepoType) > 0 {
		// on
		sb.WriteString(ansi.Color(" on ", scheme.text.Fg))
		// $VCSType
		sb.WriteString(ansi.Color(ctx.VCSInfo.RepoType+" ", scheme.vcsType.Fg+"+b"))
		// $VCSName
		sb.WriteString(ansi.Color(ctx.VCSInfo.Name+" ", scheme.vcsName.Fg+"+b"))
		// $VCSStatus
		sb.WriteString(ansi.Color(buildVCSStatusStr(ctx.VCSInfo.Status), scheme.vcsStatus.Fg+"+b"))
	}

	// next line
	sb.WriteString("\n")

	// $STATUS
	if lastStatus == "0" {
		sb.WriteString(ansi.Color(charStatus, scheme.good.Fg+"+b"))
	} else {
		sb.WriteString(ansi.Color(charStatus, scheme.bad.Fg+"+b"))
	}

	sb.WriteString(ansi.Color(" ", "reset"))

	return sb.String()
}

func buildVCSStatusStr(status int) string {
	sb := &strings.Builder{}
	sb.WriteString("[")
	if status == vcs.StatusNone {
		sb.WriteString(charVCSNone)
	} else if status == vcs.StatusClean {
		sb.WriteString(charVCSClean)
	} else {
		if status&vcs.StatusUntrack > 0 {
			sb.WriteString(charVCSUntrack)
		}
		if status&vcs.StatusUnstage > 0 {
			sb.WriteString(charVCSUnstage)
		}
		if status&vcs.StatusUncommit > 0 {
			sb.WriteString(charVCSUncommit)
		}
	}

	sb.WriteString("]")
	return sb.String()
}
