package tui

import tea "github.com/charmbracelet/bubbletea"

// handleBlink toggles the blink state for the visual alert.
func (m Model) handleBlink() (tea.Model, tea.Cmd) {
	if !m.alerting {
		return m, nil
	}
	m.blinkState = !m.blinkState
	return m, blinkCmd()
}
