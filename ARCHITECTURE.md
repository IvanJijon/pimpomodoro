# Architecture

## Overview

Pimpomodoro follows the [Elm Architecture](https://guide.elm-lang.org/architecture/) (TEA) via [Bubble Tea](https://github.com/charmbracelet/bubbletea). The application is split into domain logic and presentation, with platform-specific packages for system integration.

## Package Structure

```
pimpomodoro/
├── main.go        Entry point, CLI flags, dependency wiring
├── session/       Domain logic — session config, phase state machine
├── tui/           Terminal UI — Bubble Tea model, update, view
├── sound/         Platform-specific alarm sounds
├── notify/        Platform-specific desktop notifications
└── Makefile       Build, test, release commands
```

## Session

The `session` package is the core domain. It has no UI dependencies.

- `Config` — holds durations (work, short break, long break) and number of rounds
- `Session` — tracks current phase, pomodoro count, and durations
- `Phase` — state machine with four states: `Idle`, `Work`, `ShortBreak`, `LongBreak`

### Phase State Machine

```
Idle → Work → ShortBreak → Work → ShortBreak → ... → Work → LongBreak
                                                                 ↓
                                                          Work (cycle restarts)
```

- `NextPhase()` advances to the next phase in the cycle
- `PreviousPhase()` moves back (no-op at `Work #1` and `Idle`)
- `PhaseDuration()` returns the duration for the current phase

## TUI

The `tui` package implements the Bubble Tea application using the Elm Architecture:

- **Model** — application state (session, timer, running flag, view mode)
- **Update** — handles messages (key presses, ticks) and returns the updated model
- **View** — renders the current state to the terminal as a string

### View Modes

The UI uses a `ViewMode` to determine which screen is active and which keys are handled:

| Mode | Description |
|------|-------------|
| `ModeNormal` | Main timer screen |
| `ModeHelp` | Keybindings help screen |
| `ModeSkipConfirm` | Skip phase confirmation dialog |
| `ModeResetConfirm` | Reset phase confirmation dialog |
| `ModePreviousConfirm` | Previous phase confirmation dialog |
| `ModeQuitConfirm` | Quit confirmation dialog |

Each mode has its own update handler, keeping the key handling clean and separated.

### Timer

The countdown is driven by `tea.Tick`, which sends a `TickMsg` every second. Each tick:

1. Decrements remaining time by one second
2. Schedules the next tick
3. When time reaches zero: transitions to the next phase, plays alarm, sends notification

A **tick ID** mechanism prevents parallel tick loops. Each new tick loop gets an incremented ID. Stale ticks from old loops are ignored.

### Styling

UI styling uses [Lipgloss](https://github.com/charmbracelet/lipgloss) with phase-specific colors:

| Phase | Color |
|-------|-------|
| Work | Bordeaux |
| Short Break | Turquoise |
| Long Break | Deep Blue |
| Paused | Yellow |

## Sound & Notifications

Both packages use Go build tags (`//go:build darwin`, `linux`, `windows`) to compile platform-specific code.

### Sound

| Platform | Method |
|----------|--------|
| macOS | `afplay` with system sounds |
| Linux | `paplay` → `aplay` → terminal bell (fallback chain) |
| Windows | PowerShell `SystemSounds` |

### Notifications

| Platform | Method |
|----------|--------|
| macOS | `osascript` display notification |
| Linux | `notify-send` (if available) |
| Windows | PowerShell toast notification |

Both are injected into the model via a `Callbacks` struct, allowing them to be replaced with no-ops in tests.

## CLI

Flags are parsed in `main.go` using Go's standard `flag` package:

| Flag | Default | Description |
|------|---------|-------------|
| `--work` | 25 | Work duration in minutes |
| `--break` | 5 | Short break duration in minutes |
| `--long-break` | 15 | Long break duration in minutes |
| `--rounds` | 4 | Pomodoros before long break |
| `--no-sound` | false | Disable alarm sound |
| `--no-notify` | false | Disable desktop notifications |
| `--version` | — | Print version and exit |

Version is injected at build time from git tags via `-ldflags`.

## Build & Release

```bash
make build       # Build for current platform
make build-all   # Build for all platforms
make test        # Run all tests
make cover       # Generate coverage report
make release V=x.y.z  # Tag and build release
```
