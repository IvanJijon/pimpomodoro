# Testing Strategy

## Principles

- Test behavior, not implementation details.
- Focus testing effort where the risk is highest: domain logic and state transitions.
- Side effects are injected via callbacks and replaced with no-ops in tests.

## What We Test

| Layer | Coverage | Rationale |
|-------|----------|-----------|
| `session/` | Heavy | Pure domain logic. Phase transitions, pomodoro counting, edge cases. |
| `tui/update.go` | Selective | State transitions driven by key presses and tick messages. Testable as pure functions: `(Model, Msg) → (Model, Cmd)`. |

## What We Don't Test

| Layer | Rationale |
|-------|-----------|
| `tui/view.go` | String output is brittle. Visual correctness is validated manually. |
| `sound/`, `notify/` | OS-level side effects. Injected via `Callbacks`, stubbed in tests. |

## Future Considerations

- **`storage/` package** (planned — task persistence via JSON): Test read/write logic using a real temp directory (`t.TempDir()`), not mocks.
- **E2E with [`teatest`](https://github.com/charmbracelet/x/tree/main/exp/teatest)**: Charm's testing utility for Bubble Tea apps. Simulates key presses and asserts on rendered output. Useful for full user flows, but brittle when styling changes. Consider once the app grows. Low ROI for now since `update` unit tests already cover the logic.

## Conventions

- Table-driven tests for all unit tests.
- No magic numbers — use values derived from the system under test.
- Tests must be independent. No shared mutable state.
