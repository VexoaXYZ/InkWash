package components

import (
	"strings"

	"github.com/VexoaXYZ/inkwash/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

// Sparkline represents an inline graph component
type Sparkline struct {
	Data  []float64
	Width int
	Max   float64
}

// NewSparkline creates a new sparkline
func NewSparkline(width int) *Sparkline {
	return &Sparkline{
		Data:  make([]float64, 0),
		Width: width,
	}
}

// AddDataPoint adds a new data point to the sparkline
func (s *Sparkline) AddDataPoint(value float64) {
	s.Data = append(s.Data, value)

	// Keep only the last 'Width' points
	if len(s.Data) > s.Width {
		s.Data = s.Data[1:]
	}

	// Update max value
	if value > s.Max {
		s.Max = value
	}
}

// Render renders the sparkline
func (s *Sparkline) Render() string {
	return s.RenderWithColor(ui.ColorPrimary)
}

// RenderWithColor renders the sparkline with a specific color
func (s *Sparkline) RenderWithColor(color lipgloss.Color) string {
	// Characters for different heights (8 levels)
	bars := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	var result strings.Builder
	style := lipgloss.NewStyle().Foreground(color)

	for i := 0; i < s.Width; i++ {
		if i >= len(s.Data) {
			result.WriteRune(' ')
			continue
		}

		value := s.Data[i]

		// Handle zero or negative max
		if s.Max <= 0 {
			result.WriteRune(bars[0])
			continue
		}

		// Calculate percentage and map to bar index
		percentage := value / s.Max
		index := int(percentage * float64(len(bars)-1))

		if index < 0 {
			index = 0
		}
		if index >= len(bars) {
			index = len(bars) - 1
		}

		result.WriteRune(bars[index])
	}

	return style.Render(result.String())
}

// Clear clears all data points
func (s *Sparkline) Clear() {
	s.Data = make([]float64, 0)
	s.Max = 0
}

// SetMax manually sets the maximum value for scaling
func (s *Sparkline) SetMax(max float64) {
	s.Max = max
}
