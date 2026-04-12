package tui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
	"github.com/IvanJijon/pimpomodoro/task"
)

// handleKey dispatches key messages to the appropriate handler based on the current view mode.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Dismiss visual alert on any keypress.
	m.alerting = false
	if msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}
	if msg.String() == "q" && m.viewMode != ModeTaskAdd {
		if m.confirmEnabled {
			m.viewMode = ModeQuitConfirm
			m.running = false
			return m, nil
		}
		return m, tea.Quit
	}
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
	case ModeTaskList:
		return m.handleKeyTaskList(msg)
	case ModeTaskAdd:
		return m.handleKeyTaskAdd(msg)
	case ModeHelp:
		return m.handleKeyHelp(msg)
	default:
		return m, nil
	}
}

// handleKeyNormal processes key presses in normal mode: start/pause, reset, skip, previous, help, and quit.
func (m Model) handleKeyNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
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
	case "t":
		m.viewMode = ModeTaskList
		return m, nil
	case "?":
		m.viewMode = ModeHelp
		return m, nil
	}
	return m, nil
}

// handleKeyTaskList processes key presses in task list mode: navigation, selection, and exit.
func (m Model) handleKeyTaskList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "t", tea.KeyEsc.String():
		m.viewMode = ModeNormal
		return m, nil
	case "a":
		m.viewMode = ModeTaskAdd
		m.taskNameInput.Focus()
		return m, nil
	case "up", "k":
		if m.taskCursor > 0 {
			m.taskCursor--
		}
		return m, nil
	case "down", "j":
		if m.taskCursor < len(m.taskList.Tasks())-1 {
			m.taskCursor++
		}
		return m, nil
	case "d":
		if m.taskList.Len() == 0 {
			return m, nil
		}
		m.taskList.MarkTaskDone(m.taskList.Tasks()[m.taskCursor])
		return m, nil
	case "x":
		if m.taskList.Len() == 0 {
			return m, nil
		}
		m.taskList.Remove(m.taskList.Tasks()[m.taskCursor])
		if m.taskCursor >= m.taskList.Len() && m.taskList.Len() > 0 {
			m.taskCursor = m.taskList.Len() - 1
		}
		return m, nil
	case tea.KeyEnter.String():
		if m.taskList.Len() == 0 {
			return m, nil
		}
		m.taskList.SelectWIP(m.taskList.Tasks()[m.taskCursor])
	}
	return m, nil
}

// handleKeyTaskAdd processes key presses in task add mode.
func (m Model) handleKeyTaskAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case tea.KeyEsc.String():
		m.taskNameInput.SetValue("")
		m.taskEstimateInput.SetValue("")
		m.viewMode = ModeTaskList
		return m, nil
	case tea.KeyEnter.String():
		name := m.taskNameInput.Value()
		if name == "" {
			return m, nil
		}
		estimate, err := strconv.Atoi(m.taskEstimateInput.Value())
		if err != nil || estimate <= 0 {
			estimate = 1
		}
		m.taskList.Add(task.NewTask(name, estimate))
		m.taskNameInput.SetValue("")
		m.taskEstimateInput.SetValue("")
		m.viewMode = ModeTaskList
		return m, nil
	case tea.KeyTab.String():
		if m.taskNameInput.Focused() {
			m.taskNameInput.Blur()
			m.taskEstimateInput.Focus()
		} else {
			m.taskEstimateInput.Blur()
			m.taskNameInput.Focus()
		}
		return m, nil
	}
	var cmd tea.Cmd
	if m.taskNameInput.Focused() {
		m.taskNameInput, cmd = m.taskNameInput.Update(msg)
	} else {
		m.taskEstimateInput, cmd = m.taskEstimateInput.Update(msg)
	}
	return m, cmd
}

// handleKeyQuitConfirm processes key presses in the quit confirmation dialog.
func (m Model) handleKeyQuitConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		return m, tea.Quit
	case "n", tea.KeyEsc.String():
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
	case "n", tea.KeyEsc.String():
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
	case "n", tea.KeyEsc.String():
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
	case "n", tea.KeyEsc.String():
		m.viewMode = ModeNormal
		m.running = true
		m.tickID++
		return m, tea.Batch(tickCmd(m.tickID), m.spinner.Tick)
	}
	return m, nil
}

// handleKeyHelp processes key presses in help mode.
func (m Model) handleKeyHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "?" || msg.String() == tea.KeyEsc.String() {
		m.viewMode = ModeNormal
	}
	return m, nil
}
