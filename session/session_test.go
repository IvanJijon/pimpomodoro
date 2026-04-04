package session

import (
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	tests := []struct {
		name         string
		config       Config
		wantWork     time.Duration
		wantShort    time.Duration
		wantLong     time.Duration
		wantRounds   int
		wantPhase    Phase
		wantPomodoro int
	}{
		{
			name:         "returns default configuration",
			config:       DefaultConfig(),
			wantWork:     25 * time.Minute,
			wantShort:    5 * time.Minute,
			wantLong:     15 * time.Minute,
			wantRounds:   4,
			wantPhase:    Idle,
			wantPomodoro: 1,
		},
		{
			name: "returns custom configuration",
			config: Config{
				WorkDuration:       50 * time.Minute,
				ShortBreakDuration: 10 * time.Minute,
				LongBreakDuration:  30 * time.Minute,
				Rounds:             3,
			},
			wantWork:     50 * time.Minute,
			wantShort:    10 * time.Minute,
			wantLong:     30 * time.Minute,
			wantRounds:   3,
			wantPhase:    Idle,
			wantPomodoro: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession(tt.config)

			if s.WorkDuration != tt.wantWork {
				t.Errorf("WorkDuration = %v, want %v", s.WorkDuration, tt.wantWork)
			}
			if s.ShortBreakDuration != tt.wantShort {
				t.Errorf("ShortBreakDuration = %v, want %v", s.ShortBreakDuration, tt.wantShort)
			}
			if s.LongBreakDuration != tt.wantLong {
				t.Errorf("LongBreakDuration = %v, want %v", s.LongBreakDuration, tt.wantLong)
			}
			if s.Rounds != tt.wantRounds {
				t.Errorf("Rounds = %v, want %v", s.Rounds, tt.wantRounds)
			}
			if s.CurrentPhase != tt.wantPhase {
				t.Errorf("CurrentPhase = %v, want %v", s.CurrentPhase, tt.wantPhase)
			}
			if s.CurrentPomodoro != tt.wantPomodoro {
				t.Errorf("CurrentPomodoro = %v, want %v", s.CurrentPomodoro, tt.wantPomodoro)
			}
		})
	}
}

func TestNextPhase(t *testing.T) {
	tests := []struct {
		name         string
		phase        Phase
		pomodoro     int
		wantPhase    Phase
		wantPomodoro int
	}{
		{
			name:         "transitions from Idle to Work",
			phase:        Idle,
			pomodoro:     1,
			wantPhase:    Work,
			wantPomodoro: 1,
		},
		{
			name:         "transitions from Work to ShortBreak",
			phase:        Work,
			pomodoro:     1,
			wantPhase:    ShortBreak,
			wantPomodoro: 1,
		},
		{
			name:         "transitions from ShortBreak to Work",
			phase:        ShortBreak,
			pomodoro:     1,
			wantPhase:    Work,
			wantPomodoro: 2,
		},
		{
			name:         "transitions from Work to LongBreak on last round",
			phase:        Work,
			pomodoro:     4,
			wantPhase:    LongBreak,
			wantPomodoro: 4,
		},
		{
			name:         "transitions from LongBreak to Work and resets pomodoro",
			phase:        LongBreak,
			pomodoro:     4,
			wantPhase:    Work,
			wantPomodoro: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession(DefaultConfig())
			s.CurrentPhase = tt.phase
			s.CurrentPomodoro = tt.pomodoro

			s.NextPhase()

			if s.CurrentPhase != tt.wantPhase {
				t.Errorf("CurrentPhase = %v, want %v", s.CurrentPhase, tt.wantPhase)
			}
			if s.CurrentPomodoro != tt.wantPomodoro {
				t.Errorf("CurrentPomodoro = %v, want %v", s.CurrentPomodoro, tt.wantPomodoro)
			}
		})
	}
}

func TestPhaseDuration(t *testing.T) {
	tests := []struct {
		name         string
		phase        Phase
		wantDuration time.Duration
	}{
		{
			name:         "returns work duration when phase is Work",
			phase:        Work,
			wantDuration: 25 * time.Minute,
		},
		{
			name:         "returns short break duration when phase is ShortBreak",
			phase:        ShortBreak,
			wantDuration: 5 * time.Minute,
		},
		{
			name:         "returns long break duration when phase is LongBreak",
			phase:        LongBreak,
			wantDuration: 15 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession(DefaultConfig())
			s.CurrentPhase = tt.phase

			got := s.PhaseDuration()

			if got != tt.wantDuration {
				t.Errorf("PhaseDuration() = %v, want %v", got, tt.wantDuration)
			}
		})
	}
}

func TestPreviousPhase(t *testing.T) {
	tests := []struct {
		name         string
		phase        Phase
		pomodoro     int
		wantPhase    Phase
		wantPomodoro int
	}{
		{
			name:         "Work #1 is a no-op",
			phase:        Work,
			pomodoro:     1,
			wantPhase:    Work,
			wantPomodoro: 1,
		},
		{
			name:         "Idle is a no-op",
			phase:        Idle,
			pomodoro:     1,
			wantPhase:    Idle,
			wantPomodoro: 1,
		},
		{
			name:         "ShortBreak goes back to Work same pomodoro",
			phase:        ShortBreak,
			pomodoro:     2,
			wantPhase:    Work,
			wantPomodoro: 2,
		},
		{
			name:         "Work #2 goes back to ShortBreak and decrements pomodoro",
			phase:        Work,
			pomodoro:     2,
			wantPhase:    ShortBreak,
			wantPomodoro: 1,
		},
		{
			name:         "LongBreak goes back to Work last pomodoro",
			phase:        LongBreak,
			pomodoro:     4,
			wantPhase:    Work,
			wantPomodoro: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession(DefaultConfig())
			s.CurrentPhase = tt.phase
			s.CurrentPomodoro = tt.pomodoro

			s.PreviousPhase()

			if s.CurrentPhase != tt.wantPhase {
				t.Errorf("CurrentPhase = %v, want %v", s.CurrentPhase, tt.wantPhase)
			}
			if s.CurrentPomodoro != tt.wantPomodoro {
				t.Errorf("CurrentPomodoro = %v, want %v", s.CurrentPomodoro, tt.wantPomodoro)
			}
		})
	}
}
