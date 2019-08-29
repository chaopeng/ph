// Package main entry point of this tools
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chaopeng/ph/context"
	"github.com/chaopeng/ph/theme"
)

// ph prompt $last_status, returns prompt string
// ph tmux $pwd_path, return tmux status left string
func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalln("Invalidate Arguments: should be: `ph prompt last_status` or `ph tmux`")
	}

	ctx := context.CreateContext()

	if args[0] == "prompt" {
		prompt(args[1], ctx)
	} else if args[0] == "tmux" {
		tmuxLeft(ctx)
	} else if args[0] == "shortpath" {
		ctx.ReadPathInfo()
		fmt.Printf(ctx.PathInfo.ShorternPrefix)
		fmt.Printf(ctx.PathInfo.BaseDir)
	}
}

func prompt(lastStatus string, ctx *context.Context) {
	theme.PromptThemeRender(lastStatus, ctx)
}

func tmuxLeft(ctx *context.Context) {
	ctx.ReadCompleteInfo()
	theme.TmuxThemeRender(ctx)
}
