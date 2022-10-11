package tui

import (
	"github.com/metafates/pat/icon"
	"github.com/metafates/pat/path"
	"github.com/metafates/pat/shell"
	"github.com/metafates/pat/util"
	"github.com/samber/lo"
	"os/exec"
	"strings"
)

type item struct {
	model    *Model
	internal any
}

func (m *Model) newItem(internal any) *item {
	return &item{internal: internal, model: m}
}

func (i *item) marked() (markIcon string, marked bool) {
	switch x := i.internal.(type) {
	case shell.Shell:
		return "✓", x == i.model.selectedShell
	case *path.Path:
		pathAction, ok := i.model.onSave[x.String()]
		if !ok {
			return
		}

		switch pathAction {
		case actionAdd:
			return icon.Heart, true
		case actionDelete:
			return icon.Trash, true
		default:
			return
		}
	default:
		return
	}
}

func (i *item) FilterValue() string {
	switch i := i.internal.(type) {
	case shell.Shell:
		return i.Name()
	case *path.Path:
		return i.String()
	default:
		panic("unknown type")
	}
}

func (i *item) Title() string {
	title := strings.Builder{}
	title.WriteString(i.FilterValue())

	if markIcon, marked := i.marked(); marked {
		title.WriteString(" ")
		switch i.internal.(type) {
		case shell.Shell:
			title.WriteString(markIcon)
		case *path.Path:
			title.WriteString(markIcon)
		}
	}

	return title.String()
}

func (i *item) Description() string {
	switch i := i.internal.(type) {
	case shell.Shell:
		return lo.Must(exec.LookPath(i.Name()))
	case *path.Path:
		return util.Quantify(i.Entries(), "entry", "entries")
	default:
		panic("unknown type")
	}
}
