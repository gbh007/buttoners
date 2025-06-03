package screens

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/view/components"
)

type menuItem struct {
	Name              string
	ScreenConstructor func(shared SharedState) (tea.Model, tea.Cmd)
}

func (i menuItem) FilterValue() string { return i.Name }

type menuItemDelegate struct{}

func (d menuItemDelegate) Height() int { return 1 }

func (d menuItemDelegate) Spacing() int { return 0 }

func (d menuItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d menuItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}

	if index == m.Index() {
		_, _ = io.WriteString(w, components.CursorStyle.Bold(true).Render("> "+i.Name))
	} else {
		_, _ = io.WriteString(w, components.SubTitleStyle.Render("  "+i.Name))
	}
}

type Menu struct {
	shared SharedState
	list   list.Model
}

func NewMenu(shared SharedState) Menu {
	items := []list.Item{
		menuItem{
			Name: "Login",
			ScreenConstructor: func(shared SharedState) (tea.Model, tea.Cmd) {
				screen := NewLogin(shared)
				return screen, screen.Init()
			},
		},
		menuItem{
			Name: "Storage debug",
			ScreenConstructor: func(shared SharedState) (tea.Model, tea.Cmd) {
				screen := NewStorageDebug(shared)
				return screen, screen.Init()
			},
		},
	}

	l := list.New(items, menuItemDelegate{}, 20, 14)
	l.Title = "Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = components.TitleStyle
	l.Styles.PaginationStyle = components.BlurredStyle
	l.Styles.HelpStyle = components.HelpStyle

	return Menu{
		shared: shared,
		list:   l,
	}
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(menuItem)
			if ok {
				return i.ScreenConstructor(m.shared)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Menu) View() string {
	return m.list.View()
}
