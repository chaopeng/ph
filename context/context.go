package context

import (
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/chaopeng/ph/config"
	"github.com/chaopeng/ph/vcs"
)

type PathInfo struct {
	ShorternPrefix string
	BaseDir        string
	DangerZone     bool
}

type Context struct {
	Os string
	// is in ssh environment
	Ssh bool
	// is in tmux
	Tmux     bool
	pwd      string
	PathInfo *PathInfo
	User     *user.User
	Conf     *config.Config
	VcsInfo  *vcs.VcsInfo
}

func CreateContext() *Context {
	c := &Context{}

	c.Conf = config.ReadConfig()

	if os.Getenv("TMUX") != "" {
		c.Tmux = true
	}

	return c
}

func createPathInfo(path0 string, user *user.User, conf *config.Config) *PathInfo {
	pathInfo := &PathInfo{}
	homeDir := user.HomeDir
	path := strings.Replace(path0, homeDir, "~", 1)

	for p, short := range conf.PathShorterns {
		path = strings.Replace(path, p, "$"+short, 1)
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
	c.ReadPathInfo()

	if c.pwd != "" {
		c.VcsInfo = vcs.GetVcsInfo(c.pwd, c.User, c.Conf)
	} else {
		c.VcsInfo = &vcs.VcsInfo{}
	}

	c.Os = runtime.GOOS

	if os.Getenv("SSH_CLIENT") != "" {
		c.Ssh = true
	}
}
