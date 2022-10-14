package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/pat/path"
	"github.com/metafates/pat/shell"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Back):
			var cmd tea.Cmd

			onListBack := func(l *list.Model) tea.Cmd {
				l.ResetSelected()
				l.ResetFilter()

				return l.NewStatusMessage("")
			}

			switch m.state {
			case statePathAdd:
				m.textInputC.SetValue("")
			case stateEntriesPreview:
				if m.entriesPreviewC.FilterState() != list.Unfiltered {
					model, cmd := m.entriesPreviewC.Update(msg)
					m.entriesPreviewC = &model
					return m, cmd
				}

				cmd = onListBack(m.entriesPreviewC)
			case statePathSelect:
				if m.pathSelectC.FilterState() != list.Unfiltered {
					model, cmd := m.pathSelectC.Update(msg)
					m.pathSelectC = &model
					return m, cmd
				}

				cmd = onListBack(m.pathSelectC)
			}

			m.popState()
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	}

	switch m.state {
	case stateShellSelect:
		return m.updateShellSelect(msg)
	case statePathSelect:
		return m.updatePathSelect(msg)
	case statePathAdd:
		return m.updatePathAdd(msg)
	case stateEntriesPreview:
		return m.updateEntriesPreview(msg)
	case stateConfirmActions:
		return m.updateConfirmActions(msg)
	}

	return m, nil
}

func (m *Model) updateShellSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Select):
			item, ok := m.shellSelectC.SelectedItem().(*item)
			if !ok {
				return m, nil
			}

			m.selectedShell, ok = item.internal.(*shell.Wrapper)
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
	case *path.Path:
		m.setAction(msg, actionAdd)
		m.pathSelectC.SetItem(-1, m.newItem(msg))
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Reset):
			return m, m.reset()
		case key.Matches(msg, m.keymap.Save, m.keymap.Select):
			m.pushState(stateConfirmActions)
			return m, nil
		case key.Matches(msg, m.keymap.Add):
			m.pushState(statePathAdd)
			return m, nil
		case key.Matches(msg, m.keymap.MoveUp):
			idx := m.pathSelectC.Index()

			if idx == 0 {
				return m, nil
			}

			items := m.pathSelectC.Items()
			a, b := items[idx-1], items[idx]
			items[idx-1], items[idx] = b, a

			m.pathSelectC.SetItem(idx-1, b)
			m.pathSelectC.SetItem(idx, a)

			m.pathSelectC.CursorUp()

			m.order = mo.Some(lo.Map(items, func(i list.Item, _ int) string {
				return i.(*item).internal.(*path.Path).String()
			}))

			return m, nil
		case key.Matches(msg, m.keymap.MoveDown):
			idx := m.pathSelectC.Index()
			items := m.pathSelectC.Items()
			if idx == len(items)-1 {
				return m, nil
			}

			a, b := items[idx], items[idx+1]
			items[idx], items[idx+1] = b, a

			m.pathSelectC.SetItem(idx, b)
			m.pathSelectC.SetItem(idx+1, a)

			m.pathSelectC.CursorDown()

			m.order = mo.Some(lo.Map(items, func(i list.Item, _ int) string {
				return i.(*item).internal.(*path.Path).String()
			}))

			return m, nil
		case key.Matches(msg, m.keymap.Preview):
			item, ok := m.pathSelectC.SelectedItem().(*item)
			if !ok {
				return m, nil
			}

			p, ok := item.internal.(*path.Path)
			if !ok {
				return m, nil
			}

			m.pushState(stateEntriesPreview)
			return m, m.entriesPreviewC.SetItems(
				lo.Map(p.Entries(), func(e *path.Path, _ int) list.Item {
					return m.newItem(e)
				}),
			)
		case key.Matches(msg, m.keymap.Remove):
			index := m.pathSelectC.Index()
			item, ok := m.pathSelectC.SelectedItem().(*item)
			if !ok {
				return m, nil
			}

			p, ok := item.internal.(*path.Path)
			if !ok {
				return m, nil
			}

			pathAction, ok := m.getAction(p)
			if !ok {
				m.setAction(p, actionRemove)
			} else if pathAction == actionAdd {
				m.pathSelectC.RemoveItem(index)
				m.setAction(p, actionNone)
				return m, nil
			} else if pathAction != actionRemove {
				m.setAction(p, actionRemove)
			} else {
				m.setAction(p, actionNone)
			}

			return m, nil
		case key.Matches(msg, m.keymap.Copy):
			item, ok := m.pathSelectC.SelectedItem().(*item)
			if !ok {
				return m, nil
			}

			item.Copy()
			return m, m.pathSelectC.NewStatusMessage("Copied to clipboard")
		}
	}

	model, cmd := m.pathSelectC.Update(msg)
	m.pathSelectC = &model
	return m, cmd
}

func (m *Model) updateEntriesPreview(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Back):
			m.popState()
			return m, nil
		}
	}

	model, cmd := m.entriesPreviewC.Update(msg)
	m.entriesPreviewC = &model
	return m, cmd
}

func (m *Model) updateConfirmActions(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Reset):
			m.popState()
			return m, m.reset()
		case key.Matches(msg, m.keymap.Confirm) && m.hasUnsaved():
			err := m.save()
			if err != nil {
				m.raiseError(err)
				return m, nil
			}

			m.popState()
			return m, tea.Batch(m.reset(), m.pathSelectC.NewStatusMessage("Saved"))
		}
	}

	return m, nil
}

func (m *Model) updatePathAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Select) && m.textInputC.Value() != "":
			m.popState()
			p := path.New(m.textInputC.Value())

			m.setAction(p, actionAdd)
			m.pathSelectC.InsertItem(-1, m.newItem(p))
			m.pathSelectC.Select(0)
			m.textInputC.SetValue("")
			return m, nil
		}
	}

	model, cmd := m.textInputC.Update(msg)
	m.textInputC = &model
	return m, cmd
}
