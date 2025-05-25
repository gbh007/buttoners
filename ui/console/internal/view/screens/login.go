package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/model"
	"github.com/gbh007/buttoners/ui/console/internal/view/components"
)

type LoginEvent struct {
	err error
}

type Login struct {
	shared    SharedState
	form      components.MultiInput
	lastError error
	loader    components.Loader
}

func NewLogin(shared SharedState) Login {
	return Login{
		shared: shared,
		form: components.NewMultiInput([]components.MultiInputField{
			{
				Name:         "Хост",
				Placeholder:  "host",
				DefaultValue: "localhost:14281",
				CharLimit:    250,
			},
			{
				Name:        "Логин",
				Placeholder: "login",
				CharLimit:   50,
			},
			{
				Name:        "Пароль",
				Placeholder: "pass",
				CharLimit:   50,
				Mode:        textinput.EchoPassword,
			},
		}),
		loader: components.NewLoader("Сохранение"),
	}
}

func (m Login) Init() tea.Cmd {
	return m.form.Init()
}

func (m Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if m.form.Finished() {
				values := m.form.Values()
				var cmd tea.Cmd
				m.loader, cmd = m.loader.Activate()
				return m, tea.Batch(
					cmd,
					func() tea.Msg {
						err := m.shared.Storage.SetConnectionData(m.shared.Ctx, model.Connection{
							Addr:     values[0],
							Login:    values[1],
							Password: values[2],
						})

						return LoginEvent{
							err: err,
						}
					},
				)
			}
		}
	case LoginEvent:
		if msg.err != nil {
			var cmd tea.Cmd
			m.loader, cmd = m.loader.Deactivate()
			m.lastError = msg.err
			return m, cmd
		}

		s := NewStorageDebug(m.shared)
		return s, s.Init()
	}

	var formCmd, loaderCmd tea.Cmd

	m.form, formCmd = m.form.Update(msg)
	m.loader, loaderCmd = m.loader.Update(msg)

	return m, tea.Batch(formCmd, loaderCmd)
}

func (m Login) View() string {
	var b strings.Builder

	b.WriteString(components.TitleStyle.Render("Авторизация:"))
	b.WriteString("\n\n")
	b.WriteString(components.RenderError(m.lastError))

	if m.loader.IsActivate() {
		b.WriteString(m.loader.View())
	} else {
		b.WriteString(m.form.View())
	}

	return b.String()
}
