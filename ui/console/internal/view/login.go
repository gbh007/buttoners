package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/model"
)

type LoginViewEvent struct {
	err error
}

type LoginView struct {
	shared    Shared
	form      MultiInput
	lastError error
	loader    LoaderView
}

func NewLoginView(shared Shared) LoginView {
	v := LoginView{
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
		loader: NewLoaderView("Сохранение"),
	}

	return v
}

func (v LoginView) Init() tea.Cmd {
	return v.form.Init()
}

func (v LoginView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return v, tea.Quit

		case "enter":
			if v.form.Finished() {
				values := v.form.Values()
				var cmd tea.Cmd
				v.loader, cmd = v.loader.Activate()
				return v, tea.Batch(
					cmd,
					func() tea.Msg {
						err := v.shared.storage.SetConnectionData(v.shared.ctx, model.Connection{
							Addr:     values[0],
							Login:    values[1],
							Password: values[2],
						})

						return LoginViewEvent{
							err: err,
						}
					},
				)
			}
		}
	case LoginViewEvent:
		if msg.err != nil {
			var cmd tea.Cmd
			v.loader, cmd = v.loader.Deactivate()
			v.lastError = msg.err
			return v, cmd
		}

		s := NewStorageDebugView(v.shared)
		return s, s.Init()
	}

	var formCmd, loaderCmd tea.Cmd

	v.form, formCmd = v.form.Update(msg)
	v.loader, loaderCmd = v.loader.Update(msg)

	return v, tea.Batch(formCmd, loaderCmd)
}

func (v LoginView) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Авторизация:"))
	b.WriteString("\n\n")
	b.WriteString(renderError(v.lastError))

	if v.loader.IsActivate() {
		b.WriteString(v.loader.View())
	} else {
		b.WriteString(v.form.View())
	}

	return b.String()
}
