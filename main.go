package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type backToMenuMsg struct{}

type viewState int

const (
	viewMenu = iota
	viewConfigEditor
	viewExecution
)

type model struct {
	currentView viewState
	choices     []string         // items in the main menu
	cursor      int              // which menu item our cursor is pointing at
	selected    map[int]struct{} // which menu items are selected
	Config      cModel           // ← THIS is your cModel
}

func initialModel() model {
	return model{
		choices: []string{"Start Backup", "Edit Config"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
		Config:   configModel(),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global check for backToMenuMsg
	switch msg.(type) {
	case backToMenuMsg:
		m.currentView = viewMenu
		return m, nil
	}

	// View-specific update
	switch m.currentView {
	case viewMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			case "down":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
				return m, nil
			case "enter", " ":
				switch m.cursor {
				case 0:
					m.currentView = viewExecution
				case 1:
					m.currentView = viewConfigEditor
				}
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}

	case viewConfigEditor:
		var cmd tea.Cmd
		m.Config, cmd = m.Config.Update(msg)
		return m, cmd

	case viewExecution:
		// TODO: Implement execution.Update logic here
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	switch m.currentView {
	case viewMenu:
		// existing menu rendering code here
		s := "Back It Up v1.0.0\n\n"
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
		s += "\nPress q to quit.\n"
		return s

	case viewConfigEditor:
		// delegate to Config's View
		return m.Config.View()

	case viewExecution:
		// you can add execution view here later
		return "Execution view (not implemented yet)"

	default:
		return "Unknown view"
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
