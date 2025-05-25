package view

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/spinner"
)

type LoaderView struct {
	text    string
	active  bool
	spinner spinner.Model
}

func NewLoaderView(text string) LoaderView {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))

	return LoaderView{
		text:    text,
		spinner: s,
	}
}

func (v LoaderView) Activate() (LoaderView, tea.Cmd) {
	v.active = true
	return v, v.spinner.Tick
}

func (v LoaderView) IsActivate() bool {
	return v.active
}

func (v LoaderView) Deactivate() (LoaderView, tea.Cmd) {
	v.active = false
	return v, nil
}

func (v LoaderView) Update(msg tea.Msg) (LoaderView, tea.Cmd) {
	if !v.active {
		return v, nil
	}

	var cmd tea.Cmd
	v.spinner, cmd = v.spinner.Update(msg)

	return v, cmd

}

func (v LoaderView) View() string {
	if !v.active {
		return ""
	}

	var b strings.Builder

	b.WriteString(v.spinner.View())
	b.WriteString(" " + v.text + "\n")

	return b.String()
}
