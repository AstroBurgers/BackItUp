package tui

import (
	"BackItUp/io"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type configModel struct {
	input textinput.Model
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

func newConfigModel() configModel {
	ti := textinput.New()
	ti.Placeholder = "Enter file extensions (comma separated)"
	ti.Focus()        // Focus when this view is active
	ti.CharLimit = 64 // limit
	ti.Width = 40     // width in characters

	return configModel{
		input: ti,
	}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (configModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, func() tea.Msg { return backToMenuMsg{} }
		case "enter":
			extensions := parseExtensions(m.input.Value())
			cfg := io.Config{Extensions: extensions}
			err := io.SaveConfig(cfg)
			if err != nil {
				return m, nil
			}
			m.input.Reset()
			return m, tea.Tick(time.Second/2, func(time.Time) tea.Msg {
				return backToMenuMsg{}
			})
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m configModel) View() string {
	s := "Back It Up | Config Editor\n\n"
	s += "Enter file extensions (comma separated):\n"
	s += m.input.View() + "\n\n"
	s += "Press Enter to submit, q to quit\n"
	return s
}
