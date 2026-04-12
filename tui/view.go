package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/IvanJijon/pimpomodoro/session"
	"github.com/IvanJijon/pimpomodoro/task"
)

const viewWidth = 65

const taskLinePrefix = 18 // "> │  0/3  │ ⬜ "

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1).
			AlignHorizontal(lipgloss.Center).
			Width(viewWidth)

	centerStyle = lipgloss.NewStyle().
			Width(viewWidth).
			AlignHorizontal(lipgloss.Center)

	timerTextStyle = lipgloss.NewStyle().Bold(true)

	pausedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(yellow)

	footerStyle = lipgloss.NewStyle().
			Foreground(gray).
			Width(viewWidth).
			AlignHorizontal(lipgloss.Center)

	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Width(viewWidth - 4).
			AlignHorizontal(lipgloss.Center)
)

// View renders the current state to the terminal.
func (m Model) View() string {
	s := headerStyle.Render("\U0001f345 Pimpomodoro Timer")

	if m.viewMode == ModeHelp {
		help := "  Keybindings:\n\n"
		help += "  s  Start/pause\n"
		help += "  r  Reset current phase\n"
		help += "  n  Skip to next phase\n"
		help += "  b  Go to previous phase\n"
		help += "  t  Task list\n"
		help += "  ?  Toggle help\n"
		help += "  q  Quit"
		s += "\n" + help + "\n"
	} else if m.viewMode == ModeTaskList {
		s += "\n" + centerStyle.Render("Task List") + "\n\n"
		tasks := m.taskList.Tasks()
		if len(tasks) == 0 {
			s += centerStyle.Render(footerStyle.Render("No tasks yet")) + "\n"
		} else {
			var list string
			for i, t := range tasks {
				cursor := "  "
				if i == m.taskCursor {
					cursor = "> "
				}
				list += cursor + taskLine(t) + "\n"
			}
			s += lipgloss.NewStyle().Width(viewWidth).Render(list)
		}
		s += "\n" + footerStyle.Render("(a) add  (enter) select  (d) done") + "\n"
		s += footerStyle.Render("(x) remove  (esc) back") + "\n"
	} else if m.viewMode == ModeSwitchTaskConfirm {
		s += "\n" + centerStyle.Render("Task List") + "\n\n"
		dialog := dialogStyle.Render("Timer will be reset, switch task?\n\n(y) confirm  (n) cancel")
		s += "\n" + centerStyle.Render(dialog) + "\n"
	} else if m.viewMode == ModeTaskAdd {
		s += "\n" + centerStyle.Render("Add Task") + "\n\n"
		form := "  Name:     " + m.taskNameInput.View() + "\n"
		form += "  Estimate: " + m.taskEstimateInput.View() + "\n"
		s += lipgloss.NewStyle().Width(viewWidth).Render(form)
		s += "\n" + footerStyle.Render("(enter) confirm  (tab) next  (esc) cancel") + "\n"
	} else {
		s += "\n" + centerStyle.Render(phaseLabel(m)) + "\n"
		switch {
		case m.running:
			s += "\n" + centerStyle.Render(strings.TrimRight(m.spinner.View(), " ")+phaseTimerStyle(m).Render(formatDuration(m.remainingTime))) + "\n"
		case m.session.CurrentPhase != session.Idle:
			s += "\n" + centerStyle.Render(pausedStyle.Render("\u23f8 "+formatDuration(m.remainingTime))) + "\n"
		default:
			s += "\n" + centerStyle.Render(timerTextStyle.Render(formatDuration(m.remainingTime))) + "\n"
		}

		if wip := m.taskList.CurrentWIP(); wip != nil {
			s += "\n" + centerStyle.Render(wipLabel(wip)) + "\n"
		}

		switch m.viewMode {
		case ModeSkipConfirm:
			dialog := dialogStyle.Render("Skip to next phase?\n\n(y) confirm  (n) cancel")
			s += "\n" + centerStyle.Render(dialog) + "\n"
		case ModeQuitConfirm:
			dialog := dialogStyle.Render("Quit?\n\n(y) confirm  (n) cancel")
			s += "\n" + centerStyle.Render(dialog) + "\n"
		case ModeResetConfirm:
			dialog := dialogStyle.Render("Reset current phase?\n\n(y) confirm  (n) cancel")
			s += "\n" + centerStyle.Render(dialog) + "\n"
		case ModePreviousConfirm:
			dialog := dialogStyle.Render("Go to previous phase?\n\n(y) confirm  (n) cancel")
			s += "\n" + centerStyle.Render(dialog) + "\n"
		case ModeNormal:
			s += "\n" + footerStyle.Render(configSummary(m)) + "\n"
			s += footerStyle.Render("(?) help  (t) tasks  (q) quit") + "\n"
		default:
			s += "\n" + footerStyle.Render(configSummary(m)) + "\n"
			s += footerStyle.Render("(?) help  (t) tasks  (q) quit") + "\n"
		}
	}

	if m.width > 0 && m.height > 0 {
		if m.alerting && m.blinkState {
			borderStyle := lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(m.alertColor).
				BorderBackground(m.alertColor).
				Padding(1, 4).
				Width(viewWidth + 12).
				AlignHorizontal(lipgloss.Center)
			s = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
				borderStyle.Render(s),
				lipgloss.WithWhitespaceBackground(m.alertColor))
		} else {
			s = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
		}
	}

	return s
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func phaseLabel(m Model) string {
	color := phaseColor(m)
	style := lipgloss.NewStyle().Foreground(color).Bold(true)
	switch m.session.CurrentPhase {
	case session.Work:
		return style.Render(fmt.Sprintf("Work #%d", m.session.CurrentPomodoro))
	case session.ShortBreak:
		return style.Render("Short Break")
	case session.LongBreak:
		return style.Render("Long Break")
	case session.Idle:
		return style.Render("Idle")
	default:
		return style.Render("Idle")
	}
}

