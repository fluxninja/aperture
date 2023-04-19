package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/utils/pointer"
)

type radioButtonModel struct {
	// The currently selected option
	Selected *int
	options  []string
	msg      string
}

type optionMsg int

// InitialRadioButtonModel returns a new model with the given options and message.
func InitialRadioButtonModel(options []string, msg string) *radioButtonModel {
	return &radioButtonModel{
		options:  options,
		Selected: pointer.Int(0),
		msg:      msg,
	}
}

// Init is called when the model is first initialized.
func (m radioButtonModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received.
func (m radioButtonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			*m.Selected--
			if *m.Selected < 0 {
				*m.Selected = len(m.options) - 1
			}
		case "down":
			*m.Selected++
			if *m.Selected >= len(m.options) {
				*m.Selected = 0
			}
		case "enter":
			return m, tea.Quit
		}
	case optionMsg:
		*m.Selected = int(msg)
	}

	return m, nil
}

// View renders the UI.
func (m radioButtonModel) View() string {
	// The header
	s := fmt.Sprintf("\n%s\n\n", m.msg)

	// Iterate over our choices
	for i, choice := range m.options {

		checked := " "
		if i == *m.Selected {
			checked = "x"
		}
		// Render the row
		s += fmt.Sprintf("(%s) %s\n", checked, choice)
	}

	// The footer
	s += "\nPress 'enter' to confirm.\n"

	// Send the UI for rendering
	return s
}
