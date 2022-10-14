package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	model *Model

	Preview,

	Select,
	Confirm,
	Remove,
	Add,
	Save,

	Back,
	Reset,

	MoveUp,
	MoveDown,

	ForceQuit key.Binding
}

func (k *keymap) init() {
	k.Back = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Back", "Go back to the previous screen"),
	)
	k.ForceQuit = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Quit the program"),
	)
	k.Select = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select an item"),
	)
	k.Confirm = key.NewBinding(
		key.WithKeys("Y"),
		key.WithHelp("Y", "Confirm"),
	)
	k.Remove = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Remove an item"),
	)
	k.Add = key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add an item"),
	)
	k.Save = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Add the current list"),
	)
	k.Reset = key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "Reset the current list"),
	)
	k.MoveUp = key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "GenerateExport the selected item with the one above"),
	)
	k.MoveDown = key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "GenerateExport the selected item with the one below"),
	)
	k.Preview = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Preview at the selected item"),
	)
}

func (k *keymap) AdditionalShortHelpKeys() []key.Binding {
	return nil
}

func (k *keymap) AdditionalFullHelpKeys() []key.Binding {
	return nil
}
