package path

import (
	"github.com/dustin/go-humanize"
	"github.com/metafates/pat/filesystem"
	"github.com/metafates/pat/util"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"os"
	"path/filepath"
)

type Path struct {
	raw          string
	isDir        mo.Option[bool]
	isExecutable mo.Option[bool]
	exists       mo.Option[bool]
	size         mo.Option[int64]
	entries      mo.Option[[]*Path]
}

func (p *Path) String() string {
	return util.ResolveTilde(p.raw)
}

func New(path string) *Path {
	return &Path{
		raw:          path,
		size:         mo.None[int64](),
		entries:      mo.None[[]*Path](),
		exists:       mo.None[bool](),
		isDir:        mo.None[bool](),
		isExecutable: mo.None[bool](),
	}
}

func (p *Path) IsDir() bool {
	if p.isDir.IsPresent() {
		return p.isDir.MustGet()
	}

	isDir, _ := filesystem.Api().IsDir(p.raw)

	p.isDir = mo.Some(isDir)
	return isDir
}

func (p *Path) IsExecutable() bool {
	if p.isExecutable.IsPresent() {
		return p.isExecutable.MustGet()
	}

	stat, _ := filesystem.Api().Stat(p.raw)
	isExecutable := stat.Mode()&0111 != 0

	p.isExecutable = mo.Some(isExecutable)
	return isExecutable
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

		size += info.Size()

		return New(filepath.Join(p.raw, info.Name())), true
	})

	p.entries = mo.Some(entries)
	p.size = mo.Some(size)
	return entries
}

func (p *Path) Exists() bool {
	if p.exists.IsPresent() {
		return p.exists.MustGet()
	}

	exists, _ := filesystem.Api().Exists(p.raw)
	p.exists = mo.Some(exists)
	return exists
}

func (p *Path) Size() int64 {
	if p.size.IsPresent() {
		return p.size.MustGet()
	}

	info, err := filesystem.Api().Stat(p.raw)
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
