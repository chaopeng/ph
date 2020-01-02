package git

import (
	"testing"

	"github.com/chaopeng/ph/config"
	"github.com/google/go-cmp/cmp"
)

func Test_vcsSetting(t *testing.T) {
	conf := &config.Config{
		VCS: map[string]interface{}{
			gitSkip: []interface{}{"1", "2"},
		},
	}

	got := vcsSetting(gitSkip, conf)
	want := []string{"1", "2"}
	if d := cmp.Diff(want, got); len(d) != 0 {
		t.Errorf("vcsSetting() = (-want +got) %s", d)
	}
}
