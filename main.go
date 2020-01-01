// Package main entry point of this tools
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chaopeng/ph/context"
	"github.com/chaopeng/ph/theme"
	"github.com/chaopeng/ph/theme/powerline"
	"github.com/chaopeng/ph/theme/simple"
)

// ph prompt $last_status, returns prompt string
// ph tmux $pwd_path, return tmux status left string
func main() {
	theme.RegisterTheme("powerline", &powerline.Theme{})
	theme.RegisterTheme("simple", &simple.Theme{})

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalln("Invalidate Arguments: should be: `ph prompt last_status` or `ph tmux`")
	}

	ctx := context.CreateContext()

	if args[0] == "prompt" || args[0] == "tmux" {
		theme.Render(args[0], args[1], ctx)
	} else if args[0] == "shortpath" {
		ctx.ReadPathInfo()
		fmt.Printf(ctx.PathInfo.ShorternPrefix)
		fmt.Printf(ctx.PathInfo.BaseDir)
	}
}
