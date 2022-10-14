package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/path"
	"github.com/metafates/pat/shell"
	"github.com/metafates/pat/stack"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"golang.org/x/term"
	"os"
	"time"
)

type Styles struct {
	Padding     lipgloss.Style
	ListPadding lipgloss.Style
	Title       lipgloss.Style
}

type Model struct {
	styles *Styles

	statesHistory *stack.Stack[state]
	state         state

	shells        []*shell.Wrapper
	selectedShell *shell.Wrapper

	onSave map[string]action
	order  mo.Option[[]string]

	ttyWidth, ttyHeight int

	keymap *keymap

	err error

	pathSelectC     *list.Model
	shellSelectC    *list.Model
	entriesPreviewC *list.Model
	textInputC      *textinput.Model
	helpC           *help.Model

	toComplete mo.Option[string]
}

func NewModel() *Model {
	styles := &Styles{
		Padding: lipgloss.NewStyle().Padding(1, 2),
		Title:   lipgloss.NewStyle().Padding(0, 1).Background(color.New("62")).Foreground(color.New("230")),
	}
	styles.ListPadding = lipgloss.NewStyle().Padding(styles.Padding.GetPaddingTop(), styles.Padding.GetPaddingBottom(), 1, 0)

	model := &Model{
		keymap:        &keymap{},
		shells:        shell.AvailableShells(),
		statesHistory: stack.New[state](),
		onSave:        make(map[string]action),
		order:         mo.None[[]string](),
		styles:        styles,
	}
	model.keymap.model = model
	model.keymap.init()

	defer func() {
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			model.resize(width, height)
		}
	}()

	newList := func(title, singular, plural string) *list.Model {
		delegate := list.NewDefaultDelegate()
		delegate.Styles.SelectedTitle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("5")).
			Foreground(lipgloss.Color("5")).
			Padding(0, 0, 0, 1)
		delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Copy().Foreground(color.White)

		delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.Copy()

		l := list.New(make([]list.Item, 0), delegate, 0, 0)
		l.Title = title
		l.SetStatusBarItemName(singular, plural)
		l.AdditionalShortHelpKeys = model.keymap.AdditionalShortHelpKeys
		l.AdditionalFullHelpKeys = model.keymap.AdditionalFullHelpKeys
		l.StatusMessageLifetime = time.Millisecond * 1500
		l.Styles.NoItems = model.styles.Padding
		return &l
	}

	newTextInput := func() *textinput.Model {
		t := textinput.New()
		return &t
	}

	newHelp := func() *help.Model {
		h := help.New()
		return &h
	}

	model.pathSelectC = newList("Paths", "path", "paths")
	model.pathSelectC.SetFilteringEnabled(false)

	model.shellSelectC = newList("Select Shell", "shell", "shells")
	model.shellSelectC.SetFilteringEnabled(false)

	model.entriesPreviewC = newList("Preview", "entry", "entries")
	model.textInputC = newTextInput()

	model.helpC = newHelp()

	return model
}

func (m *Model) resize(width, height int) {
	pX, pY := m.styles.Padding.GetFrameSize()
	lX, lY := m.styles.ListPadding.GetFrameSize()

	m.ttyWidth, m.ttyHeight = width-pX, height-pY

	width, height = width-lX, height-lY

	m.shellSelectC.SetSize(width, height)
	m.pathSelectC.SetSize(width, height)
	m.entriesPreviewC.SetSize(width, height)

	m.textInputC.Width = width
	m.helpC.Width = width
}

func (m *Model) pushState(s state) {
	if m.state == s {
		return
	}

	ignoredStates := []state{stateError}
	if !lo.Contains(ignoredStates, m.state) {
		m.statesHistory.Push(m.state)
	}

	m.state = s
}

func (m *Model) popState() {
	m.state = m.statesHistory.Pop()
}

func (m *Model) getAction(p *path.Path) (action, bool) {
	pathAction, ok := m.onSave[p.String()]
	return pathAction, ok
}

func (m *Model) setAction(p *path.Path, a action) {
	m.onSave[p.String()] = a
}

func (m *Model) raiseError(err error) {
	m.pushState(stateError)
	m.err = err
}

func (m *Model) reset() tea.Cmd {
	for _, it := range m.pathSelectC.Items() {
		item, ok := it.(*item)
		if !ok {
			continue
		}

		p, ok := item.internal.(*path.Path)
		if !ok {
			continue
		}

		m.setAction(p, actionNone)
	}
	paths, err := m.selectedShell.Paths()
	if err != nil {
		m.raiseError(err)
		return nil
	}
	m.pathSelectC.Select(0)
	m.order = mo.None[[]string]()
	return m.pathSelectC.SetItems(
		lo.Map(paths, func(p string, _ int) list.Item {
			return m.newItem(path.New(p))
		}),
	)
}

func (m *Model) hasUnsaved() bool {
	for _, a := range m.onSave {
		if a != actionNone {
			return true
		}
	}

	return m.order.IsPresent()
}

func (m *Model) save() (err error) {
	for p, a := range m.onSave {
		switch a {
		case actionRemove:
			err = m.selectedShell.RemovePath(p)
		case actionAdd:
			err = m.selectedShell.AddPath(p)
		}

		if err != nil {
			return
		}
	}

	if m.order.IsPresent() {
		err = m.selectedShell.Export(m.order.MustGet())
	}

	return
}
