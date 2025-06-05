package screens

import (
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/view/components"
)

type MainButtonEvent struct {
	isInit bool
	data   struct {
		Count     int64
		LastPress time.Time
	}
	err error
}

type MainButton struct {
	shared *SharedState
	data   struct {
		Count     int64
		LastPress time.Time
	}
	form      components.MultiInput
	loader    components.Loader
	lastError error
}

func NewMainButton(shared *SharedState) MainButton {
	return MainButton{
		shared: shared,
		loader: components.NewLoader("Загрузка"),
		form: components.NewMultiInput([]components.MultiInputField{
			{
				Name:         "Продолжительность [1;60]",
				DefaultValue: "3",
				CharLimit:    10,
			},
			{
				Name:         "Шанс провала [0;100]",
				DefaultValue: "20",
				CharLimit:    10,
			},
		}, "Нажать"),
	}
}

func (m MainButton) refreshData() (MainButton, tea.Cmd) {
	var cmd tea.Cmd
	m.loader, cmd = m.loader.Activate()

	cmd = tea.Batch(cmd, func() tea.Msg {
		var event MainButtonEvent
		event.data.Count, event.data.LastPress, event.err = m.shared.GateClient.Activity(m.shared.Ctx, m.shared.GateToken)
		return event
	})

	return m, cmd
}

func (m MainButton) pressButton(duration int64, chance int64) (MainButton, tea.Cmd) {
	var cmd tea.Cmd
	m.loader, cmd = m.loader.Activate()

	cmd = tea.Batch(cmd, func() tea.Msg {
		var event MainButtonEvent
		err := m.shared.GateClient.ButtonClick(m.shared.Ctx, m.shared.GateToken, duration, chance)
		if err != nil {
			return MainButtonEvent{
				err: err,
			}
		}

		event.data.Count, event.data.LastPress, event.err = m.shared.GateClient.Activity(m.shared.Ctx, m.shared.GateToken)
		return event
	})

	return m, cmd
}

func (m MainButton) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return MainButtonEvent{
				isInit: true,
			}
		},
		m.form.Init(),
	)
}

func (m MainButton) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			nextScreen := NewMenu(m.shared)
			return nextScreen, nextScreen.Init()

		case "enter":
			if m.form.Finished() {
				values := m.form.Values()

				duration, err := strconv.ParseInt(values[0], 10, 64)
				if err != nil {
					return m, func() tea.Msg {
						return MainButtonEvent{
							err: err,
						}
					}
				}

				chance, err := strconv.ParseInt(values[1], 10, 64)
				if err != nil {
					return m, func() tea.Msg {
						return MainButtonEvent{
							err: err,
						}
					}
				}

				return m.pressButton(duration, chance)
			}
		}
	case MainButtonEvent:
		if msg.isInit {
			return m.refreshData()
		}

		var cmd tea.Cmd
		m.loader, cmd = m.loader.Deactivate()
		m.data = msg.data
		m.lastError = msg.err
		return m, cmd
	}

	var loaderCmd tea.Cmd
	m.loader, loaderCmd = m.loader.Update(msg)

	var formCmd tea.Cmd
	m.form, formCmd = m.form.Update(msg)

	return m, tea.Batch(loaderCmd, formCmd)
}

func (m MainButton) View() string {
	var b strings.Builder

	b.WriteString(components.TitleStyle.Render("Нажимай кнопку всласть:"))
	b.WriteString("\n\n")

	b.WriteString(components.RenderError(m.lastError))

	if m.loader.IsActivate() {
		b.WriteString(m.loader.View())
	} else {
		b.WriteString(components.SubTitleStyle.Render("Запросов:"))
		b.WriteString(" " + strconv.FormatInt(m.data.Count, 10) + " ")

		b.WriteString(components.SubTitleStyle.Render("Последний:"))
		b.WriteString(" " + m.data.LastPress.In(time.Local).Format(time.DateTime) + "\n")
	}

	b.WriteString("\n")

	b.WriteString(m.form.View())

	return b.String()
}
