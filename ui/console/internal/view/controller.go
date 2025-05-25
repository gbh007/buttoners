package view

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gbh007/buttoners/ui/console/internal/storage"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("165"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle   = focusedStyle
	defaultStyle  = lipgloss.NewStyle()
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("21"))
	subTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("117"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("125"))
	helpStyle     = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("66"))
)

type Shared struct {
	storage *storage.Storage
	ctx     context.Context
}

func Run() error {
	shared := Shared{
		storage: storage.New(),
		ctx:     context.Background(),
	}

	startView := NewLoginView(shared)
	p := tea.NewProgram(startView)
	_, err := p.Run()
	if err != nil {
		return err
	}

	return nil
}

/* Для копипасты

type ViewTemplate struct{
	shared Shared
}


func NewViewTemplate(shared Shared) *ViewTemplate {
	return &ViewTemplate{
		shared: shared,
	}
}

func (v *ViewTemplate) Init() tea.Cmd {
	return nil
}

func (v *ViewTemplate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v *ViewTemplate) View() string {
	return ""
}

*/

func renderError(err error) string {
	if err == nil {
		return ""
	}

	return errorStyle.Render("Ошибка: "+err.Error()) + "\n\n"
}
