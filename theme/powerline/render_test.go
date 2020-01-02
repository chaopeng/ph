package powerline

import (
	"testing"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/context"
	"github.com/chaopeng/ph/theme"
	"github.com/chaopeng/ph/vcs"
)

func contextWithDefaultConfig() *context.Context {
	return &context.Context{
		Conf: &config.Config{
			Scheme: map[string]map[string]config.Color{
				ThemeName: map[string]config.Color{
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
				},
			},
		},
	}
}

type testcase struct {
	name       string
	ssh        bool
	status     string
	pathDanger bool
	vcsStatus  int
	wantPrompt string
	wantTmux   string
}

func testContext(tc testcase) *context.Context {
	ctx := contextWithDefaultConfig()
	ctx.ReadCompleteInfo()
	ctx.OS = "Linux"
	ctx.SSH = tc.ssh
	ctx.Tmux = true
	ctx.PathInfo = &context.PathInfo{
		ShorternPrefix: "/path/a/",
		BaseDir:        "b",
		DangerZone:     tc.pathDanger,
	}
	// no vcs info
	if tc.vcsStatus == -2 {
		ctx.VCSInfo = &vcs.VCSInfo{}
	} else {
		ctx.VCSInfo = &vcs.VCSInfo{
			RepoType: "git",
			Name:     "br-1",
			Status:   tc.vcsStatus,
		}
	}

	return ctx
}

