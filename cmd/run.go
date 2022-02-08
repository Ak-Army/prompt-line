package cmd

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Ak-Army/cli"

	"github.com/Ak-Army/prompt-line/modules"
	"github.com/Ak-Army/prompt-line/prompt"
)

func init() {
	cli.RootCommand().AddCommand("run", &Run{})
}

type Run struct {
	Base
	Width         int     `flag:"width, width of the terminal"`
	Error         int     `flag:"error, last command error"`
	ExecutionTime float64 `flag:"execution-time, last command execution time"`
}

func (c *Run) Help() string {
	return `Usage: run [command options]`
}

func (c *Run) Synopsis() string {
	return "Print prompt-line"
}

func (c *Run) Run(ctx context.Context) error {
	c.initLogger()
	p, err := prompt.New(c.Config)
	if err != nil {
		return err
	}
	modules.Get().SetExitCode(c.Error)
	modules.Get().SetExecutionTime(c.ExecutionTime)

	line := p.Print(c.Width)
	switch c.Shell {
	case "bash":
		regex := regexp.MustCompile(`(\x1b\[[^m]+m)`)
		fmt.Print(regex.ReplaceAllString(line, "\\[${1}\\]"))
	default:
		fmt.Print(line)
	}
	return nil
}
