# Architecture

## Overview

Pimpomodoro follows the [Elm Architecture](https://guide.elm-lang.org/architecture/) (TEA) via [Bubble Tea](https://github.com/charmbracelet/bubbletea). The application is split into domain logic and presentation, with platform-specific packages for system integration.

## Package Structure

```
pimpomodoro/
├── main.go              Entry point, CLI flags, parseAppConfig()
├── session/             Domain logic — session config, phase state machine
│   ├── session.go       Session, Config, state machine
│   └── phase.go         Phase type and constants
├── task/                Task domain model
│   ├── task.go          Task entity: status, pomodoros, state transitions
│   └── tasklist.go      TaskList aggregate root: add, remove, select WIP, mark done, sort
├── tui/                 Terminal UI — Bubble Tea model, update, view
│   ├── model.go         Model struct, AppConfig, Callbacks, ViewMode, NewModel, Init
│   ├── update.go        Update dispatcher (thin router to handlers)
│   ├── handle_key.go    Global key handling (q, ctrl-c) + per-mode key handlers
│   ├── handle_tick.go   Window size, spinner tick, timer tick handlers
│   ├── handle_blink.go  Visual alert blink handler
│   ├── tick.go          TickMsg and tickCmd (1s interval)
│   ├── blink.go         BlinkMsg and blinkCmd (500ms interval)
│   ├── view.go          View rendering + styles + display helpers
│   ├── spinner.go       Spinner component config
│   ├── colors.go        Color constants
│   ├── update_test.go   Table-driven tests for update logic
│   └── view_test.go     Tests for formatDuration
├── sound/               Platform-specific alarm sounds
├── notify/              Platform-specific desktop notifications
└── Makefile             Build, test, release commands
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

## Task

The `task` package contains the task domain model, independent of UI.

### Task (Entity)

Fields: Name, EstimatedPomos, ActualPomos, Status (Pending/InProgress/Done)

State transitions:
- `StartWork()` — Pending → InProgress
- `StopWork()` — InProgress → Pending
- `MarkDone()` — InProgress → Done
- `IncreaseActualPomos()` — increments count (only when InProgress)

### TaskList (Aggregate Root)

Ordered collection of `*Task` with invariant: only one task can be InProgress at a time.

- `Add(*Task)`, `Remove(*Task)`, `Len()`, `Tasks()`, `CurrentWIP()`
- `SelectWIP(*Task)` — starts task, stops previous WIP, re-sorts
- `MarkTaskDone(*Task)` — marks done, re-sorts (validates membership)
- Sort order: InProgress → Pending → Done (via `statusPriority` map)

## TUI

The `tui` package implements the Bubble Tea application using the Elm Architecture:

- **Model** — application state, grouped by subject (timer, visual alert, UI, config & callbacks)
- **Update** — thin dispatcher that routes messages to dedicated handler files
- **View** — renders the current state to the terminal as a string

### Update Dispatcher

`update.go` routes each message type to its handler:

| Message | Handler file | Handler method |
|---------|-------------|----------------|
| `tea.WindowSizeMsg` | `handle_tick.go` | `handleWindowSize` |
| `spinner.TickMsg` | `handle_tick.go` | `handleSpinnerTick` |
| `TickMsg` | `handle_tick.go` | `handleTick` |
| `BlinkMsg` | `handle_blink.go` | `handleBlink` |
| `tea.KeyMsg` | `handle_key.go` | `handleKey` → per-mode handlers |

### Global Key Handling

`handleKey` processes these before dispatching to mode-specific handlers:
- `ctrl-c` — force quit from any mode
- `q` — quit with confirmation from any mode except ModeTaskAdd

### View Modes

| Mode | Description | Back key |
|------|-------------|----------|
| `ModeNormal` | Main timer screen | — |
| `ModeHelp` | Keybindings help | `?` or `esc` |
| `ModeTaskList` | Task list with navigation | `t` or `esc` |
| `ModeTaskAdd` | Add task form (text inputs) | `esc` |
| `ModeSkipConfirm` | Skip phase confirmation | `n` or `esc` |
| `ModeResetConfirm` | Reset phase confirmation | `n` or `esc` |
| `ModePreviousConfirm` | Previous phase confirmation | `n` or `esc` |
| `ModeQuitConfirm` | Quit confirmation | `n` or `esc` |

### Timer

The countdown is driven by `tea.Tick`, which sends a `TickMsg` every second. Each tick:

1. Decrements remaining time by one second
2. Schedules the next tick
3. When time reaches zero: transitions to the next phase, plays alarm, sends notification, starts blink loop if visual alert is enabled

A **tick ID** mechanism prevents parallel tick loops. Each new tick loop gets an incremented ID. Stale ticks from old loops are ignored.

### Visual Alert

Accessibility feature for users with hearing difficulties. Enabled via `--visual-alert` flag (off by default).

- On timer expiry: `alerting = true`, `alertColor` captures the completed phase's color
- `BlinkMsg` fires every 500ms via `blinkCmd`, toggling `blinkState`
- View applies phase-colored background to the entire terminal when `alerting && blinkState`
- Any keypress dismisses the alert (cleared at top of `handleKey`)

### Task List UI

Keybindings in ModeTaskList:
- `↑`/`↓` or `j`/`k` — navigate
- `enter` — select task at cursor as WIP
- `d` — mark task at cursor as done
- `x` — remove task at cursor
- `a` — enter ModeTaskAdd
- `esc` or `t` — back to ModeNormal

Task display: `│ 0/3  │ ⬜ Task name` with dynamic truncation based on `viewWidth`.

### Task Add UI

Two `bubbles/textinput` fields:
- Name (50 char limit, `│` prompt, underlined)
- Estimate (2 char limit, `│` prompt, underlined)
- `tab` switches focus between inputs
- `enter` confirms (empty name rejected, invalid estimate defaults to 1)
- `esc` cancels

### Normal Mode WIP Display

When a task is InProgress, normal mode shows it below the timer: `🔧 Task name [actual/estimated]`

### Confirmation Dialogs

Skip, reset, previous, and quit all follow the same pattern:
1. Key press → set `viewMode` to confirm mode, pause timer
2. View shows dialog with `(y) confirm  (n) cancel`
3. `y` → execute action
4. `n` or `esc` → cancel, resume timer

Confirmations can be disabled globally via `--no-confirm` flag.

### Styling

UI styling uses [Lipgloss](https://github.com/charmbracelet/lipgloss) with phase-specific colors:

| Phase | Color |
|-------|-------|
| Work | Bordeaux |
| Short Break | Turquoise |
| Long Break | Deep Blue |
| Paused | Yellow |

View width is 65 chars. All content is centered in the terminal via `lipgloss.Place`. Task list and help text are left-aligned within the centered block.

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

All flags are parsed in `main.go` via `parseAppConfig()` → `tui.AppConfig`.

| Flag | Default | Description |
|------|---------|-------------|
| `--work` | 25 | Work duration in minutes |
| `--break` | 5 | Short break duration in minutes |
| `--long-break` | 15 | Long break duration in minutes |
| `--rounds` | 4 | Pomodoros before long break |
| `--no-sound` | false | Disable alarm sound |
| `--no-notify` | false | Disable desktop notifications |
| `--no-confirm` | false | Disable confirmation dialogs |
| `--visual-alert` | false | Enable visual alert (blinking) |
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
