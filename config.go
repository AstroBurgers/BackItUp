package main

import (
	"BackItUp/IO"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type cModel struct {
	choices  []string         // items in the main menu
	cursor   int              // which menu item our cursor is pointing at
	selected map[int]struct{} // which menu items are selected
	input    textinput.Model  // file extension input
}

func parseExtensions(input string) []string {
	parts := strings.Split(input, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
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

		case "ctrl+c", "q", "esc":
			return m, func() tea.Msg { return backToMenuMsg{} }

		case "enter":
			extensions := parseExtensions(m.input.Value())
			cfg := IO.Config{Extensions: extensions}

			err := IO.SaveConfig(cfg)
			if err != nil {
				return m, nil
			}
			m.input.Reset()

			return m, tea.Tick(time.Second/2, func(time.Time) tea.Msg {
				return backToMenuMsg{}
			})
		}
	}

	// Let textinput handle everything else
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
