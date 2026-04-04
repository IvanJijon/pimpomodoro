package session

import "time"

type Phase int

const (
	Idle Phase = iota
	Work
	ShortBreak
	LongBreak
)

type Session struct {
	WorkDuration       time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
	Rounds             int
	CurrentPhase       Phase
	CurrentPomodoro    int
}

// NewSession returns a session with default Pomodoro configuration.
func NewSession() Session {
	return Session{
		WorkDuration:       25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		Rounds:             4,
		CurrentPomodoro:    1,
	}
}

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
