package bash

import (
	"os/exec"
)

type Bash struct{}

func (b *Bash) Name() string {
	return "bash"
}

func New() *Bash {
	return &Bash{}
}

func (b *Bash) AddPath(path string) error    { return nil }
func (b *Bash) RemovePath(path string) error { return nil }
func (b *Bash) Paths() ([]string, error)     { return nil, nil }
func (b *Bash) Available() bool {
	if _, err := exec.LookPath(b.Name()); err != nil {
		return false
	}

	return true
}
