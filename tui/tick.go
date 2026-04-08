package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg is sent by the Bubble Tea runtime every second via tea.Tick.
type TickMsg struct {
	id int
}

func tickCmd(id int) tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return TickMsg{id: id}
	})
}
