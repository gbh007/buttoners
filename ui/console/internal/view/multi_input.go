package view

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

func NewMultiInput(fields []MultiInputField) *MultiInput {
	v := &MultiInput{
		inputs:        make([]textinput.Model, len(fields)),
		fields:        fields,
		focusedButton: focusedStyle.Render("[ Готово ]"),
		blurredButton: fmt.Sprintf("[ %s ]", blurredStyle.Render("Готово")),
	}

	var t textinput.Model
	for i, field := range fields {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = field.CharLimit
		t.Placeholder = field.Placeholder
		t.SetValue(field.DefaultValue)
		t.EchoMode = field.Mode
		t.EchoCharacter = '•'

		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		v.inputs[i] = t
	}

	return v
}

func (v *MultiInput) Init() tea.Cmd {
	return textinput.Blink
}

func (v *MultiInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				v.focusIndex--
			} else {
				v.focusIndex++
			}

			if v.focusIndex > len(v.inputs) {
				v.focusIndex = 0
			} else if v.focusIndex < 0 {
				v.focusIndex = len(v.inputs)
			}

			cmds := make([]tea.Cmd, len(v.inputs))
			for i := 0; i <= len(v.inputs)-1; i++ {
				if i == v.focusIndex {
					cmds[i] = v.inputs[i].Focus()
					v.inputs[i].PromptStyle = focusedStyle
					v.inputs[i].TextStyle = focusedStyle
					continue
				}

				v.inputs[i].Blur()
				v.inputs[i].PromptStyle = defaultStyle
				v.inputs[i].TextStyle = defaultStyle
			}

			return v, tea.Batch(cmds...)
		}
	}

	cmds := make([]tea.Cmd, len(v.inputs))

	for i := range v.inputs {
		v.inputs[i], cmds[i] = v.inputs[i].Update(msg)
	}

	return v, tea.Batch(cmds...)

}

func (v *MultiInput) View() string {
	var b strings.Builder

	for i := range v.inputs {
		if i != 0 {
			b.WriteRune('\n')
		}

		fmt.Fprintf(&b, "%s\n", subTitleStyle.Render(v.fields[i].Name+":"))
		b.WriteString(v.inputs[i].View())
	}

	button := v.blurredButton
	if v.focusIndex == len(v.inputs) {
		button = v.focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n", button)

	return b.String()
}

func (v *MultiInput) Values() []string {
	res := make([]string, len(v.inputs))

	for i := range v.inputs {
		res[i] = v.inputs[i].Value()
	}

	return res
}

func (v *MultiInput) Finished() bool {
	return v.focusIndex == len(v.inputs)
}
