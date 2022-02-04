package modules

import (
	"os"
	"path/filepath"
	"strings"
)

func init() {
	addModule("path", func() Module {
		return &path{
			Base: Get(),
		}
	})
}

type path struct {
	*Base

	// Input
	Aliases map[string]string
	Length  int

	// Output
	Path string
}

func (p *path) Init() error {
	if p.Aliases == nil {
		p.Aliases = make(map[string]string)
	}
	home := p.Base.homeDir()
	if home != "" {
		if s, ok := p.Aliases["home"]; ok {
			p.Aliases[home] = s
		} else if s, ok := p.Aliases["~"]; ok {
			p.Aliases[home] = s
		}
	}
	var err error
	p.Path, err = os.Getwd()
	if err != nil {
		return err
	}
	for f, i := range p.Aliases {
		prefix := filepath.Clean(f)
		if strings.HasPrefix(p.Path, prefix) {
			p.Path = i + p.Path[len(prefix):]
			break
		}
	}
	if p.Length != 0 {
		dirs := strings.Split(p.Path, string(os.PathSeparator))
		l := len(dirs)
		if p.Length+1 < l {

			p.Path = strings.Join(dirs[0:p.Length], string(os.PathSeparator)) +
				string(os.PathSeparator) +
				"..." +
				string(os.PathSeparator) +
				dirs[l-1]
		}
	}
	return nil
}
