package main

import (
	"BackItUp/IO"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type bModel struct {
	status      string
	done        bool
	failed      bool
	err         error
	filesDone   int
	filesTotal  int
	progressBar progress.Model
	spinner     spinner.Model
	progressCh  chan IO.ProgressMsg
}

type backupDoneMsg struct {
	err error
}

func listenProgress(progressCh <-chan IO.ProgressMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-progressCh
		if !ok {
			// channel closed, stop listening
			return nil
		}
		return msg
	}
}

// startBackup now accepts a progress channel to send progress updates
func startBackup(cfg IO.Config, progressCh chan IO.ProgressMsg) tea.Cmd {
	return func() tea.Msg {
		now := time.Now()
		timestamp := now.Format("2006-01-02_15-04-05")
		zipName := fmt.Sprintf("backitup_%s.zip", timestamp)

		err := IO.ZipWithExtensions(zipName, cfg.Extensions, progressCh)
		return backupDoneMsg{err: err}
	}
}

func newBModel(cfg IO.Config) (bModel, tea.Cmd) {
	prog := progress.New(progress.WithDefaultGradient())
	spin := spinner.New()
	spin.Spinner = spinner.Line
	progressCh := make(chan IO.ProgressMsg)
	m := bModel{
		status:      "🔄 Starting backup...",
		progressBar: prog,
		spinner:     spin,
		progressCh:  progressCh,
	}

	cmd := tea.Batch(
		startBackup(cfg, progressCh),
		listenProgress(progressCh),
	)

	return m, cmd
}

func (m bModel) Init() tea.Cmd {
	return nil
}

func (m bModel) Update(msg tea.Msg) (bModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case IO.ProgressMsg:
		m.filesDone = msg.Done
		m.filesTotal = msg.Total
		m.status = fmt.Sprintf("Zipping file %d of %d...", m.filesDone, m.filesTotal)
		return m, listenProgress(m.progressCh)

	case backupDoneMsg:
		if msg.err != nil {
			m.status = "❌ Backup failed"
			m.failed = true
			m.err = msg.err
		} else {
			m.status = "✅ Backup complete"
			m.done = true
		}

	case tea.KeyMsg:
		if m.done || m.failed {
			if msg.String() == "q" || msg.String() == "enter" {
				return m, func() tea.Msg { return backToMenuMsg{} }
			}
		}
	}

	m.progressBar.Update(msg)
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m bModel) View() string {
	s := "Back It Up | Execution\n\n"
	s += m.spinner.View() + " " + m.status + "\n\n"

	if m.filesTotal > 0 {
		percent := float64(m.filesDone) / float64(m.filesTotal)
		s += m.progressBar.ViewAs(percent) + "\n\n"
	}

	if m.err != nil {
		s += fmt.Sprintf("Error: %v\n", m.err)
	}

	if m.done || m.failed {
		s += "\nPress q or Enter to return to main menu.\n"
	}

	return s
}
