// Package vcs stores helper for getting vcs information
package vcs

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/chaopeng/ph/config"

	"gopkg.in/src-d/go-git.v4"
)

type VcsStatus int

const (
	StatusNone VcsStatus = iota
	StatusClean
	StatusDirty
)

type VcsInfo struct {
	// git or "", "" means not a repo
	Repo string
	// git branch name or git commit id
	Name string
	// git status
	Status VcsStatus
}

func GetVcsInfo(path string, user *user.User, conf *config.Config) *VcsInfo {
	return getGitInformation(path, user, conf)
}

// shouldCallGitInformation getGitInformation is not cheap, we should avoid this in some case.
// eg, chromium src.
func shouldCallGitInformation(path string, user *user.User, conf *config.Config) bool {
	if strings.HasPrefix(path, user.HomeDir) {
		return true
	}
	for _, include := range conf.Includes {
		if strings.HasPrefix(path, include) {
			return true
		}
	}

	return false
}

// getGitInformation returns needed information from git
func getGitInformation(path string, user *user.User, conf *config.Config) (info *VcsInfo) {
	info = &VcsInfo{}
	// find the .git first
	for {
		if !shouldCallGitInformation(path, user, conf) {
			return
		}
		if _, err := os.Stat(path + "/.git"); !os.IsNotExist(err) {
			break
		}
		if path == "/" {
			return
		}

		var err error
		path, err = filepath.Abs(path + "/..")
		if err != nil {
			return
		}
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		return
	}
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
	info.Status = StatusDirty
	if st.IsClean() {
		info.Status = StatusClean
	}

	return
}
