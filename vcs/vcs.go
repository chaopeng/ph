// Package vcs stores helper for getting vcs information
package vcs

import (
	"os/user"

	"github.com/chaopeng/ph/config"
)

const (
	// don't have status or don't care status.
	StatusNone     = -1
	StatusClean    = 0
	StatusUntrack  = 1
	StatusUnstage  = 1 << 1
	StatusUncommit = 1 << 2
)

type VcsInfo struct {
	// git or "", "" means not a repo
	RepoType string
	// p4 client, git clone name
	Workspace string
	// git branch name, git commit id, hg bookmark
	Name string
	// git status, hg status
	Status int
}

func (s *VcsInfo) StatusDirty() bool {
	return s.Status > 0
}

type VCS interface {
	GetVcsInfo(path string, user *user.User, conf *config.Config) *VcsInfo
}

func GetVcsInfo(vs []VCS, path string, user *user.User, conf *config.Config) *VcsInfo {
	for _, v := range vs {
		res := v.GetVcsInfo(path, user, conf)
		if res != nil {
			return res
		}
	}

	return &VcsInfo{}
}
