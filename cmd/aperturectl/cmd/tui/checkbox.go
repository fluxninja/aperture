package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type CheckBoxModel struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	Selected map[int]struct{} // which to-do items are selected
	msg      string
}

func InitialCheckboxModel(options []string, msg string) *CheckBoxModel {
	return &CheckBoxModel{
		// Our to-do list is a grocery list
		choices: options,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		Selected: make(map[int]struct{}),

		msg: msg,
	}
}

func (m CheckBoxModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m CheckBoxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "c", "C", "q", "Q":
			return m, tea.Quit

		// The "d" key to de-select all
		case "d", "D":
			for option := range m.choices {
				delete(m.Selected, option)
			}

		// The "s" key to select all
		case "s", "S":
			for option := range m.choices {
				m.Selected[option] = struct{}{}
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.Selected[m.cursor]
			if ok {
				delete(m.Selected, m.cursor)
			} else {
				m.Selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m CheckBoxModel) View() string {
	// The header
	s := fmt.Sprintf("\n%s\n\n", m.msg)

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.Selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress s to select all.\n"
	s += "Press d to de-select all.\n"
	s += "Press c to confirm.\n"

	// Send the UI for rendering
	return s
}