func phaseColor(m Model) lipgloss.Color {
	switch m.session.CurrentPhase {
	case session.Work:
		return bordeaux
	case session.LongBreak:
		return deepBlue
	case session.ShortBreak:
		return turquoise
	case session.Idle:
		return gray
	default:
		return gray
	}
}

func phaseTimerStyle(m Model) lipgloss.Style {
	return timerTextStyle.Foreground(phaseColor(m))
}

func wipLabel(wip *task.Task) string {
	return fmt.Sprintf("\U0001f527 %s [%d/%d]", wip.Name, wip.ActualPomos, wip.EstimatedPomos)
}

func configSummary(m Model) string {
	w := int(m.session.WorkDuration.Minutes())
	s := int(m.session.ShortBreakDuration.Minutes())
	l := int(m.session.LongBreakDuration.Minutes())
	r := m.session.Rounds
	return fmt.Sprintf("w:%d  b:%d  lb:%d  rounds:%d", w, s, l, r)
}

func taskLine(t *task.Task) string {
	maxName := viewWidth - taskLinePrefix
	name := t.Name
	if len(name) > maxName {
		name = name[:maxName-1] + "…"
	}
	pomos := fmt.Sprintf("%2d/%-2d", t.ActualPomos, t.EstimatedPomos)
	switch t.Status {
	case task.InProgress:
		return fmt.Sprintf("│ %s │ \U0001f527 %s", pomos, name)
	case task.Done:
		return fmt.Sprintf("│ %s │ \u2705 %s", pomos, name)
	default:
		return fmt.Sprintf("│ %s │ \u2b1c %s", pomos, name)
	}
}

func phaseCompleteMessage(phase session.Phase) string {
	switch phase {
	case session.Work:
		return "Work session complete! Time for a break."
	case session.ShortBreak:
		return "Short break is over! Back to work."
	case session.LongBreak:
		return "Long break is over! Ready for a new cycle."
	case session.Idle:
		return "Timer is idle."
	default:
		return "Time's up!"
	}
}
