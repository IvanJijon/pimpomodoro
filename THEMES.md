# Themes

Pimpomodoro supports custom color themes via a YAML file at `~/.pimpom/theme.yaml`.

## YAML Format

```yaml
work: "#722F37"
short-break: "#40E0D0"
long-break: "#1E3A5F"
paused: "#FFD700"
subtle: "#666666"
```

| Key | Description | Default |
|-----|-------------|---------|
| `work` | Work phase color | `#722F37` (bordeaux) |
| `short-break` | Short break phase color | `#40E0D0` (turquoise) |
| `long-break` | Long break phase color | `#1E3A5F` (deep blue) |
| `paused` | Paused indicator color | `#FFD700` (gold) |
| `subtle` | Footer and secondary text | `#666666` (gray) |

## Partial Themes

You only need to specify the colors you want to change. Missing fields fall back to defaults.

```yaml
# Only override work and paused colors
work: "#E74C3C"
paused: "#F39C12"
```

## Example: Solarized Dark

```yaml
work: "#DC322F"
short-break: "#2AA198"
long-break: "#268BD2"
paused: "#B58900"
subtle: "#586E75"
```

## Note on Background Color

Pimpomodoro inherits the terminal's background color. The theme controls text and accent colors only — the background is determined by your terminal emulator's settings.
