package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
)

// handleKey dispatches key messages to the appropriate handler based on the current view mode.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Dismiss visual alert on any keypress.
	m.alerting = false
	switch m.viewMode {
	case ModeNormal:
		return m.handleKeyNormal(msg)
	case ModeQuitConfirm:
		return m.handleKeyQuitConfirm(msg)
	case ModeSkipConfirm:
		return m.handleKeySkipConfirm(msg)
	case ModeResetConfirm:
		return m.handleKeyResetConfirm(msg)
	case ModePreviousConfirm:
		return m.handleKeyPreviousConfirm(msg)
	case ModeHelp:
		return m.handleKeyHelp(msg)
	default:
		return m, nil
	}
}

// handleKeyNormal processes key presses in normal mode: start/pause, reset, skip, previous, help, and quit.
func (m Model) handleKeyNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl-c":
		return m, tea.Quit
	case "q":
		if m.confirmEnabled {
			m.viewMode = ModeQuitConfirm
			m.running = false
			return m, nil
		}
		return m, tea.Quit
	case "s":
		if m.running {
			m.running = false
			return m, nil
		}
		if m.session.CurrentPhase == session.Idle {
			m.session.NextPhase()
			m.remainingTime = m.session.PhaseDuration()
		}
		m.running = true
		m.tickID++
		return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
	case "r":
		if m.confirmEnabled {
			m.viewMode = ModeResetConfirm
			m.running = false
			return m, nil
		}
		m.remainingTime = m.session.PhaseDuration()
		m.running = false
		return m, nil
	case "n":
		if m.session.CurrentPhase == session.Idle {
			return m, nil
		}
		if m.confirmEnabled {
			m.viewMode = ModeSkipConfirm
			m.running = false
			return m, nil
		}
		m.session.NextPhase()
		m.remainingTime = m.session.PhaseDuration()
		m.running = false
		return m, nil
	case "b":
		if m.session.CurrentPhase == session.Idle {
			return m, nil
		}
		if m.confirmEnabled {
			m.viewMode = ModePreviousConfirm
			m.running = false
			return m, nil
		}
		m.session.PreviousPhase()
		m.remainingTime = m.session.PhaseDuration()
		m.running = false
		return m, nil
	case "?":
		m.viewMode = ModeHelp
		return m, nil
	}
	return m, nil
}

// handleKeyQuitConfirm processes key presses in the quit confirmation dialog.
func (m Model) handleKeyQuitConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		return m, tea.Quit
	case "n":
		m.viewMode = ModeNormal
		if m.session.CurrentPhase != session.Idle {
			m.running = true
			m.tickID++
			return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
		}
		return m, nil
	}
	return m, nil
}

// handleKeySkipConfirm processes key presses in the skip phase confirmation dialog.
func (m Model) handleKeySkipConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.session.NextPhase()
		m.remainingTime = m.session.PhaseDuration()
		m.viewMode = ModeNormal
	case "n":
		m.viewMode = ModeNormal
		m.running = true
		m.tickID++
		return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
	}
	return m, nil
}

// handleKeyResetConfirm processes key presses in the reset confirmation dialog.
func (m Model) handleKeyResetConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.remainingTime = m.session.PhaseDuration()
		m.viewMode = ModeNormal
		return m, nil
	case "n":
		m.viewMode = ModeNormal
		m.running = true
		m.tickID++
		return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
	}
	return m, nil
}

// handleKeyPreviousConfirm processes key presses in the previous phase confirmation dialog.
func (m Model) handleKeyPreviousConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.session.PreviousPhase()
		m.remainingTime = m.session.PhaseDuration()
		m.viewMode = ModeNormal
		return m, nil
	case "n":
		m.viewMode = ModeNormal
		m.running = true
		m.tickID++
		return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
	}
	return m, nil
}

// handleKeyHelp processes key presses in help mode.
func (m Model) handleKeyHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "?" {
		m.viewMode = ModeNormal
	}
	return m, nil
}
