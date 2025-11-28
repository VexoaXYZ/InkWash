package components

import (
	"time"

	"github.com/VexoaXYZ/inkwash/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Spinner represents a loading spinner component
type Spinner struct {
	Frames []string
	FPS    time.Duration
	index  int
	Tier   ui.AnimationTier
}

// NewSpinner creates a new spinner based on animation tier
func NewSpinner(tier ui.AnimationTier) *Spinner {
	s := &Spinner{
		Tier: tier,
	}

	switch tier {
	case ui.TierMinimal:
		s.Frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		s.FPS = 80 * time.Millisecond

	case ui.TierBalanced:
		s.Frames = []string{"▰▱▱▱▱▱", "▰▰▱▱▱▱", "▰▰▰▱▱▱", "▰▰▰▰▱▱", "▰▰▰▰▰▱", "▰▰▰▰▰▰"}
		s.FPS = 100 * time.Millisecond

	case ui.TierFull:
		s.Frames = []string{
			"◐", "◓", "◑", "◒",
		}
		s.FPS = 80 * time.Millisecond

	default:
		return NewSpinner(ui.TierBalanced)
	}

	return s
}

// Tick advances the spinner to the next frame
func (s *Spinner) Tick() {
	s.index = (s.index + 1) % len(s.Frames)
}

// View renders the current frame of the spinner
func (s *Spinner) View() string {
	return lipgloss.NewStyle().
		Foreground(ui.ColorPrimary).
		Render(s.Frames[s.index])
}

// TickCmd returns a tea.Cmd that sends a tick message
func (s *Spinner) TickCmd() tea.Cmd {
	return tea.Tick(s.FPS, func(t time.Time) tea.Msg {
		return SpinnerTickMsg(t)
	})
}

// SpinnerTickMsg is sent when the spinner should advance
type SpinnerTickMsg time.Time
