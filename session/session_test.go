package session

import (
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	tests := []struct {
		name         string
		wantWork     time.Duration
		wantShort    time.Duration
		wantLong     time.Duration
		wantRounds   int
		wantPhase    Phase
		wantPomodoro int
	}{
		{
			name:         "returns default configuration",
			wantWork:     25 * time.Minute,
			wantShort:    5 * time.Minute,
			wantLong:     15 * time.Minute,
			wantRounds:   4,
			wantPhase:    Idle,
			wantPomodoro: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession()

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
			s := NewSession()
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
