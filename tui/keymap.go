package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/samber/lo"
)

type keymap struct {
	model *Model

	Preview,

	AcceptCompletion,

	Select,
	Confirm,
	Remove,
	Add,
	Save,
	Copy,

	Back,
	Reset,

	MoveUp,
	MoveDown,

	ForceQuit key.Binding
}

func (k *keymap) ShortHelp() []key.Binding {
	switch k.model.state {
	case stateError:
		return []key.Binding{k.Back, k.ForceQuit}
	case stateShellSelect:
		return []key.Binding{k.Select, k.Back}
	case stateConfirmActions:
		return []key.Binding{k.Confirm, k.Back, k.ForceQuit}
	case statePathSelect:
		return []key.Binding{k.Select, k.Remove, k.Add, k.Save, k.Copy, k.Back}
	case statePathAdd:
		s := k.Select
		s.SetHelp("enter", "add the current path")
		return []key.Binding{s, k.AcceptCompletion, k.Back}
	case stateEntriesPreview:
		return []key.Binding{k.Back}
	default:
		return nil
	}
}

func (k *keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func (k *keymap) init() {
	k.Back = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	)
	k.ForceQuit = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
	k.AcceptCompletion = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "accept completion"),
	)
	k.Select = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	)
	k.Confirm = key.NewBinding(
		key.WithKeys("Y"),
		key.WithHelp("Y", "confirm"),
	)
	k.Remove = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "remove"),
	)
	k.Copy = key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy"),
	)
	k.Add = key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	)
	k.Save = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save"),
	)
	k.Reset = key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "reset changes"),
	)
	k.MoveUp = key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "swap the selected item with the one above"),
	)
	k.MoveDown = key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "swap the selected item with the one below"),
	)
	k.Preview = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "preview at the selected item"),
	)
}

func (k *keymap) AdditionalShortHelpKeys() []key.Binding {
	return k.ShortHelp()
}

func (k *keymap) AdditionalFullHelpKeys() []key.Binding {
	return lo.Flatten(k.FullHelp())
}
