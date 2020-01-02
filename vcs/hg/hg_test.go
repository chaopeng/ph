package hg

import (
	"github.com/chaopeng/ph/config"
	"testing"
)

func Test_enabled(t *testing.T) {
	tests := []struct {
		name string
		conf *config.Config
		want bool
	}{
		{
			name: "not set",
			conf: &config.Config{VCS: map[string]interface{}{}},
			want: false,
		},
		{
			name: "true",
			conf: &config.Config{VCS: map[string]interface{}{
				"hg_enable": true,
			}},
			want: true,
		},
		{
			name: "false",
			conf: &config.Config{VCS: map[string]interface{}{
				"hg_enable": false,
			}},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if enabled(tc.conf) != tc.want {
				t.Errorf("wants %v", tc.want)
			}
		})
	}
}
