package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/model"
)

type LoginView struct {
	shared    Shared
	form      *MultiInput
	lastError error
}

func NewLoginView(shared Shared) *LoginView {
	v := &LoginView{
		shared: shared,
		form: NewMultiInput([]MultiInputField{
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
	}

	return v
}

func (v *LoginView) Init() tea.Cmd {
	return v.form.Init()
}

func (v *LoginView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return v, tea.Quit

		case "enter":
			if v.form.Finished() {
				values := v.form.Values()
				v.lastError = v.shared.storage.SetConnectionData(v.shared.ctx, model.Connection{
					Addr:     values[0],
					Login:    values[1],
					Password: values[2],
				})
				if v.lastError == nil {
					return NewStorageDebugView(v.shared), nil
				}
			}
		}
	}

	_, cmd := v.form.Update(msg)

	return v, cmd
}

func (v *LoginView) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Авторизация:"))
	b.WriteString("\n\n")
	b.WriteString(renderError(v.lastError))
	b.WriteString(v.form.View())

	return b.String()
}
