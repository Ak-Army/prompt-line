package modules

import (
	"os"
	"runtime"
	"strings"
)

func init() {
	addModule("session", func() Module {
		return &session{
			Base: Get(),
		}
	})
}

type session struct {
	*Base

	// Output
	UserName     string
	ComputerName string
	IsSSHSession bool
	IsRoot       bool
}

func (s *session) Init() error {
	s.IsSSHSession = s.activeSSHSession()
	s.UserName = s.getUserName()
	s.ComputerName = s.getComputerName()
	s.IsRoot = os.Geteuid() == 0

	return nil
}

func (s *session) getUserName() string {
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	username := strings.TrimSpace(user)
	if runtime.GOOS == "windows" && strings.Contains(username, "\\") {
		username = strings.Split(username, "\\")[1]
	}
	return username
}

func (s *session) getComputerName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return strings.TrimSpace(hostname)
}

func (s *session) activeSSHSession() bool {
	keys := []string{
		"SSH_CONNECTION",
		"SSH_CLIENT",
	}
	for _, key := range keys {
		content := os.Getenv(key)
		if content != "" {
			return true
		}
	}
	return false
}
