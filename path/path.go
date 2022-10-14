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
	isSymlink    mo.Option[bool]
	exists       mo.Option[bool]
	size         mo.Option[int64]
	executables  mo.Option[[]*Path]
}

func (p *Path) String() string {
	path := util.ResolveTilde(p.raw)
	if filepath.IsAbs(path) {
		return path
	}

	cwd, _ := os.Getwd()
	return filepath.Join(cwd, path)
}

func New(path string) *Path {
	return &Path{
		raw:          path,
		size:         mo.None[int64](),
		executables:  mo.None[[]*Path](),
		exists:       mo.None[bool](),
		isDir:        mo.None[bool](),
		isExecutable: mo.None[bool](),
		isSymlink:    mo.None[bool](),
	}
}

func (p *Path) IsDir() bool {
	if p.isDir.IsPresent() {
		return p.isDir.MustGet()
	}

	isDir, _ := filesystem.Api().IsDir(p.String())

	p.isDir = mo.Some(isDir)
	return isDir
}

func (p *Path) IsExecutable() bool {
	if p.isExecutable.IsPresent() {
		return p.isExecutable.MustGet()
	}

	stat, _ := filesystem.Api().Stat(p.String())
	isExecutable := stat.Mode()&0111 != 0

	p.isExecutable = mo.Some(isExecutable)
	return isExecutable
}

func (p *Path) IsSymLink() bool {
	if p.isSymlink.IsPresent() {
		return p.isSymlink.MustGet()
	}

	stat, err := filesystem.Api().Stat(p.String())
	if err != nil {
		p.isSymlink = mo.Some(false)
		return false
	}

	isSymlink := stat.Mode()&os.ModeSymlink != 0
	p.isSymlink = mo.Some(isSymlink)
	return isSymlink
}

func (p *Path) Executables() []*Path {
	if p.executables.IsPresent() {
		return p.executables.MustGet()
	}

	var path string
	if p.IsSymLink() {
		// find the target of the symlink
		realPath, err := filepath.EvalSymlinks(p.String())
		if err != nil {
			p.executables = mo.Some([]*Path{})
			return p.executables.MustGet()
		}

		path = realPath
	} else {
		path = p.String()
	}

	descriptors, err := filesystem.Api().ReadDir(path)
	if err != nil {
		return make([]*Path, 0)
	}

	var size int64

	executables := lo.FilterMap(descriptors, func(info os.FileInfo, _ int) (*Path, bool) {
		// we don't care about directories since PATH entry is supposed to be a list of executables,
		// so we filter them out
		if info.IsDir() {
			return nil, false
		}

		// remove non executable files
		if info.Mode()&0111 == 0 {
			return nil, false
		}

		size += info.Size()

		return New(filepath.Join(p.String(), info.Name())), true
	})

	p.executables = mo.Some(executables)
	p.size = mo.Some(size)
	return executables
}

func (p *Path) Exists() bool {
	if p.exists.IsPresent() {
		return p.exists.MustGet()
	}

	if p.IsSymLink() {
		p.exists = mo.Some(true)
		return true
	}

	exists, _ := filesystem.Api().Exists(p.String())
	p.exists = mo.Some(exists)
	return exists
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
		for _, entry := range p.Executables() {
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
