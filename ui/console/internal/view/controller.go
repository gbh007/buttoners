package view

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/storage"
	"github.com/gbh007/buttoners/ui/console/internal/view/screens"
)

func Run() error {
	shared := screens.SharedState{
		Ctx:     context.Background(),
		Storage: storage.New(),
	}

	startView := screens.NewLogin(shared)
	p := tea.NewProgram(startView)
	_, err := p.Run()
	if err != nil {
		return err
	}

	return nil
}

/* Для копипасты экранов

type ScreenTemplate struct{
	shared Shared
}


func NewScreenTemplate(shared Shared) ScreenTemplate {
	return ScreenTemplate{
		shared: shared,
	}
}

func (m ScreenTemplate) Init() tea.Cmd {
	return nil
}

func (m ScreenTemplate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ScreenTemplate) View() string {
	return ""
}

*/

/* Для копипасты компонентов

type ComponentTemplate struct{
	shared Shared
}


func NewComponentTemplate(shared Shared) ComponentTemplate {
	return ComponentTemplate{
		shared: shared,
	}
}

func (m ComponentTemplate) Init() tea.Cmd {
	return nil
}

func (m ComponentTemplate) Update(msg tea.Msg) (ComponentTemplate, tea.Cmd) {
	return m, nil
}

func (m ComponentTemplate) View() string {
	return ""
}

*/
