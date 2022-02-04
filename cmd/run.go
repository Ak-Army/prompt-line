package cmd

import (
	"context"

	"github.com/Ak-Army/cli"

	"github.com/Ak-Army/prompt-line/internal/prompt"
	"github.com/Ak-Army/prompt-line/modules"
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
	prompt, err := prompt.New(c.Config)
	if err != nil {
		return err
	}
	modules.Get().SetExitCode(c.Error)
	modules.Get().SetExecutionTime(c.ExecutionTime)
	return prompt.Print(c.Width)
}
