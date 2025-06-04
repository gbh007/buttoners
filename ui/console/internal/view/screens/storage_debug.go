package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/model"
	"github.com/gbh007/buttoners/ui/console/internal/view/components"
)

type StorageDebugEvent struct {
	isInit bool
	data   model.Connection
	err    error
}

type StorageDebug struct {
	shared    *SharedState
	data      model.Connection
	loader    components.Loader
	lastError error
}

func NewStorageDebug(shared *SharedState) StorageDebug {
	return StorageDebug{
		shared: shared,
		loader: components.NewLoader("Загрузка"),
	}
}

func (m StorageDebug) refreshData() (StorageDebug, tea.Cmd) {
	var cmd tea.Cmd
	m.loader, cmd = m.loader.Activate()

	cmd = tea.Batch(cmd, func() tea.Msg {
		var event StorageDebugEvent
		event.data, event.err = m.shared.Storage.GetConnectionData(m.shared.Ctx)
		return event
	})

	return m, cmd
}

func (m StorageDebug) Init() tea.Cmd {
	return func() tea.Msg {
		return StorageDebugEvent{
			isInit: true,
		}
	}
}

func (m StorageDebug) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			nextScreen := NewMenu(m.shared)
			return nextScreen, nextScreen.Init()

		case "ctrl+r":
			return m.refreshData()
		}
	case StorageDebugEvent:
		if msg.isInit {
			return m.refreshData()
		}

		var cmd tea.Cmd
		m.loader, cmd = m.loader.Deactivate()
		m.data = msg.data
		m.lastError = msg.err
		return m, cmd
	}

	var cmd tea.Cmd
	m.loader, cmd = m.loader.Update(msg)

	return m, cmd
}

func (m StorageDebug) View() string {
	var b strings.Builder

	b.WriteString(components.TitleStyle.Render("Данные хранилища:"))
	b.WriteString("\n\n")

	b.WriteString(components.RenderError(m.lastError))

	if m.loader.IsActivate() {
		b.WriteString(m.loader.View())
	} else {
		b.WriteString(components.SubTitleStyle.Render("Хост:"))
		b.WriteString(" " + m.data.Addr + "\n")

		b.WriteString(components.SubTitleStyle.Render("Логин:"))
		b.WriteString(" " + m.data.Login + "\n")

		b.WriteString(components.SubTitleStyle.Render("Пароль:"))
		b.WriteString(" " + m.data.Password + "\n")
	}

	b.WriteString("\n")

	b.WriteString(components.HelpStyle.Render("Для обновления нажми ctrl+r"))

	return b.String()
}
