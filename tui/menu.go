package tui

import (
	"BackItUp/io"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState int

const (
	viewMenu viewState = iota
	viewConfigEditor
	viewExecution
)

type Model struct {
	currentView viewState
	choices     []string
	cursor      int
	selected    map[int]struct{}
	config      configModel
	exec        backupModel
}

func InitialModel() Model {
	return Model{
		choices:  []string{"Start Backup", "Edit Config"},
		selected: make(map[int]struct{}),
		config:   newConfigModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					cfg, err := io.LoadConfig()
					if err != nil {
						fmt.Println("Failed to load config:", err)
						cfg = io.Default()
					}
					exec, cmd := newBackupModel(cfg)
					m.exec = exec
					m.currentView = viewExecution
					return m, cmd
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
		m.config, cmd = m.config.Update(msg)
		return m, cmd

	case viewExecution:
		var cmd tea.Cmd
		m.exec, cmd = m.exec.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	switch m.currentView {
	case viewMenu:
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
		return m.config.View()
	case viewExecution:
		return m.exec.View()
	default:
		return "Unknown view"
	}
}
