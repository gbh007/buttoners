package view

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/model"
)

type StorageDebugEvent struct {
	isInit bool
	data   model.Connection
	err    error
}

type StorageDebugView struct {
	shared    Shared
	data      model.Connection
	loader    LoaderView
	lastError error
}

func NewStorageDebugView(shared Shared) StorageDebugView {
	return StorageDebugView{
		shared: shared,
		loader: NewLoaderView("Загрузка"),
	}
}

func (v StorageDebugView) refreshData() (StorageDebugView, tea.Cmd) {
	var cmd tea.Cmd
	v.loader, cmd = v.loader.Activate()

	cmd = tea.Batch(cmd, func() tea.Msg {
		var event StorageDebugEvent
		event.data, event.err = v.shared.storage.GetConnectionData(v.shared.ctx)
		return event
	})

	return v, cmd
}

func (v StorageDebugView) Init() tea.Cmd {
	return func() tea.Msg {
		return StorageDebugEvent{
			isInit: true,
		}
	}
}

func (v StorageDebugView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return v, tea.Quit
		case "ctrl+r":
			return v.refreshData()
		}
	case StorageDebugEvent:
		if msg.isInit {
			return v.refreshData()
		}

		var cmd tea.Cmd
		v.loader, cmd = v.loader.Deactivate()
		v.data = msg.data
		v.lastError = msg.err
		return v, cmd
	}

	var cmd tea.Cmd
	v.loader, cmd = v.loader.Update(msg)

	return v, cmd
}

func (v StorageDebugView) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Данные хранилища:"))
	b.WriteString("\n\n")

	b.WriteString(renderError(v.lastError))

	if v.loader.IsActivate() {
		b.WriteString(v.loader.View())
	} else {
		b.WriteString(subTitleStyle.Render("Хост:"))
		b.WriteString(" " + v.data.Addr + "\n")

		b.WriteString(subTitleStyle.Render("Логин:"))
		b.WriteString(" " + v.data.Login + "\n")

		b.WriteString(subTitleStyle.Render("Пароль:"))
		b.WriteString(" " + v.data.Password + "\n")
	}

	b.WriteString("\n")

	b.WriteString(helpStyle.Render("Для обновления нажми ctrl+r"))

	return b.String()
}
