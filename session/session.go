package session

import "time"

type Session struct {
	WorkDuration       time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
	Rounds             int
	CurrentPhase       Phase
	CurrentPomodoro    int
}

// Config holds the durations and rounds for a Pomodoro session.
type Config struct {
	WorkDuration       time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
	Rounds             int
}

// DefaultConfig returns the standard Pomodoro configuration.
func DefaultConfig() Config {
	return Config{
		WorkDuration:       25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		Rounds:             4,
		
	}
}

// NewSession returns a session with the given configuration.
func NewSession(cfg Config) Session {
	return Session{
		WorkDuration:       cfg.WorkDuration,
		ShortBreakDuration: cfg.ShortBreakDuration,
		LongBreakDuration:  cfg.LongBreakDuration,
		Rounds:             cfg.Rounds,
		CurrentPomodoro:    1,
	}
}

// NextPhase advances the session to the next phase in the Pomodoro cycle.
func (s *Session) NextPhase() {
	switch s.CurrentPhase {
	case Idle:
		s.CurrentPhase = Work
	case Work:
		if s.CurrentPomodoro == s.Rounds {
			s.CurrentPhase = LongBreak
		} else {
			s.CurrentPhase = ShortBreak
		}
	case ShortBreak:
		s.CurrentPhase = Work
		s.CurrentPomodoro++
	case LongBreak:
		s.CurrentPhase = Work
		s.CurrentPomodoro = 1
	}
}

// PreviousPhase moves the session back to the previous phase in the Pomodoro cycle.
func (s *Session) PreviousPhase() {
	switch s.CurrentPhase {
	case Work:
		if s.CurrentPomodoro > 1 {
			s.CurrentPhase = ShortBreak
			s.CurrentPomodoro--
		}
	case ShortBreak:
		s.CurrentPhase = Work
	case LongBreak:
		s.CurrentPhase = Work
	}
}

// PhaseDuration returns the duration of the current phase.
func (s *Session) PhaseDuration() time.Duration {
	switch s.CurrentPhase {
	case ShortBreak:
		return s.ShortBreakDuration
	case LongBreak:
		return s.LongBreakDuration
	default:
		return s.WorkDuration
	}
}
