// Package git contains vcs info for git
package git

import (
	"fmt"
	"log"
	"os/user"
	"strings"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/vcs"

	"gopkg.in/src-d/go-git.v4"
)

const (
	gitSkip       = "git_skip"
	gitStatusSkip = "git_status_skip"
)

type Git struct {
}

func (s *Git) GetVCSInfo(path string, user *user.User, conf *config.Config) *vcs.VCSInfo {
	return gitInformation(path, user, conf)
}

func vcsSetting(key string, conf *config.Config) []string {
	v, ok := conf.VCS[key]
	if !ok {
		return nil
	}

	s, ok := v.([]interface{})
	if !ok {
		log.Printf("conf.VCS[%s] type incorrect", key)
		return nil
	}

	res := make([]string, len(s))
	for i, v := range s {
		res[i] = fmt.Sprint(v)
	}

	return res
}

// skip gitInformation or gitStatus, git gitInformation is not cheap,
// we may want to avoid this in some case, eg. Chromium.
func skip(key, path string, user *user.User, conf *config.Config) bool {
	for _, exclude := range vcsSetting(key, conf) {
		if strings.HasPrefix(path, exclude) {
			return true
		}
	}

	return false
}

func gitStatus(r *git.Repository) (int, error) {
	w, err := r.Worktree()
	if err != nil {
		return vcs.StatusNone, err
	}

	st, err := w.Status()
	if err != nil {
		return vcs.StatusNone, err
	}

	res := vcs.StatusClean
	for _, s := range st {
		// untrack
		if s.Worktree == git.Untracked {
			res |= vcs.StatusUntrack
		}
		// unstage
		if s.Worktree != s.Staging {
			res |= vcs.StatusUnstage
		}
		// uncommit
		if s.Staging != git.Unmodified {
			res |= vcs.StatusUncommit
		}
	}

	return res, nil
}

func gitInformation(path string, user *user.User, conf *config.Config) (info *vcs.VCSInfo) {
	if skip(gitSkip, path, user, conf) {
		return nil
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		return nil
	}

	// This is a git repo.
	info = &vcs.VCSInfo{}
	info.RepoType = "git"

	ref, err := r.Head()
	if err != nil {
		return
	}

	if ref.Name().IsBranch() {
		info.Name = ref.Name().Short()
	} else {
		info.Name = ref.Hash().String()[0:7]
	}
	info.Status = vcs.StatusNone

	if skip(gitStatusSkip, path, user, conf) {
		return
	}

	status, err := gitStatus(r)
	if err != nil {
		return
	}

	info.Status = status
	return
}
