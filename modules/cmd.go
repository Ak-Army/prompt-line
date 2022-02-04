package modules

import (
	"errors"
	"strings"
)

func init() {
	addModule("cmd", func() Module {
		return &cmd{
			Base: Get(),
		}
	})
}

type cmd struct {
	*Base
	// Input
	Shell   string
	Command string

	// Output
	Value string
	Err   []string
}

func (c *cmd) Init() error {
	if c.Shell == "" || c.Command == "" {
		return errors.New("undefined shell or command")
	}
	if strings.Contains(c.Command, "||") {
		commands := strings.Split(c.Command, "||")
		for _, cmd := range commands {
			output, err := c.Base.runCommand(c.Shell, cmd)
			if err != nil {
				c.Err = append(c.Err, err.Error())
				continue
			}
			if output != "" {
				c.Value = output
				return nil
			}
		}
	}
	if strings.Contains(c.Command, "&&") {
		commands := strings.Split(c.Command, "&&")
		for _, cmd := range commands {
			output, err := c.Base.runCommand(c.Shell, cmd)
			if err != nil {
				c.Err = append(c.Err, err.Error())
				continue
			}
			c.Value += output
		}
		return nil
	}
	output, err := c.Base.runCommand(c.Shell, c.Command)
	if err != nil {
		c.Err = append(c.Err, err.Error())
	}
	c.Value += output

	return nil
}
