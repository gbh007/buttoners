package components

import "github.com/charmbracelet/lipgloss"

var (
	FocusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("165"))
	BlurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle   = FocusedStyle
	DefaultStyle  = lipgloss.NewStyle()
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("21"))
	SubTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("117"))
	ErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("125"))
	HelpStyle     = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("66"))
)
