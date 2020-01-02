package context

import (
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/vcs"
)

type PathInfo struct {
	Orignal        string
	ShorternPrefix string
	BaseDir        string
	DangerZone     bool
}

type Context struct {
	OS string
	// is in ssh environment
	SSH bool
	// is in tmux
	Tmux     bool
	pwd      string
	PathInfo *PathInfo
	User     *user.User
	Conf     *config.Config
	VCSInfo  *vcs.VCSInfo
	complete bool
}

func CreateContext() *Context {
	c := &Context{}

	c.Conf = config.ReadConfig()

	if len(os.Getenv("TMUX")) > 0 {
		c.Tmux = true
	}

	return c
}

func createPathInfo(path0 string, user *user.User, conf *config.Config) *PathInfo {
	pathInfo := &PathInfo{Orignal: path0}
	homeDir := user.HomeDir
	path := strings.Replace(path0, homeDir, "~", 1)

	for p, short := range conf.PathShorterns {
		r := regexp.MustCompile(p)
		path = r.ReplaceAllString(path, "@"+short)
	}

	// Consider the path is not user's directory if we don't have shortern setting.
	if path == path0 {
		pathInfo.DangerZone = true
	}

	ss := strings.Split(path, "/")
	prePath := ss[0 : len(ss)-1]
	pathInfo.BaseDir = ss[len(ss)-1]

	if len(prePath) == 0 {
		pathInfo.ShorternPrefix = ""
	} else if len(prePath) > 5 {
		// use ... if path too deep
		pathInfo.ShorternPrefix = prePath[0] + "/" + prePath[1][0:1] + "/.../"
	} else {
		for i := 1; i < len(prePath); i++ {
			prePath[i] = prePath[i][0:1]
		}
		pathInfo.ShorternPrefix = strings.Join(prePath, "/") + "/"
	}
	return pathInfo
}

func (c *Context) ReadPathInfo() {
	var err error

	c.User, err = user.Current()
	if err != nil {
		log.Fatalln("user.Current() failed. ", err)
	}

	// pwd not exists, maybe just call rm, we should still handle this case.
	c.pwd, _ = os.Getwd()

	if c.pwd != "" {
		c.PathInfo = createPathInfo(c.pwd, c.User, c.Conf)
	} else {
		c.PathInfo = &PathInfo{
			BaseDir:    "???",
			DangerZone: true,
		}
	}
}

func (c *Context) ReadCompleteInfo() {
	if c.complete {
		return
	}
	c.complete = true

	c.ReadPathInfo()

	if len(c.pwd) > 0 {
		c.VCSInfo = vcs.GetVCSInfo(c.pwd, c.User, c.Conf)
	} else {
		c.VCSInfo = &vcs.VCSInfo{}
	}

	c.OS = c.Conf.HostName
	if len(c.OS) == 0 {
		c.OS = runtime.GOOS
	}

	if len(os.Getenv("SSH_CLIENT")) > 0 {
		c.SSH = true
	}
}
