package cmd

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/Ak-Army/xlog"
)

type Base struct {
	Verbose bool   `flag:"v, print debug and info messages"`
	Shell   string `flag:"shell, shell name"`
	Config  string `flag:"config, config file"`
	Log     string `flag:"log, log path"`
}

func (b *Base) initLogger() xlog.Logger {
	xlog.SetLogger(xlog.NopLogger)

	multiOutput := xlog.MultiOutput{}
	if b.Verbose {
		multiOutput = append(multiOutput, xlog.NewConsoleOutput())
	} else {
		multiOutput = append(multiOutput, xlog.LevelOutput{
			Fatal: xlog.NewConsoleOutputW(os.Stderr, xlog.NewLogfmtOutput(os.Stderr)),
		})
	}
	if b.Log != "" {
		logfile, err := os.OpenFile(b.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err == nil {
			multiOutput = append(multiOutput, xlog.NewLogfmtOutput(logfile))
		}
	}
	conf := xlog.Config{
		Output: multiOutput,
	}
	log.SetFlags(0)
	l := xlog.New(conf)
	xlog.SetLogger(l)
	log.SetOutput(l)

	bi, ok := debug.ReadBuildInfo()
	version := "Unknown version"
	if ok {
		version = bi.Main.Version
	}
	l.SetField("version", version)

	return l
}
