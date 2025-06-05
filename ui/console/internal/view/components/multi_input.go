package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MultiInputField struct {
	Name         string
	Placeholder  string
	DefaultValue string
	CharLimit    int
	Mode         textinput.EchoMode
}

type MultiInput struct {
	fields        []MultiInputField
	focusIndex    int
	inputs        []textinput.Model
	focusedButton string
	blurredButton string
}

func NewMultiInput(fields []MultiInputField, buttonText string) MultiInput {
	m := MultiInput{
		inputs:        make([]textinput.Model, len(fields)),
		fields:        fields,
		focusedButton: FocusedStyle.Render("[ " + buttonText + " ]"),
		blurredButton: fmt.Sprintf("[ %s ]", BlurredStyle.Render(buttonText)),
	}

	var t textinput.Model
	for i, field := range fields {
		t = textinput.New()
		t.Cursor.Style = CursorStyle
		t.CharLimit = field.CharLimit
		t.Placeholder = field.Placeholder
		t.SetValue(field.DefaultValue)
		t.EchoMode = field.Mode
		t.EchoCharacter = '•'

		if i == 0 {
			t.Focus()
			t.PromptStyle = FocusedStyle
			t.TextStyle = FocusedStyle
		}

		m.inputs[i] = t
	}

	return m
}

func (m MultiInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m MultiInput) Update(msg tea.Msg) (MultiInput, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = FocusedStyle
					m.inputs[i].TextStyle = FocusedStyle
					continue
				}

				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = DefaultStyle
				m.inputs[i].TextStyle = DefaultStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)

}

func (m MultiInput) View() string {
	var b strings.Builder

	for i := range m.inputs {
		if i != 0 {
			b.WriteRune('\n')
		}

		fmt.Fprintf(&b, "%s\n", SubTitleStyle.Render(m.fields[i].Name+":"))
		b.WriteString(m.inputs[i].View())
	}

	button := m.blurredButton
	if m.focusIndex == len(m.inputs) {
		button = m.focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n", button)

	return b.String()
}

func (m MultiInput) Values() []string {
	res := make([]string, len(m.inputs))

	for i := range m.inputs {
		res[i] = m.inputs[i].Value()
	}

	return res
}

func (m MultiInput) Finished() bool {
	return m.focusIndex == len(m.inputs)
}
