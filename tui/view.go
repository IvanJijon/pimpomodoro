package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/IvanJijon/pimpomodoro/session"
)

const viewWidth = 40

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
		s += "\n  Keybindings:\n\n"
		s += "  s  Start/resume\n"
		s += "  p  Pause\n"
		s += "  r  Reset current phase\n"
		s += "  n  Skip to next phase\n"
		s += "  b  Go to previous phase\n"
		s += "  ?  Toggle help\n"
		s += "  q  Quit\n"
		return s
	}

	s += "\n" + centerStyle.Render(phaseLabel(m)) + "\n"
	switch {
	case m.running:
		s += "\n" + centerStyle.Render(strings.TrimRight(m.spinner.View(), " ")+phaseTimerStyle(m).Render(formatDuration(m.remainingTime))) + "\n"
	case m.session.CurrentPhase != session.Idle:
		s += "\n" + centerStyle.Render(pausedStyle.Render("\u23f8 "+formatDuration(m.remainingTime))) + "\n"
	default:
		s += "\n" + centerStyle.Render(timerTextStyle.Render(formatDuration(m.remainingTime))) + "\n"
	}

	switch m.viewMode {
	case ModeSkipConfirm:
		dialog := dialogStyle.Render("Skip to next phase?\n\n(y) confirm  (x) cancel")
		s += "\n" + centerStyle.Render(dialog) + "\n"
	case ModeQuitConfirm:
		dialog := dialogStyle.Render("Quit?\n\n(y) confirm  (x) cancel")
		s += "\n" + centerStyle.Render(dialog) + "\n"
	case ModeResetConfirm:
		dialog := dialogStyle.Render("Reset current phase?\n\n(y) confirm  (x) cancel")
		s += "\n" + centerStyle.Render(dialog) + "\n"
	case ModeNormal:
		s += "\n" + footerStyle.Render("(?) help  (q) quit") + "\n"
	case ModeHelp:
		// ModeHelp is handled at the top of the View function, so nothing needed here.
	default:
		s += "\n" + footerStyle.Render("(?) help  (q) quit") + "\n"
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
