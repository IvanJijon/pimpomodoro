package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
)

// Model holds the application state.
type Model struct {
	session         session.Session
	header          *Header
	spinner         spinner.Model
	remainingTime   time.Duration
	running         bool
	showHelp        bool
	showSkipConfirm bool
	tickID          int
}

// NewModel returns a Model with default session and UI.
func NewModel() Model {
	return Model{
		session: session.NewSession(),
		header:  NewHeader(),
		spinner: newSpinner(),
	}
}

// Init returns the initial command. No command is needed at startup.
func (m Model) Init() tea.Cmd {
	return nil
}
