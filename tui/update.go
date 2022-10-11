package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/pat/path"
	"github.com/metafates/pat/shell"
	"github.com/samber/lo"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Back):
			m.popState()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.Resize(msg.Width, msg.Height)
	}

	switch m.state {
	case stateShellSelect:
		return m.updateShellSelect(msg)
	case statePathSelect:
		return m.updatePathSelect(msg)
	}

	return m, nil
}

func (m *Model) updateShellSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Confirm):
			item, ok := m.shellSelectC.SelectedItem().(*item)
			if !ok {
				return m, nil
			}

			m.selectedShell, ok = item.internal.(shell.Shell)
			if !ok {
				return m, nil
			}

			paths, err := m.selectedShell.Paths()
			if err != nil {
				m.raiseError(err)
				return m, nil
			}

			m.pushState(statePathSelect)
			return m, m.pathSelectC.SetItems(
				lo.Map(paths, func(p string, _ int) list.Item {
					return m.newItem(path.New(p))
				}),
			)
		}
	}

	model, cmd := m.shellSelectC.Update(msg)
	m.shellSelectC = &model
	return m, cmd
}

func (m *Model) updatePathSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Reset):
			for _, it := range m.pathSelectC.Items() {
				item, ok := it.(*item)
				if !ok {
					continue
				}

				p, ok := item.internal.(*path.Path)
				if !ok {
					continue
				}

				m.onSave[p.String()] = actionNone
			}
		case key.Matches(msg, m.keymap.Delete):
			index := m.pathSelectC.Index()
			item, ok := m.pathSelectC.SelectedItem().(*item)
			if !ok {
				return m, nil
			}

			p, ok := item.internal.(*path.Path)
			if !ok {
				return m, nil
			}

			pathAction, ok := m.onSave[p.String()]
			if !ok {
				m.onSave[p.String()] = actionDelete
			} else if pathAction == actionAdd {
				m.pathSelectC.RemoveItem(index)
				return m, nil
			} else if pathAction != actionDelete {
				m.onSave[p.String()] = actionDelete
			} else {
				m.onSave[p.String()] = actionNone
			}

			return m, nil
		}
	}

	model, cmd := m.pathSelectC.Update(msg)
	m.pathSelectC = &model
	return m, cmd
}
