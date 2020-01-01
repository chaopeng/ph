package simple

import (
	"testing"
	"time"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/context"
	"github.com/chaopeng/ph/theme"
)

func TestRender(t *testing.T) {
	ctx := &context.Context{
		Conf: &config.Config{
			Scheme: map[string]config.Color{
				"simple/time": config.Color{
					Fg: "15",
				},
				"simple/good": config.Color{
					Fg: "2",
				},
				"simple/bad": config.Color{
					Fg: "1",
				},
			},
		},
	}
	s := Theme{}
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	got := s.Render(theme.Prompt, "0", ctx)
	want := "[1;38;5;15m" + timeStr + " [0m[1;38;5;2m>[0m [0m"
	if got != want {
		t.Errorf("s.Render(theme.Prompt, 0, ctx) = %s, wants %s", got, want)
	}

	got = s.Render(theme.Prompt, "1", ctx)
	want = "[1;38;5;15m" + timeStr + " [0m[1;38;5;1m>[0m [0m"
	if got != want {
		t.Errorf("s.Render(theme.Prompt, 1, ctx) = %s, wants %s", got, want)
	}

	got = s.Render(theme.Tmux, "1", ctx)
	want = ""
	if got != want {
		t.Errorf("s.Render(theme.Tmux, 1, ctx) = %s, wants %s", got, want)
	}
}
