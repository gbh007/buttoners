package view

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbh007/buttoners/ui/console/internal/view/screens"
)

func Run() error {
	shared := &screens.SharedState{
		Ctx: context.Background(),
	}

	startView := screens.NewMenu(shared)
	p := tea.NewProgram(
		startView,
		tea.WithContext(shared.Ctx),
		tea.WithAltScreen(),
	)
	_, err := p.Run()
	if err != nil {
		return err
	}

	return nil
}
