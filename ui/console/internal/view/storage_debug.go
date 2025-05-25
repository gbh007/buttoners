package view

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type StorageDebugView struct {
	shared Shared
}

func NewStorageDebugView(shared Shared) *StorageDebugView {
	return &StorageDebugView{
		shared: shared,
	}
}

func (v *StorageDebugView) Init() tea.Cmd {
	return nil
}

func (v *StorageDebugView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return v, tea.Quit
		}
	}

	return v, nil
}

func (v *StorageDebugView) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Данные хранилища:"))
	b.WriteString("\n\n")

	// FIXME: убрать в обновления
	data, _ := v.shared.storage.GetConnectionData(v.shared.ctx)

	b.WriteString(subTitleStyle.Render("Хост:"))
	b.WriteString(" " + data.Addr + "\n")

	b.WriteString(subTitleStyle.Render("Логин:"))
	b.WriteString(" " + data.Login + "\n")

	b.WriteString(subTitleStyle.Render("Пароль:"))
	b.WriteString(" " + data.Password + "\n")

	return b.String()
}
