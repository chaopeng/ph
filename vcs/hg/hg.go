// Package hg contains vcs info for hg
package hg

import (
	"log"
	"os/exec"
	"os/user"
	"strings"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/vcs"
)

type HG struct {
}

func commandExec(command string, args ...string) (string, error) {
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// Status returns hg status to vcs status, error if it is not a hg repo.
func Status() (int, error) {
	// hg status, it is not a hg repo, if error.
	status, err := commandExec("hg", "status")
	if err != nil {
		return vcs.StatusNone, nil
	}

	st := vcs.StatusClean
	if len(status) == 0 {
		return st, nil
	}

	for _, s := range strings.Split(status, "\n") {
		if len(s) == 0 {
			continue
		}
		code := s[0:1]
		if code == "C" {
			continue
		}
		if code == "?" {
			st |= vcs.StatusUntrack
		} else {
			// hg seems do not have unstage state
			st |= vcs.StatusUncommit
		}
	}

	return st, nil
}

// CurrentBookmark returns current bookmark, returns "" if not in bookmark.
func CurrentBookmark() (string, error) {
	// hg bookmark
	bookmarks, err := commandExec("hg", "bookmark")
	if err != nil {
		return "", err
	}

	if len(bookmarks) == 0 {
		return "", nil
	}
	for _, s := range strings.Split(bookmarks, "\n") {
		if len(s) == 0 {
			continue
		}
		s = strings.TrimLeft(s, " ")
		if s[0:1] == "*" {
			return strings.SplitN(s[2:], " ", 2)[0], nil
		}
	}
	return "", nil
}

func (s *HG) GetVCSInfo(path string, user *user.User, conf *config.Config) *vcs.VCSInfo {
	st, err := Status()
	if err != nil {
		return nil
	}

	v := &vcs.VCSInfo{
		RepoType: "hg",
		Status:   st,
	}

	n, err := CurrentBookmark()
	if err != nil {
		log.Printf("CurrentBookmark failed: %v", err)
		return v
	}

	if len(n) == 0 {
		n = "default"
	}
	v.Name = n

	return v
}
