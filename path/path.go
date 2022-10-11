package path

import (
	"github.com/dustin/go-humanize"
	"github.com/metafates/pat/filesystem"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"os"
	"path/filepath"
)

type Path struct {
	raw     string
	size    mo.Option[int64]
	entries mo.Option[[]*Path]
}

func (p *Path) String() string {
	return p.raw
}

func New(path string) *Path {
	return &Path{
		raw:     path,
		size:    mo.None[int64](),
		entries: mo.None[[]*Path](),
	}
}

func (p *Path) Entries() []*Path {
	if p.entries.IsPresent() {
		return p.entries.MustGet()
	}

	descriptors, err := filesystem.Api().ReadDir(p.String())
	if err != nil {
		return make([]*Path, 0)
	}

	var size int64

	entries := lo.FilterMap(descriptors, func(info os.FileInfo, _ int) (*Path, bool) {
		// we don't care about directories since PATH entry is supposed to be a list of executables,
		// so we filter them out
		if info.IsDir() {
			return nil, false
		}

		// resolve symlinks
		realPath, err := filepath.EvalSymlinks(info.Name())
		if err != nil {
			size += info.Size()
		} else {
			realInfo, err := filesystem.Api().Stat(realPath)
			if err != nil {
				size += info.Size()
			} else {
				size += realInfo.Size()
			}
		}

		return New(info.Name()), true
	})

	p.entries = mo.Some(entries)
	p.size = mo.Some(size)
	return entries
}

func (p *Path) Exists() bool {
	exists, err := filesystem.Api().Exists(p.String())
	return err == nil && exists
}

func (p *Path) Size() int64 {
	if p.size.IsPresent() {
		return p.size.MustGet()
	}

	info, err := filesystem.Api().Stat(p.String())
	if err != nil {
		return 0
	}

	var size int64

	if info.IsDir() {
		for _, entry := range p.Entries() {
			size += entry.Size()
		}
	} else {
		size = info.Size()
	}

	p.size = mo.Some(size)
	return size
}

func (p *Path) SizeHuman() string {
	return humanize.Bytes(uint64(p.Size()))
}

func (p *Path) Eq(other *Path) bool {
	return p.String() == other.String()
}
