package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type cModel struct {
	choices  []string         // items in the main menu
	cursor   int              // which menu item our cursor is pointing at
	selected map[int]struct{} // which menu items are selected
	input    textinput.Model  // file extension input
}

func configModel() cModel {
	ti := textinput.New()
	ti.Placeholder = "Enter file extensions (comma separated)"
	ti.Focus()        // Focus when this view is active
	ti.CharLimit = 64 // optional limit
	ti.Width = 40     // width in characters

	return cModel{
		choices:  []string{"Leave", "Leave"},
		selected: make(map[int]struct{}),
		input:    ti,
	}
}

func (m cModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m cModel) Update(msg tea.Msg) (cModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			// You could handle "submit" here, e.g., parse input.Value()
			// For now, just return the model
			return m, nil
		}
	}

	// Let textinput handle the rest of the messages
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m cModel) View() string {
	s := "Back It Up | Config Editor\n\n"

	s += "Enter file extensions (comma separated):\n"
	s += m.input.View() + "\n\n"

	s += "Press Enter to submit, q to quit\n"

	return s
}
