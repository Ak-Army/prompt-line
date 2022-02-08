package prompt

import (
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/Ak-Army/xlog"
	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
	"github.com/gookit/color"

	"github.com/Ak-Army/prompt-line/modules"
)

type Prompt struct {
	ConsoleTemplate string  `toml:"console_template"`
	ConsoleTitle    bool    `toml:"console_title"`
	FinalSpace      bool    `toml:"final_space"`
	Lines           []*Line `toml:"lines"`

	metaData toml.MetaData
}

type Line struct {
	Alignment string           `toml:"alignment"`
	Modules   []toml.Primitive `toml:"module"`
	NewLine   bool             `toml:"new_line,omitempty"`

	modules []*Module
}

type Module struct {
	Name     string `toml:"name"`
	Template string `toml:"template"`

	mod modules.Module
}

func New(file string) (*Prompt, error) {
	var err error
	p := &Prompt{}
	p.metaData, err = toml.DecodeFile(file, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Prompt) Print(width int) string {
	now := time.Now()
	wg := sync.WaitGroup{}
	for _, l := range p.Lines {
		for _, m := range l.Modules {
			module := &Module{}
			if err := p.metaData.PrimitiveDecode(m, module); err != nil {
				xlog.Error("Unable to parse module %s", m)
				continue
			}
			mod := modules.GetModule(module.Name)
			if mod == nil {
				xlog.Error("Unknown module: %s", module.Name)
				continue
			}
			if err := p.metaData.PrimitiveDecode(m, mod); err != nil {
				xlog.Error("Unable to parse module params %s", module.Name)
				continue
			}
			l.modules = append(l.modules, module)
			wg.Add(1)
			go func(m *Module) {
				n := time.Now()
				err := mod.Init()
				if err != nil {
					xlog.Warnf("Unable to init module %s", m.Name, err)
				} else {
					m.mod = mod
				}
				xlog.Infof("%s module init time: %s", m.Name, time.Since(n))
				wg.Done()
			}(module)
		}
	}
	wg.Wait()
	buffer := &strings.Builder{}
	if p.ConsoleTitle && p.ConsoleTemplate != "" {
		buffer.WriteString("\033]0;")
		t, err := template.New("consolTemplate").Parse(p.ConsoleTemplate)
		if err == nil {
			if err := t.Execute(buffer, p); err != nil {
				xlog.Error("Unable to render console title", err)
			}
		}
		buffer.WriteString("\a")
	}
	var re = regexp.MustCompile(`\x1b\[[^m]+m`)
	builder := &strings.Builder{}
	render := template.New("template").Funcs(sprig.TxtFuncMap())

	for _, l := range p.Lines {
		var left string
		if l.Alignment == "right" {
			left = builder.String()
		} else if l.NewLine {
			buffer.WriteString("\n")
		}
		builder = &strings.Builder{}

		for _, m := range l.modules {
			if m.mod == nil {
				continue
			}
			t, err := render.Parse(m.Template)
			if err != nil {
				xlog.Errorf("Unable to parse %s module template", m.Name, err)
				continue
			}
			b := &strings.Builder{}
			if err := t.Execute(b, m.mod); err != nil {
				xlog.Errorf("Unable to render %s module template", m.Name, err)
			}
			color.Fprint(builder, b.String())
		}
		content := strings.ReplaceAll(strings.ReplaceAll(builder.String(), "\n", ""), "\r", "")
		if l.Alignment == "right" && len(content) > 0 {
			leftLength := utf8.RuneCountInString(re.ReplaceAllString(left, ``))
			rightLength := utf8.RuneCountInString(re.ReplaceAllString(content, ``))
			count := width - leftLength - rightLength

			if count < 0 {
				count = 5
			}
			buffer.WriteString(strings.Repeat(" ", count))
			buffer.WriteString(content)
			buffer.WriteString("\n")
		} else {
			buffer.WriteString(content)
		}
	}
	if p.FinalSpace {
		buffer.WriteString(" ")
	}
	xlog.Infof("Render time: %s", time.Since(now))
	return buffer.String()
}
