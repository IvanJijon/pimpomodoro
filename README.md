# Pimpomodoro

A lightweight, terminal-based Pomodoro timer.

## Why the name?

Since this is a lightweight Pomodoro tool, it is named after the smallest tomato — *Solanum pimpinellifolium*, native to Ecuador and other countries in South America. `pimp` + `pomodoro` = `pimpomodoro`.

Command: `pimpom`

## Usage

```bash
pimpom [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--work` | 25 | Work duration in minutes |
| `--break` | 5 | Short break duration in minutes |
| `--long-break` | 15 | Long break duration in minutes |
| `--rounds` | 4 | Number of pomodoros before long break |
| `--no-sound` | false | Disable alarm sound |
| `--no-notify` | false | Disable desktop notifications |
| `--no-confirm` | false | Disable confirmation dialogs |
| `--visual-alert` | false | Enable visual alert (blinking) when timer expires |

### Examples

```bash
# Default Pomodoro (25/5/15, 4 rounds)
pimpom

# Custom durations
pimpom --work 50 --break 10 --long-break 30 --rounds 3

# Silent mode
pimpom --no-sound --no-notify
```

### Keybindings

| Key | Action |
|-----|--------|
| `s` | Start/pause |
| `r` | Reset current phase |
| `n` | Skip to next phase |
| `b` | Go to previous phase |
| `t` | Toggle task list |
| `?` | Toggle help |
| `q` | Quit |
| `esc` | Go back |

### Task List Keybindings

| Key | Action |
|-----|--------|
| `↑`/`↓` or `j`/`k` | Navigate tasks |
| `enter` | Select task as WIP |
| `a` | Add new task |
| `d` | Mark task as done |
| `x` | Remove task |


## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for technical details on how Pimpomodoro is built.
