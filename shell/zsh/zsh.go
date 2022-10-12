package zsh

import (
	"os/exec"
	"strings"
)

type Zsh struct{}

func (z *Zsh) Name() string {
	return "zsh"
}

func New() *Zsh {
	return &Zsh{}
}

func (z *Zsh) cmd(code string) *exec.Cmd {
	return exec.Command("zsh", "-c", code)
}

func (z *Zsh) AddPath(path string) error    { return nil }
func (z *Zsh) RemovePath(path string) error { return nil }
func (z *Zsh) Paths() ([]string, error) {
	cmd := z.cmd("echo $PATH")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(out)), ":"), nil
}
func (z *Zsh) Overwrite(paths []string) error {
	return nil
}
func (z *Zsh) Available() bool {
	if _, err := exec.LookPath(z.Name()); err != nil {
		return false
	}

	return true
}
