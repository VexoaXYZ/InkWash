package animation

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TransitionMsg is a message sent during transitions
type TransitionMsg struct {
	ID       string
	Progress float64
	Complete bool
}

// Transition represents an animated transition
type Transition struct {
	ID         string
	Duration   time.Duration
	Easing     EasingFunc
	OnUpdate   func(progress float64)
	OnComplete func()
}

// Start begins the transition animation
func (t *Transition) Start() tea.Cmd {
	if t.Easing == nil {
		t.Easing = EaseInOutCubic
	}

	return func() tea.Msg {
		start := time.Now()
		ticker := time.NewTicker(16 * time.Millisecond) // 60fps
		defer ticker.Stop()

		for {
			<-ticker.C

			elapsed := time.Since(start)
			progress := float64(elapsed) / float64(t.Duration)

			if progress >= 1.0 {
				progress = 1.0

				if t.OnComplete != nil {
					t.OnComplete()
				}

				return TransitionMsg{
					ID:       t.ID,
					Progress: 1.0,
					Complete: true,
				}
			}

			easedProgress := t.Easing(progress)

			if t.OnUpdate != nil {
				t.OnUpdate(easedProgress)
			}

			return TransitionMsg{
				ID:       t.ID,
				Progress: easedProgress,
				Complete: false,
			}
		}
	}
}

// TickMsg is a generic tick message for animations
type TickMsg time.Time

// Tick returns a command that sends a TickMsg after the specified duration
func Tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// EveryTick returns a command that continuously ticks
func EveryTick(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return TickMsg(time.Now())
	}
}

// FadeInText creates a fade-in animation for text
func FadeInText(text string, duration time.Duration) *Transition {
	return &Transition{
		ID:       "fade-in-text",
		Duration: duration,
		Easing:   EaseInOutSine,
	}
}

// SlideIn creates a slide-in animation
func SlideIn(duration time.Duration) *Transition {
	return &Transition{
		ID:       "slide-in",
		Duration: duration,
		Easing:   EaseOutExpo,
	}
}

// PulseAnimation creates a pulsing animation effect
func PulseAnimation(duration time.Duration) *Transition {
	return &Transition{
		ID:       "pulse",
		Duration: duration,
		Easing:   EaseInOutSine,
	}
}
