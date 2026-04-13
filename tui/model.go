package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/IvanJijon/pimpomodoro/notify"
	"github.com/IvanJijon/pimpomodoro/session"
	"github.com/IvanJijon/pimpomodoro/sound"
	"github.com/IvanJijon/pimpomodoro/task"
	"github.com/IvanJijon/pimpomodoro/theme"
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
	ModeTaskList
	ModeTaskAdd
	ModeTaskEdit
	ModeSwitchTaskConfirm
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

// AppConfig holds all configuration for the application.
type AppConfig struct {
	Session        session.Config
	Callbacks      Callbacks
	ConfirmEnabled bool
	VisualAlert    bool
	Theme          theme.Theme
}

// Model holds the application state.
type Model struct {
	// Timer
	session       session.Session
	remainingTime time.Duration
	running       bool
	tickID        int

	// Visual alert
	visualAlert bool
	alerting    bool
	blinkState  bool
	alertColor  lipgloss.Color

	// UI
	width  int
	height int
	theme  theme.Theme

	spinner  spinner.Model
	viewMode ViewMode

	taskList          *task.TaskList
	taskCursor        int
	taskNameInput     textinput.Model
	taskEstimateInput textinput.Model

	// Config & callbacks
	confirmEnabled bool
	callbacks      Callbacks
}

// NewModel returns a Model with the given application configuration.
func NewModel(cfg AppConfig) Model {
	estimation := func() textinput.Model {
		ti := textinput.New()
		ti.Placeholder = "1"
		ti.CharLimit = 2
		ti.Width = 2
		ti.Prompt = "│ "
		ti.TextStyle = lipgloss.NewStyle().Underline(true)
		return ti
	}

	taskName := func() textinput.Model {
		ti := textinput.New()
		ti.Placeholder = "task name"
		ti.CharLimit = 50
		ti.Width = 30
		ti.Prompt = "│ "
		ti.TextStyle = lipgloss.NewStyle().Underline(true)
		return ti
	}

	var zeroTheme theme.Theme
	if cfg.Theme == zeroTheme {
		cfg.Theme = theme.DefaultTheme()
	}

	return Model{
		session:           session.NewSession(cfg.Session),
		spinner:           newSpinner(),
		callbacks:         cfg.Callbacks,
		confirmEnabled:    cfg.ConfirmEnabled,
		visualAlert:       cfg.VisualAlert,
		taskList:          task.NewTaskList(),
		taskCursor:        0,
		taskNameInput:     taskName(),
		taskEstimateInput: estimation(),
		theme:             cfg.Theme,
	}
}

// Init returns the initial command. No command is needed at startup.
func (m Model) Init() tea.Cmd {
	return nil
}
