// Package vcs stores helper for getting vcs information
package vcs

import (
	"os/user"

	"github.com/chaopeng/ph/config"
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
