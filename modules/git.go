package modules

import (
	"strconv"
	"strings"
	"time"
)

func init() {
	addModule("git", func() Module {
		return &git{
			Base: Get(),
		}
	})
}

type git struct {
	*Base

	// Output
	Ahead  int
	Behind int

	Branch    string
	Tag       string
	Detached  bool
	HasRemote bool
	Staged    int
	Conflicts int
	Changed   int
	Untracked int
	Stashes   int
}

func (g *git) Init() error {
	gitDir, err := g.Base.runCommand("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return err
	}
	lastFetch, err := g.Base.runCommand("stat", "-c", `%Y`, strings.TrimSpace(gitDir)+`/.git/FETCH_HEAD`)
	if err == nil {
		i, err := strconv.ParseInt(strings.TrimSpace(lastFetch), 10, 64)
		if err == nil && time.Since(time.Unix(i, 0)).Hours() > 1 {
			g.Base.runCommand("git", "fetch")
		}
	}
	status, err := g.Base.runCommand("git", "status", "--porcelain", "-b", "--ignore-submodules")
	if err != nil {
		return err
	}
	lines := strings.Split(status, "\n")

	g.parseBranchLine(lines[0])
	g.collectChanges(lines[1:])
	return nil
}

func (g *git) getStashCount() {
	stdout, err := g.Base.runCommand("git", "log", "--format=\"%%gd: %%gs\"", "-g", "--first-parent", "-m", "refs/stash", "--")
	if err != nil {
		return
	}

	g.Stashes = len(strings.Split(stdout, "\n")) - 1
}

func (g *git) parseBranchLine(line string) {
	if strings.Contains(line, "no Branch") {
		g.Detached = true
		g.Branch = "no Branch"
		hash, _ := g.Base.runCommand("git", "rev-parse", "--short", "HEAD")
		g.Branch += ":" + strings.TrimSpace(hash[0:len(hash)-1])
	} else if strings.Contains(line, "...") {
		g.HasRemote = true

		lines := strings.Split(line, " ")
		g.Branch = strings.Split(lines[1], "...")[0]

		if len(lines) >= 3 {
			for _, l := range lines {
				if strings.HasPrefix(l, "ahead ") {
					g.Ahead, _ = strconv.Atoi(strings.TrimPrefix(l, "ahead "))
				} else if strings.HasPrefix(l, "behind ") {
					g.Ahead, _ = strconv.Atoi(strings.TrimPrefix(l, "behind "))
					break
				}
			}
		}
	} else {
		g.Branch = strings.TrimSpace(strings.Split(line, " ")[1])
	}
}

func (g *git) collectChanges(lines []string) {
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		idxStatus := line[0]
		wtStatus := line[1]

		// " M  hoge.txt" , "AM ahoo.png" , ...
		if wtStatus != ' ' && wtStatus != '?' {
			g.Changed++
		}

		// "MT hoge.cpp" , "A  fuga.txt" , ...
		if idxStatus != ' ' && idxStatus != '?' {
			g.Staged++
		}

		// "?? hoge.txt", ...
		if idxStatus == '?' && wtStatus == '?' {
			g.Untracked++
		}

		// "UU hogehoge.txt" ...
		if idxStatus == 'U' && wtStatus == 'U' {
			g.Conflicts++
		}
	}
}

func (g *git) runCommand(command string, args ...string) (string, error) {
	args = append([]string{"--no-optional-locks", "-c", "core.quotepath=false", "-c", "color.status=false"}, args...)
	return g.Base.runCommand(command, args...)
}

func (g *git) getTag() {
	tag, _ := g.Base.runCommand("git", "describe", "--tags", "--exact")
	if tag != "" {
		g.Tag = strings.TrimSpace(tag[0 : len(tag)-1])
	}
}
