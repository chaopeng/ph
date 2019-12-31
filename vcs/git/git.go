// Package git contains vcs info for git
package git

import (
	"os/user"
	"strings"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/vcs"

	"gopkg.in/src-d/go-git.v4"
)

type Git struct {
}

func (s *Git) GetVCSInfo(path string, user *user.User, conf *config.Config) *vcs.VCSInfo {
	return gitInformation(path, user, conf)
}

// skip gitInformation is not cheap, we may want to avoid this in some case.
// eg, chromium src.
func skip(path string, user *user.User, conf *config.Config) bool {
	for _, exclude := range conf.VCSSkip {
		if strings.HasPrefix(path, exclude) {
			return true
		}
	}

	if strings.HasPrefix(path, user.HomeDir) {
		return false
	}

	for _, include := range conf.VCSIncludes {
		if strings.HasPrefix(path, include) {
			return false
		}
	}

	return true
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
	if skip(path, user, conf) {
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

	status, err := gitStatus(r)
	if err != nil {
		return
	}

	info.Status = status
	return
}
