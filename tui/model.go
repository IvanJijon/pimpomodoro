package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/notify"
	"github.com/IvanJijon/pimpomodoro/session"
	"github.com/IvanJijon/pimpomodoro/sound"
)

// ViewMode represents the current UI mode.
type ViewMode int

const (
	ModeNormal ViewMode = iota
	ModeHelp
	ModeSkipConfirm
	ModeQuitConfirm
	ModeResetConfirm
	ModePreviousConfirm
)

// Callbacks holds external side-effect functions injected into the model.
type Callbacks struct {
	PlayAlarm  func()
	SendNotify func(string, string)
}

// DefaultCallbacks returns callbacks wired to real implementations.
func DefaultCallbacks() Callbacks {
	return Callbacks{
		PlayAlarm:  sound.PlayAlarm,
		SendNotify: notify.Send,
	}
}

// Model holds the application state.
type Model struct {
	session       session.Session
	spinner       spinner.Model
	remainingTime time.Duration
	running       bool
	viewMode      ViewMode
	tickID        int
	width         int
	height        int
	confirmEnabled bool

	// not model but injected callbacks for side effects
	callbacks Callbacks
}

// NewModel returns a Model with default session and UI.
func NewModel(cfg session.Config, cb Callbacks) Model {
	return Model{
		session:        session.NewSession(cfg),
		spinner:        newSpinner(),
		callbacks:      cb,
		confirmEnabled: true,
	}
}

// Init returns the initial command. No command is needed at startup.
func (m Model) Init() tea.Cmd {
	return nil
}
