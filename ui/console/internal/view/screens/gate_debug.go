package screens

import (
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/view/components"
)

type GateDebugEvent struct {
	isInit bool
	data   struct {
		Count     int64
		LastPress time.Time
	}
	err error
}

type GateDebug struct {
	shared *SharedState
	data   struct {
		Count     int64
		LastPress time.Time
	}
	loader    components.Loader
	lastError error
}

func NewGateDebug(shared *SharedState) GateDebug {
	return GateDebug{
		shared: shared,
		loader: components.NewLoader("Загрузка"),
	}
}

func (m GateDebug) refreshData() (GateDebug, tea.Cmd) {
	var cmd tea.Cmd
	m.loader, cmd = m.loader.Activate()

	cmd = tea.Batch(cmd, func() tea.Msg {
		var event GateDebugEvent
		event.data.Count, event.data.LastPress, event.err = m.shared.GateClient.Activity(m.shared.Ctx, m.shared.GateToken)
		return event
	})

	return m, cmd
}

func (m GateDebug) pressButton() (GateDebug, tea.Cmd) {
	var cmd tea.Cmd
	m.loader, cmd = m.loader.Activate()

	cmd = tea.Batch(cmd, func() tea.Msg {
		var event GateDebugEvent
		err := m.shared.GateClient.ButtonClick(m.shared.Ctx, m.shared.GateToken, 1, 100)
		if err != nil {
			return GateDebugEvent{
				err: err,
			}
		}

		event.data.Count, event.data.LastPress, event.err = m.shared.GateClient.Activity(m.shared.Ctx, m.shared.GateToken)
		return event
	})

	return m, cmd
}

func (m GateDebug) Init() tea.Cmd {
	return func() tea.Msg {
		return GateDebugEvent{
			isInit: true,
		}
	}
}

func (m GateDebug) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+r":
			return m.refreshData()
		case "ctrl+v":
			return m.pressButton()
		}
	case GateDebugEvent:
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

func (m GateDebug) View() string {
	var b strings.Builder

	b.WriteString(components.TitleStyle.Render("Данные хранилища:"))
	b.WriteString("\n\n")

	b.WriteString(components.RenderError(m.lastError))

	if m.loader.IsActivate() {
		b.WriteString(m.loader.View())
	} else {
		b.WriteString(components.SubTitleStyle.Render("Нажато:"))
		b.WriteString(" " + strconv.FormatInt(m.data.Count, 10) + " ")

		b.WriteString(components.SubTitleStyle.Render("Последнее:"))
		b.WriteString(" " + m.data.LastPress.Format(time.DateTime) + "\n")
	}

	b.WriteString("\n")

	b.WriteString(components.HelpStyle.Render("Для обновления нажми ctrl+r\n"))
	b.WriteString(components.HelpStyle.Render("Для нажатия на кнопку нажми ctrl+v"))

	return b.String()
}