func TestRender(t *testing.T) {
	tests := []testcase{
		{
			name:       "normal",
			ssh:        false,
			status:     "0",
			pathDanger: false,
			vcsStatus:  vcs.StatusClean,
			wantPrompt: "[1;38;5;15;48;5;33m Linux [38;5;33;48;5;240m[1;38;5;252;48;5;240m /path/a/[1;38;5;15;48;5;240mb [38;5;240;48;5;2m[1;38;5;238;48;5;2m git:br-1[1;38;5;2;48;5;2m [0m[38;5;2m[0m ",
			wantTmux: "#[fg=colour15]#[bg=colour33] Linux #[fg=colour33]#[bg=colour240]#[fg=colour252]#[bg=colour240] /path/a/#[fg=colour15]#[bg=colour240]b #[fg=colour240]#[bg=colour2]#[fg=colour238]#[bg=colour2] git:br-1#[fg=colour2]#[bg=colour2] #[fg=colour2]#[bg=default]",
		},
		{
			name:       "no vcs info",
			ssh:        false,
			status:     "0",
			pathDanger: false,
			vcsStatus:  -2,
			wantPrompt: "[1;38;5;15;48;5;33m Linux [38;5;33;48;5;240m[1;38;5;252;48;5;240m /path/a/[1;38;5;15;48;5;240mb [0m[38;5;240m[0m ",
			wantTmux: "#[fg=colour15]#[bg=colour33] Linux #[fg=colour33]#[bg=colour240]#[fg=colour252]#[bg=colour240] /path/a/#[fg=colour15]#[bg=colour240]b #[fg=colour240]#[bg=default]",
		},
		{
			name:       "vcs status not clean",
			ssh:        false,
			status:     "0",
			pathDanger: false,
			vcsStatus:  vcs.StatusUncommit,
			wantPrompt: "[1;38;5;15;48;5;33m Linux [38;5;33;48;5;240m[1;38;5;252;48;5;240m /path/a/[1;38;5;15;48;5;240mb [38;5;240;48;5;1m[1;38;5;15;48;5;1m git:br-1[1;38;5;15;48;5;1m*[1;38;5;1;48;5;1m [0m[38;5;1m[0m ",
			wantTmux: "#[fg=colour15]#[bg=colour33] Linux #[fg=colour33]#[bg=colour240]#[fg=colour252]#[bg=colour240] /path/a/#[fg=colour15]#[bg=colour240]b #[fg=colour240]#[bg=colour1]#[fg=colour15]#[bg=colour1] git:br-1#[fg=colour15]#[bg=colour1]*#[fg=colour1]#[bg=colour1] #[fg=colour1]#[bg=default]",
		},
		{
			name:       "vcs status not none",
			ssh:        false,
			status:     "0",
			pathDanger: false,
			vcsStatus:  vcs.StatusNone,
			wantPrompt: "[1;38;5;15;48;5;33m Linux [38;5;33;48;5;240m[1;38;5;252;48;5;240m /path/a/[1;38;5;15;48;5;240mb [38;5;240;48;5;3m[1;38;5;238;48;5;3m git:br-1[1;38;5;3;48;5;3m [0m[38;5;3m[0m ",
			wantTmux: "#[fg=colour15]#[bg=colour33] Linux #[fg=colour33]#[bg=colour240]#[fg=colour252]#[bg=colour240] /path/a/#[fg=colour15]#[bg=colour240]b #[fg=colour240]#[bg=colour3]#[fg=colour238]#[bg=colour3] git:br-1#[fg=colour3]#[bg=colour3] #[fg=colour3]#[bg=default]",
		},
		{
			name:       "in ssh",
			ssh:        true,
			status:     "0",
			pathDanger: false,
			vcsStatus:  vcs.StatusClean,
			wantPrompt: "[1;38;5;15;48;5;33m Linux [38;5;33;48;5;240m[1;38;5;252;48;5;240m /path/a/[1;38;5;15;48;5;240mb [38;5;240;48;5;2m[1;38;5;238;48;5;2m git:br-1[1;38;5;2;48;5;2m [0m[38;5;2m[0m ",
			wantTmux: "#[fg=colour252]#[bg=colour240]  #[fg=colour240]#[bg=colour33]#[fg=colour15]#[bg=colour33] Linux #[fg=colour33]#[bg=colour240]#[fg=colour252]#[bg=colour240] /path/a/#[fg=colour15]#[bg=colour240]b #[fg=colour240]#[bg=colour2]#[fg=colour238]#[bg=colour2] git:br-1#[fg=colour2]#[bg=colour2] #[fg=colour2]#[bg=default]",
		},
		{
			name:       "bad status",
			ssh:        false,
			status:     "1",
			pathDanger: false,
			vcsStatus:  vcs.StatusClean,
			wantPrompt: "[1;38;5;1;48;5;15m ! [38;5;15;48;5;33m[1;38;5;15;48;5;33m Linux [38;5;33;48;5;240m[1;38;5;252;48;5;240m /path/a/[1;38;5;15;48;5;240mb [38;5;240;48;5;2m[1;38;5;238;48;5;2m git:br-1[1;38;5;2;48;5;2m [0m[38;5;2m[0m ",
			wantTmux: "#[fg=colour15]#[bg=colour33] Linux #[fg=colour33]#[bg=colour240]#[fg=colour252]#[bg=colour240] /path/a/#[fg=colour15]#[bg=colour240]b #[fg=colour240]#[bg=colour2]#[fg=colour238]#[bg=colour2] git:br-1#[fg=colour2]#[bg=colour2] #[fg=colour2]#[bg=default]",
		},
		{
			name:       "danger path",
			ssh:        false,
			status:     "0",
			pathDanger: true,
			vcsStatus:  vcs.StatusClean,
			wantPrompt: "[1;38;5;15;48;5;33m Linux [38;5;33;48;5;124m[1;38;5;252;48;5;124m /path/a/[1;38;5;15;48;5;124mb [38;5;124;48;5;2m[1;38;5;238;48;5;2m git:br-1[1;38;5;2;48;5;2m [0m[38;5;2m[0m ",
			wantTmux: "#[fg=colour15]#[bg=colour33] Linux #[fg=colour33]#[bg=colour124]#[fg=colour252]#[bg=colour124] /path/a/#[fg=colour15]#[bg=colour124]b #[fg=colour124]#[bg=colour2]#[fg=colour238]#[bg=colour2] git:br-1#[fg=colour2]#[bg=colour2] #[fg=colour2]#[bg=default]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testContext(tc)
			p := Theme{}
			pGot := p.Render(theme.Prompt, tc.status, ctx)
			if pGot != tc.wantPrompt {
				t.Errorf("prompt = %s, wants %s", pGot, tc.wantPrompt)
			}
			tGot := p.Render(theme.Tmux, tc.status, ctx)
			if tGot != tc.wantTmux {
				t.Errorf("tmux left = %s, wants %s", tGot, tc.wantTmux)
			}
		})
	}
}
