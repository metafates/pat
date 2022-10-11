package zsh

import (
	"os/exec"
)

type Zsh struct{}

func (z *Zsh) Name() string {
	return "zsh"
}

func New() *Zsh {
	return &Zsh{}
}

func (z *Zsh) AddPath(path string) error    { return nil }
func (z *Zsh) RemovePath(path string) error { return nil }
func (z *Zsh) Paths() ([]string, error)     { return nil, nil }
func (z *Zsh) Available() bool {
	if _, err := exec.LookPath(z.Name()); err != nil {
		return false
	}

	return true
}
