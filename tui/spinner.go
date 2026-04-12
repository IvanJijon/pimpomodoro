package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func newSpinner() spinner.Model {
	s := spinner.New(WithDot(), WithStyle())

	return s
}

// WithDot returns a spinner option that sets the spinner to the "Dot" style.
func WithDot() spinner.Option {
	return func(m *spinner.Model) {
		m.Spinner = spinner.Dot
	}
}

// WithStyle returns a spinner option that applies a custom style to the spinner.
func WithStyle() spinner.Option {
	return func(m *spinner.Model) {
		m.Style = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#C15C5C"))
	}
}
