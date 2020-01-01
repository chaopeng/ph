package context

import (
	"os/user"
	"testing"

	"github.com/chaopeng/ph/config"
	"github.com/google/go-cmp/cmp"
)

func Test_createPathInfo(t *testing.T) {
	tests := []struct {
		name string
		path string
		want *PathInfo
	}{
		{
			name: "normal",
			path: "/path/a/b",
			want: &PathInfo{
				Orignal:        "/path/a/b",
				ShorternPrefix: "/p/a/",
				BaseDir:        "b",
				DangerZone:     true,
			},
		},
		{
			name: "home",
			path: "/home/user/b",
			want: &PathInfo{
				Orignal:        "/home/user/b",
				ShorternPrefix: "~/",
				BaseDir:        "b",
				DangerZone:     false,
			},
		},
		{
			name: "long path in home",
			path: "/home/user/abcde/f",
			want: &PathInfo{
				Orignal:        "/home/user/abcde/f",
				ShorternPrefix: "~/a/",
				BaseDir:        "f",
				DangerZone:     false,
			},
		},
		{
			name: "deep path in home",
			path: "/home/user/a/b/c/d/e/f",
			want: &PathInfo{
				Orignal:        "/home/user/a/b/c/d/e/f",
				ShorternPrefix: "~/a/.../",
				BaseDir:        "f",
				DangerZone:     false,
			},
		},
		{
			name: "shortern",
			path: "/path/long/b",
			want: &PathInfo{
				Orignal:        "/path/long/b",
				ShorternPrefix: "$short/",
				BaseDir:        "b",
				DangerZone:     false,
			},
		},
	}

	conf := &config.Config{
		PathShorterns: map[string]string{
			"/path/long": "short",
		},
	}

	user := &user.User{
		HomeDir: "/home/user",
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := createPathInfo(tc.path, user, conf)
			if d := cmp.Diff(tc.want, got); len(d) != 0 {
				t.Errorf("createPathInfo(%s) = (-want, + got) %s", tc.path, d)
			}
		})
	}
}
