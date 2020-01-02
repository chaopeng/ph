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

var (
	vcses = []VCS{}
)

type VCSInfo struct {
	// git or "", "" means not a repo
	RepoType string
	// git branch name, git commit id, hg bookmark
	Name string
	// git status, hg status
	Status int
}

func (s *VCSInfo) StatusDirty() bool {
	return s.Status > 0
}

func RegisterVCS(v VCS) {
	vcses = append(vcses, v)
}

type VCS interface {
	GetVCSInfo(path string, user *user.User, conf *config.Config) *VCSInfo
}

func GetVCSInfo(path string, user *user.User, conf *config.Config) *VCSInfo {
	for _, v := range vcses {
		res := v.GetVCSInfo(path, user, conf)
		if res != nil {
			return res
		}
	}

	return &VCSInfo{}
}
