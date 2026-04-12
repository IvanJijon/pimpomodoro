package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
)

// handleWindowSize stores the terminal dimensions when the window is resized.
func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	return m, nil
}

// handleSpinnerTick forwards spinner ticks while the timer is running.
func (m Model) handleSpinnerTick(msg spinner.TickMsg) (tea.Model, tea.Cmd) {
	if m.running {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

// handleTick processes timer ticks: decrements remaining time or triggers phase transition on expiry.
func (m Model) handleTick(msg TickMsg) (tea.Model, tea.Cmd) {
	if !m.running || msg.id != m.tickID {
		return m, nil
	}
	if m.remainingTime <= 0 {
		if m.session.CurrentPhase == session.Work {
			if tsk := m.taskList.CurrentWIP(); tsk != nil {
				tsk.IncreaseActualPomos()
			}
		}
		notifyMsg := phaseCompleteMessage(m.session.CurrentPhase)
		if m.visualAlert {
			m.alertColor = phaseColor(m)
			m.alerting = true
		}
		m.session.NextPhase()
		m.remainingTime = m.session.PhaseDuration()
		m.running = false
		m.callbacks.PlayAlarm()
		m.callbacks.SendNotify("Pimpomodoro", notifyMsg)
		if m.alerting {
			return m, blinkCmd()
		}
		return m, nil
	}
	m.remainingTime -= time.Second
	return m, tickCmd(m.tickID)
}
