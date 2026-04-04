package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func newSpinner() spinner.Model {
	s := spinner.New()

	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C15C5C"))

	return s
}
