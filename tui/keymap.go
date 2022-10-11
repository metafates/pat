package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	model *Model

	Confirm,
	Delete,
	Add,
	Save,

	Back,

	Quit,
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
	k.Quit = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Quit the program"),
	)

	k.Confirm = key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter", "Select an item"),
	)
	k.Delete = key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete an item"),
	)
	k.Add = key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add an item"),
	)
	k.Save = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Add the current list"),
	)
}

func (k *keymap) AdditionalShortHelpKeys() []key.Binding {
	return nil
}

func (k *keymap) AdditionalFullHelpKeys() []key.Binding {
	return nil
}
