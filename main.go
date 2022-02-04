package main

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/Ak-Army/cli"

	_ "github.com/Ak-Army/prompt-line/cmd"
)

func main() {
	bi, ok := debug.ReadBuildInfo()
	version := "Unknown version"
	if ok {
		version = bi.Main.Version
	}
	c := cli.New("prompt-line", version)
	cli.RootCommand().Authors = []string{"Ak-Army"}
	c.Run(context.Background(), os.Args)
}
