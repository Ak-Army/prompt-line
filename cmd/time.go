package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/Ak-Army/cli"
)

func init() {
	cli.RootCommand().AddCommand("time", &Time{})
}

type Time struct {
	Base
}

func (c *Time) Help() string {
	return `Usage: time [command options]`
}

func (c *Time) Synopsis() string {
	return "Return time in millisecond"
}

func (c *Time) Run(_ context.Context) error {
	c.initLogger()
	fmt.Print(time.Now().UnixNano() / 1000000)
	return nil
}
