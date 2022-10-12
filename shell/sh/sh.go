package sh

import (
	"os/exec"
	"strings"
)

type Sh struct{}

func New() *Sh {
	return &Sh{}
}

func (s *Sh) Name() string {
	return "sh"
}

func (s *Sh) cmd(code string) *exec.Cmd {
	return exec.Command("sh", "-c", code)
}

func (s *Sh) AddPath(path string) error    { return nil }
func (s *Sh) RemovePath(path string) error { return nil }
func (s *Sh) Paths() ([]string, error) {
	cmd := s.cmd("echo $PATH")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(out)), ":"), nil
}
func (s *Sh) Overwrite(paths []string) error {
	return nil
}
func (s *Sh) Available() bool {
	if _, err := exec.LookPath(s.Name()); err != nil {
		return false
	}

	return true
}
