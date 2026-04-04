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
	centerStyle = lipgloss.NewStyle().
			Width(viewWidth).
			AlignHorizontal(lipgloss.Center)

	timerTextStyle = lipgloss.NewStyle().Bold(true)

	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Width(viewWidth - 4).
			AlignHorizontal(lipgloss.Center)
)

// View renders the current state to the terminal.
func (m Model) View() string {
	s := m.header.View()

	if m.showHelp {
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
	if m.running {
		s += "\n" + centerStyle.Render(strings.TrimRight(m.spinner.View(), " ")+timerTextStyle.Render(formatDuration(m.remainingTime))) + "\n"
	} else if m.session.CurrentPhase != session.Idle {
		s += "\n" + centerStyle.Render(timerTextStyle.Render("⏸ "+formatDuration(m.remainingTime))) + "\n"
	} else {
		s += "\n" + centerStyle.Render(timerTextStyle.Render(formatDuration(m.remainingTime))) + "\n"
	}

	if m.showSkipConfirm {
		dialog := dialogStyle.Render("Skip to next phase?\n\n(y) confirm  (x) cancel")
		s += "\n" + centerStyle.Render(dialog) + "\n"
	} else {
		s += "\n" + centerStyle.Render("? help  q quit") + "\n"
	}

	return s
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func phaseLabel(m Model) string {
	switch m.session.CurrentPhase {
	case session.Work:
		return fmt.Sprintf("Work #%d", m.session.CurrentPomodoro)
	case session.ShortBreak:
		return "Short Break"
	case session.LongBreak:
		return "Long Break"
	default:
		return "Idle"
	}
}
