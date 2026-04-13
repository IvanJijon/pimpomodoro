package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
	"github.com/IvanJijon/pimpomodoro/theme"
	"github.com/IvanJijon/pimpomodoro/tui"
)

var (
	version = "dev"

	work          = flag.Int("work", 25, "work duration in minutes")
	brk           = flag.Int("break", 5, "short break duration in minutes")
	longBrk       = flag.Int("long-break", 15, "long break duration in minutes")
	rounds        = flag.Int("rounds", 4, "number of pomodoros before long break")
	noSound       = flag.Bool("no-sound", false, "disable alarm sound")
	noNotify      = flag.Bool("no-notify", false, "disable desktop notifications")
	noConfirm     = flag.Bool("no-confirm", false, "disable confirmation dialogs")
	noVisualAlert = flag.Bool("no-visual-alert", false, "disable visual alert (blinking) when timer expires")
	showVer       = flag.Bool("version", false, "print version and exit")
)

func main() {
	flag.Parse()

	if *showVer {
		fmt.Printf("pimpom %s\n", version)
		return
	}

	p := tea.NewProgram(tui.NewModel(parseAppConfig()), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oops! There's been an error: %v", err)
		os.Exit(1)
	}
}

func parseAppConfig() tui.AppConfig {
	cb := tui.DefaultCallbacks()
	if *noSound {
		cb.PlayAlarm = func() {}
	}
	if *noNotify {
		cb.SendNotify = func(_, _ string) {}
	}

	var th theme.Theme
	home, err := os.UserHomeDir()
	if err != nil {
		th = theme.DefaultTheme()
	} else {
		th = theme.LoadFromFile(filepath.Join(home, ".config", "pimpom", "theme.yaml"))
	}

	return tui.AppConfig{
		Theme: th,
		Session: session.Config{
			WorkDuration:       time.Duration(*work) * time.Minute,
			ShortBreakDuration: time.Duration(*brk) * time.Minute,
			LongBreakDuration:  time.Duration(*longBrk) * time.Minute,
			Rounds:             *rounds,
		},
		Callbacks:      cb,
		ConfirmEnabled: !*noConfirm,
		VisualAlert:    !*noVisualAlert,
	}
}
