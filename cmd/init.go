package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/Ak-Army/cli"
)

func init() {
	cli.RootCommand().AddCommand("init", &Init{})
}

//go:embed init/bash.templ
var bashInit string

//go:embed init/zsh.templ
var zshInit string

//go:embed init/fish.templ
var fishInit string

type Init struct {
	Base
	App string `flag:"-"`
}

func (c *Init) Help() string {
	return `Usage: Init [command options]`
}

func (c *Init) Synopsis() string {
	return "Initialize the shell"
}

func (c *Init) Run(_ context.Context) error {
	c.initLogger()
	if c.Shell == "" {
		return errors.New("undefined shell params")
	}
	if c.Shell == "" {
		return errors.New("undefined shell params")
	}
	c.App = os.Args[0]

	var t *template.Template
	var err error
	switch c.Shell {
	case "bash":
		t, err = template.New("init").Parse(bashInit)
	case "zsh":
		t, err = template.New("init").Parse(zshInit)
	case "fish":
		t, err = template.New("init").Parse(fishInit)
	default:
		err = errors.New("unknown shell")
	}
	if err != nil {
		fmt.Println("")
		return err
	}
	buffer := new(bytes.Buffer)
	defer buffer.Reset()
	err = t.Execute(buffer, c)
	if err != nil {
		return err
	}
	fmt.Println(buffer.String())
	return nil
}
