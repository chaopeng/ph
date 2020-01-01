package simpleass

import (
	"testing"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/context"
	"github.com/chaopeng/ph/theme"
	"github.com/chaopeng/ph/vcs"
)

func Test_buildVCSStatusStr(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  string
	}{
		{
			name:  "none",
			input: vcs.StatusNone,
			want:  charVCSNone,
		},
		{
			name:  "clean",
			input: vcs.StatusClean,
			want:  charVCSClean,
		},
		{
			name:  "dirty",
			input: vcs.StatusUncommit | vcs.StatusUnstage | vcs.StatusUntrack,
			want:  charVCSUntrack + charVCSUnstage + charVCSUncommit,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildVCSStatusStr(tc.input)
			want := "[" + tc.want + "]"
			if got != want {
				t.Errorf("buildVCSStatusStr() = %s, wants %s", got, want)
			}
		})
	}
}

func contextWithDefaultConfig() *context.Context {
	return &context.Context{
		Conf: &config.Config{
			Scheme: map[string]config.Color{
				"simpleass/text": config.Color{
					Fg: "15",
				},
				"simpleass/good": config.Color{
					Fg: "15",
				},
				"simpleass/bad": config.Color{
					Fg: "1",
				},
				"simpleass/ssh": config.Color{
					Fg: "11",
				},
				"simpleass/os": config.Color{
					Fg: "130",
				},
				"simpleass/pre_pwd": config.Color{
					Fg: "2",
				},
				"simpleass/pwd": config.Color{
					Fg: "10",
				},
				"simpleass/danger_pre_pwd": config.Color{
					Fg: "161",
				},
				"simpleass/danger_pwd": config.Color{
					Fg: "196",
				},
				"simpleass/vcs_type": config.Color{
					Fg: "6",
				},
				"simpleass/vcs_name": config.Color{
					Fg: "5",
				},
				"simpleass/vcs_status": config.Color{
					Fg: "12",
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
	want       string
}

func testContext(tc testcase) *context.Context {
	ctx := contextWithDefaultConfig()
	ctx.ReadCompleteInfo()
	ctx.Os = "Linux"
	ctx.Ssh = tc.ssh
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
			want: "[38;5;15mat [0m[1;38;5;130mlinux[0m[38;5;15m in [0m[1;38;5;2m/path/a/[0m[1;38;5;10mb[0m[38;5;15m on [0m[1;38;5;6mgit [0m[1;38;5;5mbr-1 [0m[1;38;5;12m[✓][0m\n[1;38;5;15m↪[0m[0m [0m",
		},
		{
			name:       "no vcs info",
			ssh:        false,
			status:     "0",
			pathDanger: false,
			vcsStatus:  -2,
			want: "[38;5;15mat [0m[1;38;5;130mlinux[0m[38;5;15m in [0m[1;38;5;2m/path/a/[0m[1;38;5;10mb[0m\n[1;38;5;15m↪[0m[0m [0m",
		},
		{
			name:       "vcs status not clean",
			ssh:        false,
			status:     "0",
			pathDanger: false,
			vcsStatus:  vcs.StatusUncommit,
			want: "[38;5;15mat [0m[1;38;5;130mlinux[0m[38;5;15m in [0m[1;38;5;2m/path/a/[0m[1;38;5;10mb[0m[38;5;15m on [0m[1;38;5;6mgit [0m[1;38;5;5mbr-1 [0m[1;38;5;12m[+][0m\n[1;38;5;15m↪[0m[0m [0m",
		},
		{
			name:       "vcs status not none",
			ssh:        false,
			status:     "0",
			pathDanger: false,
			vcsStatus:  vcs.StatusNone,
			want: "[38;5;15mat [0m[1;38;5;130mlinux[0m[38;5;15m in [0m[1;38;5;2m/path/a/[0m[1;38;5;10mb[0m[38;5;15m on [0m[1;38;5;6mgit [0m[1;38;5;5mbr-1 [0m[1;38;5;12m[*][0m\n[1;38;5;15m↪[0m[0m [0m",
		},
		{
			name:       "in ssh",
			ssh:        true,
			status:     "0",
			pathDanger: false,
			vcsStatus:  vcs.StatusClean,
			want: "[38;5;15mvia [0m[38;5;11mssh [0m[38;5;15mat [0m[1;38;5;130mlinux[0m[38;5;15m in [0m[1;38;5;2m/path/a/[0m[1;38;5;10mb[0m[38;5;15m on [0m[1;38;5;6mgit [0m[1;38;5;5mbr-1 [0m[1;38;5;12m[✓][0m\n[1;38;5;15m↪[0m[0m [0m",
		},
		{
			name:       "bad status",
			ssh:        false,
			status:     "1",
			pathDanger: false,
			vcsStatus:  vcs.StatusClean,
			want: "[38;5;15mat [0m[1;38;5;130mlinux[0m[38;5;15m in [0m[1;38;5;2m/path/a/[0m[1;38;5;10mb[0m[38;5;15m on [0m[1;38;5;6mgit [0m[1;38;5;5mbr-1 [0m[1;38;5;12m[✓][0m\n[1;38;5;1m↪[0m[0m [0m",
		},
		{
			name:       "danger path",
			ssh:        false,
			status:     "0",
			pathDanger: true,
			vcsStatus:  vcs.StatusClean,
			want: "[38;5;15mat [0m[1;38;5;130mlinux[0m[38;5;15m in [0m[1;38;5;161m/path/a/[0m[1;38;5;196mb[0m[38;5;15m on [0m[1;38;5;6mgit [0m[1;38;5;5mbr-1 [0m[1;38;5;12m[✓][0m\n[1;38;5;15m↪[0m[0m [0m",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testContext(tc)
			p := Theme{}
			got := p.Render(theme.Prompt, tc.status, ctx)
			if got != tc.want {
				t.Errorf("prompt = %s, wants %s", got, tc.want)
			}
			tGot := p.Render(theme.Tmux, tc.status, ctx)
			if len(tGot) != 0 {
				t.Errorf("tmux left = %s, wants empty", tGot)
			}
		})
	}
}
