package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/spinner"
)

type Loader struct {
	text    string
	active  bool
	spinner spinner.Model
}

func NewLoader(text string) Loader {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))

	return Loader{
		text:    text,
		spinner: s,
	}
}

func (m Loader) Activate() (Loader, tea.Cmd) {
	m.active = true
	return m, m.spinner.Tick
}

func (m Loader) IsActivate() bool {
	return m.active
}

func (m Loader) Deactivate() (Loader, tea.Cmd) {
	m.active = false
	return m, nil
}

func (m Loader) Update(msg tea.Msg) (Loader, tea.Cmd) {
	if !m.active {
		return m, nil
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	return m, cmd

}

func (m Loader) View() string {
	if !m.active {
		return ""
	}

	var b strings.Builder

	b.WriteString(m.spinner.View())
	b.WriteString(" " + m.text + "\n")

	return b.String()
}
