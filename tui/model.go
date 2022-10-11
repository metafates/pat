package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/metafates/pat/shell"
	"github.com/metafates/pat/stack"
	"github.com/samber/lo"
	"golang.org/x/term"
	"os"
)

type Model struct {
	statesHistory *stack.Stack[state]
	state         state

	shells        []shell.Shell
	selectedShell shell.Shell

	onSave map[string]action

	ttyWidth, ttyHeight int

	keymap *keymap

	err error

	pathSelectC  *list.Model
	shellSelectC *list.Model
	textInputC   *textinput.Model
}

func NewModel() *Model {
	model := &Model{
		keymap:        &keymap{},
		shells:        shell.AvailableShells(),
		statesHistory: stack.New[state](),
		onSave:        make(map[string]action),
	}
	model.keymap.init()

	defer func() {
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			model.Resize(width, height)
		}
	}()

	newList := func(title, singular, plural string) *list.Model {
		l := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
		l.Title = title
		l.SetStatusBarItemName(singular, plural)
		l.AdditionalShortHelpKeys = model.keymap.AdditionalShortHelpKeys
		l.AdditionalFullHelpKeys = model.keymap.AdditionalFullHelpKeys
		return &l
	}

	newTextInput := func() *textinput.Model {
		t := textinput.New()
		return &t
	}

	model.pathSelectC = newList("Paths", "path", "paths")
	model.shellSelectC = newList("Shell", "shell", "shells")
	model.textInputC = newTextInput()

	return model
}

func (m *Model) Resize(width, height int) {
	m.ttyWidth = width
	m.ttyHeight = height

	m.shellSelectC.SetWidth(width)
	m.shellSelectC.SetHeight(height)

	m.pathSelectC.SetWidth(width)
	m.pathSelectC.SetHeight(height)

	m.textInputC.Width = width
}

func (m *Model) pushState(s state) {
	if m.state == s {
		return
	}

	ignoredStates := []state{stateQuit, stateError}
	if !lo.Contains(ignoredStates, m.state) {
		m.statesHistory.Push(m.state)
	}

	m.state = s
}

func (m *Model) popState() {
	m.state = m.statesHistory.Pop()
}

func (m *Model) raiseError(err error) {
	m.pushState(stateError)
	m.err = err
}
