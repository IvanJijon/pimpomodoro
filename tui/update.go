package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
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

// Update handles messages and updates the model accordingly.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.running {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	case TickMsg:
		if !m.running || msg.id != m.tickID {
			return m, nil
		}

		// running
		if m.remainingTime <= 0 {
			m.session.NextPhase()
			m.remainingTime = m.session.PhaseDuration()
			m.running = false
			return m, nil
		}
		m.remainingTime -= time.Second
		return m, tickCmd(m.tickID)
	case tea.KeyMsg:
		if m.showSkipConfirm {
			switch msg.String() {
			case "y":
				m.session.NextPhase()
				m.remainingTime = m.session.PhaseDuration()
				m.showSkipConfirm = false
			case "x":
				m.showSkipConfirm = false
				m.running = true
				m.tickID++
				return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		case "s":
			if m.running {
				return m, nil
			}
			if m.session.CurrentPhase == session.Idle {
				m.session.NextPhase()
				m.remainingTime = m.session.PhaseDuration()
			}
			m.running = true
			m.tickID++
			return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
		case "p":
			m.running = false
			return m, nil
		case "r":
			m.remainingTime = m.session.PhaseDuration()
			m.running = false
			return m, nil
		case "n":
			if m.session.CurrentPhase == session.Idle {
				return m, nil
			}
			m.showSkipConfirm = true
			m.running = false
			return m, nil
		case "b":
			if m.session.CurrentPhase == session.Idle {
				return m, nil
			}
			m.session.PreviousPhase()
			m.remainingTime = m.session.PhaseDuration()
			m.running = false
			return m, nil
		case "?":
			if m.showHelp {
				m.showHelp = false
			} else {
				m.showHelp = true
			}
			return m, nil
		}
	}
	return m, nil
}
