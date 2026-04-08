package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
)

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
		if m.remainingTime <= 0 {
			msg := phaseCompleteMessage(m.session.CurrentPhase)
			m.session.NextPhase()
			m.remainingTime = m.session.PhaseDuration()
			m.running = false
			m.callbacks.PlayAlarm()
			m.callbacks.SendNotify("Pimpomodoro", msg)
			return m, nil
		}
		m.remainingTime -= time.Second
		return m, tickCmd(m.tickID)
	case tea.KeyMsg:
		switch m.viewMode {
		case ModeNormal:
			return m.updateNormal(msg)
		case ModeQuitConfirm:
			return m.updateQuitConfirm(msg)
		case ModeSkipConfirm:
			return m.updateSkipConfirm(msg)
		case ModeResetConfirm:
			return m.updateResetConfirm(msg)
		case ModeHelp:
			return m.updateHelp(msg)
		default:
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl-c":
		return m, tea.Quit
	case "q":
		m.viewMode = ModeQuitConfirm
		m.running = false
		return m, nil
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
		m.viewMode = ModeResetConfirm
		m.running = false
		return m, nil
	case "n":
		if m.session.CurrentPhase == session.Idle {
			return m, nil
		}
		m.viewMode = ModeSkipConfirm
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
		m.viewMode = ModeHelp
		return m, nil
	}
	return m, nil
}

func (m Model) updateQuitConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		return m, tea.Quit
	case "x":
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

func (m Model) updateSkipConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.session.NextPhase()
		m.remainingTime = m.session.PhaseDuration()
		m.viewMode = ModeNormal
	case "x":
		m.viewMode = ModeNormal
		m.running = true
		m.tickID++
		return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
	}
	return m, nil
}

func (m Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "?" {
		m.viewMode = ModeNormal
	}
	return m, nil
}

func (m Model) updateResetConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.remainingTime = m.session.PhaseDuration()
		m.viewMode = ModeNormal
		return m, nil
	case "x":
		m.viewMode = ModeNormal
		m.running = true
		m.tickID++
		return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
	}
	return m, nil
}
