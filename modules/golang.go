package modules

import (
	"strings"

	"golang.org/x/mod/semver"
)

func init() {
	addModule("golang", func() Module {
		return &golang{
			Base: Get(),
		}
	})
}

type golang struct {
	*Base

	// Output
	GoModVersion string
	GoVersion    string
	Compare      int
}

func (g *golang) Init() error {
	goModVersion, err := g.Base.runCommand("go", "list", "-m", "-f={{.GoVersion}}")
	if err != nil {
		return nil
	}
	if goModVersion != "" {
		g.GoModVersion = "v" + goModVersion
	}
	goVersion, err := g.Base.runCommand("go", "version")
	if err != nil {
		return nil
	}
	gv := strings.Split(goVersion, " ")
	if len(gv) > 2 {
		g.GoVersion = "v" + gv[2][2:]
	}
	g.Compare = semver.Compare(g.GoModVersion, semver.MajorMinor(g.GoVersion))

	return nil
}
