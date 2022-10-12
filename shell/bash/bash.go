package bash

import (
	"github.com/metafates/pat/log"
	"os/exec"
	"strings"
)

type Bash struct{}

func (b *Bash) Name() string {
	return "bash"
}

func New() *Bash {
	return &Bash{}
}

func (b *Bash) cmd(code string) *exec.Cmd {
	cmd := exec.Command(b.Name(), "-c", code)
	log.Infof("executing %s", code)
	return cmd
}

func (b *Bash) AddPath(path string) error { return nil }

func (b *Bash) RemovePath(path string) error { return nil }

func (b *Bash) Paths() ([]string, error) {
	cmd := b.cmd("echo $PATH")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(out)), ":"), nil
}

func (b *Bash) Overwrite(paths []string) error {
	return nil
}

func (b *Bash) Available() bool {
	if _, err := exec.LookPath(b.Name()); err != nil {
		return false
	}

	return true
}
