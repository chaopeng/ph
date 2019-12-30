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

func (s *Git) GetVcsInfo(path string, user *user.User, conf *config.Config) *vcs.VcsInfo {
	return getGitInformation(path, user, conf)
}

// shouldCallGitInformation getGitInformation is not cheap, we should avoid this in some case.
// eg, chromium src.
func shouldCallGitInformation(path string, user *user.User, conf *config.Config) bool {
	if strings.HasPrefix(path, user.HomeDir) {
		return true
	}
	for _, include := range conf.VCSIncludes {
		if strings.HasPrefix(path, include) {
			return true
		}
	}

	return false
}

// getGitInformation returns needed information from git
func getGitInformation(path string, user *user.User, conf *config.Config) (info *vcs.VcsInfo) {
	if !shouldCallGitInformation(path, user, conf) {
		return nil
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		return nil
	}

	// This is a git repo.
	info = &vcs.VcsInfo{}
	info.Repo = "git"

	ref, err := r.Head()
	if err != nil {
		return
	}

	if ref.Name().IsBranch() {
		info.Name = ref.Name().Short()
	} else {
		info.Name = ref.Hash().String()[0:7]
	}

	w, err := r.Worktree()
	if err != nil {
		return
	}
	st, err := w.Status()
	if err != nil {
		return
	}
	info.Status = vcs.StatusDirty
	if st.IsClean() {
		info.Status = vcs.StatusClean
	}

	return
}
