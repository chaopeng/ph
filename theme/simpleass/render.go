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
	charStatus      = "↪"
	charVCSClean    = "✓"
	charVCSNone     = "*"
	charVCSUntrack  = "?"
	charVCSUnstage  = "!"
	charVCSUncommit = "+"
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

// Render simpleass theme: SSH at OS in PATH on VCS_NAME [VCS_STATUS] \n↪
func (s *Theme) Render(place string, lastStatus string, ctx *context.Context) string {
	if place == theme.Tmux {
		log.Println("simpleass theme does not support tmux")
		return ""
	}

	ctx.ReadCompleteInfo()

	scheme := scheme{
		text:         ctx.Conf.Scheme["simpleass/text"],
		good:         ctx.Conf.Scheme["simpleass/good"],
		bad:          ctx.Conf.Scheme["simpleass/bad"],
		ssh:          ctx.Conf.Scheme["simpleass/ssh"],
		os:           ctx.Conf.Scheme["simpleass/os"],
		prePWD:       ctx.Conf.Scheme["simpleass/pre_pwd"],
		pwd:          ctx.Conf.Scheme["simpleass/pwd"],
		dangerPrePWD: ctx.Conf.Scheme["simpleass/danger_pre_pwd"],
		dangerPWD:    ctx.Conf.Scheme["simpleass/danger_pwd"],
		vcsType:      ctx.Conf.Scheme["simpleass/vcs_type"],
		vcsName:      ctx.Conf.Scheme["simpleass/vcs_name"],
		vcsStatus:    ctx.Conf.Scheme["simpleass/vcs_status"],
	}

	sb := &strings.Builder{}

	// $SSH
	if ctx.Ssh {
		// via
		sb.WriteString(ansi.Color("via ", scheme.text.Fg))
		// ssh
		sb.WriteString(ansi.Color("ssh ", scheme.ssh.Fg))
	}

	// at
	sb.WriteString(ansi.Color("at ", scheme.text.Fg))

	// $OS
	os := ctx.Os
	if ctx.Os == "darwin" {
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
