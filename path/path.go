package path

import (
	"github.com/metafates/pat/filesystem"
)

type Path struct {
	raw string
}

func (p *Path) String() string {
	return p.raw
}

func New(path string) *Path {
	return &Path{raw: path}
}

func (p *Path) Entries() int {
	entries, err := filesystem.Api().ReadDir(p.String())
	if err != nil {
		return 0
	}

	var n int

	for _, entry := range entries {
		if !entry.IsDir() {
			n++
		}
	}

	return n
}

func (p *Path) Exists() bool {
	exists, err := filesystem.Api().Exists(p.String())
	return err == nil && exists
}

func (p *Path) Eq(other *Path) bool {
	return p.String() == other.String()
}
