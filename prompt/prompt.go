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

var re = regexp.MustCompile(`\x1b\[[^m]+m`)

type Prompt struct {
	ConsoleTemplate string  `toml:"console_template"`
	ConsoleTitle    bool    `toml:"console_title"`
	FinalSpace      bool    `toml:"final_space"`
	Lines           []*Line `toml:"lines"`

	metaData toml.MetaData
	lastLine *strings.Builder
	render   *template.Template
	buffer   *strings.Builder
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
	wg := &sync.WaitGroup{}
	for _, l := range p.Lines {
		for _, m := range l.Modules {
			p.initModule(m, l, wg)
		}
	}
	wg.Wait()
	p.buffer = &strings.Builder{}
	p.renderConsoleTitle(p.buffer)

	p.lastLine = &strings.Builder{}
	p.render = template.New("template").Funcs(sprig.TxtFuncMap())

	for _, l := range p.Lines {
		p.renderModules(l, width)
	}
	if p.FinalSpace {
		p.buffer.WriteString(" ")
	}
	xlog.Infof("Render time: %s", time.Since(now))
	return p.buffer.String()
}

func (p *Prompt) renderConsoleTitle(buffer *strings.Builder) {
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
}

func (p *Prompt) initModule(m toml.Primitive, l *Line, wg *sync.WaitGroup) {
	module := &Module{}
	if err := p.metaData.PrimitiveDecode(m, module); err != nil {
		xlog.Error("Unable to parse module %s", m)
		return
	}
	mod := modules.GetModule(module.Name)
	if mod == nil {
		xlog.Error("Unknown module: %s", module.Name)
		return
	}
	if err := p.metaData.PrimitiveDecode(m, mod); err != nil {
		xlog.Error("Unable to parse module params %s", module.Name)
		return
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

func (p *Prompt) renderModules(l *Line, width int) {
	var left string
	if l.Alignment == "right" {
		left = p.lastLine.String()
	} else if l.NewLine {
		p.buffer.WriteString("\n")
	}
	p.lastLine = &strings.Builder{}

	for _, m := range l.modules {
		if m.mod == nil {
			continue
		}
		t, err := p.render.Parse(m.Template)
		if err != nil {
			xlog.Errorf("Unable to parse %s module template", m.Name, err)
			continue
		}
		b := &strings.Builder{}
		if err := t.Execute(b, m.mod); err != nil {
			xlog.Errorf("Unable to render %s module template", m.Name, err)
		}
		color.Fprint(p.lastLine, b.String())
	}
	content := strings.ReplaceAll(strings.ReplaceAll(p.lastLine.String(), "\n", ""), "\r", "")
	if l.Alignment == "right" && len(content) > 0 {
		leftLength := utf8.RuneCountInString(re.ReplaceAllString(left, ``))
		rightLength := utf8.RuneCountInString(re.ReplaceAllString(content, ``))
		count := width - leftLength - rightLength

		if count < 0 {
			count = 5
		}
		p.buffer.WriteString(strings.Repeat(" ", count))
		p.buffer.WriteString(content)
		p.buffer.WriteString("\n")
	} else {
		p.buffer.WriteString(content)
	}
}
