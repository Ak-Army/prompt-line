package modules

import (
	"bytes"
	"os/exec"
	"os/user"
	"strings"
)

type Module interface {
	Init() error
}

type initFunc func() Module

type Base struct {
	exitCode      int
	executionTime float64
}

var b *Base
var modules map[string]initFunc

func Get() *Base {
	if b == nil {
		b = &Base{}
	}
	return b
}

func GetModule(name string) Module {
	if m, ok := modules[name]; ok {
		return m()
	}
	return nil
}

func addModule(name string, fn initFunc) {
	if modules == nil {
		modules = make(map[string]initFunc)
	}
	if _, ok := modules[name]; !ok {
		modules[name] = fn
	}
}

func (b *Base) runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmdErr := cmd.Run()
	if cmdErr != nil {
		output := err.String()
		return output, cmdErr
	}
	output := strings.TrimSuffix(out.String(), "\n")

	return output, nil
}

func (b *Base) homeDir() string {
	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir
	}
	return ""
}

func (b *Base) SetExitCode(code int) {
	b.exitCode = code
}

func (b *Base) SetExecutionTime(executionTime float64) {
	b.executionTime = executionTime
}
