package sh

import (
	"os/exec"
)

type Sh struct{}

func New() *Sh {
	return &Sh{}
}

func (s *Sh) Name() string {
	return "sh"
}
func (s *Sh) AddPath(path string) error    { return nil }
func (s *Sh) RemovePath(path string) error { return nil }
func (s *Sh) Paths() ([]string, error)     { return nil, nil }
func (s *Sh) Overwrite(paths []string) error {
	return nil
}
func (s *Sh) Available() bool {
	if _, err := exec.LookPath(s.Name()); err != nil {
		return false
	}

	return true
}
