package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// BlinkMsg drives the visual alert blink cycle.
type BlinkMsg struct{}

func blinkCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(_ time.Time) tea.Msg {
		return BlinkMsg{}
	})
}
