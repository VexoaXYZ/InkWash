package components

import (
	"strings"

	"github.com/VexoaXYZ/inkwash/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectorItem represents an item in the selector
type SelectorItem struct {
	Label       string
	Description string
	Value       interface{}
}

// Selector represents an interactive list selector
type Selector struct {
	Title     string
	Items     []SelectorItem
	Selected  int
	Confirmed bool
	Focused   bool
	MaxHeight int // Maximum visible items (0 = show all)
	offset    int // Scroll offset for large lists
}

// NewSelector creates a new selector
func NewSelector(title string, items []SelectorItem) *Selector {
	return &Selector{
		Title:     title,
		Items:     items,
		Selected:  0,
		Confirmed: false,
		Focused:   false,
		MaxHeight: 0,
		offset:    0,
	}
}

// Focus sets the selector as focused
func (s *Selector) Focus() {
	s.Focused = true
}

// Blur removes focus from the selector
func (s *Selector) Blur() {
	s.Focused = false
}

// Update handles keyboard navigation
func (s *Selector) Update(msg tea.Msg) tea.Cmd {
	if !s.Focused || s.Confirmed {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.Selected > 0 {
				s.Selected--
				s.adjustOffset()
			}

		case "down", "j":
			if s.Selected < len(s.Items)-1 {
				s.Selected++
				s.adjustOffset()
			}

		case "home":
			s.Selected = 0
			s.offset = 0

		case "end":
			s.Selected = len(s.Items) - 1
			s.adjustOffset()

		case "enter":
			s.Confirmed = true
		}
	}

	return nil
}

// adjustOffset adjusts the scroll offset to keep selected item visible
func (s *Selector) adjustOffset() {
	if s.MaxHeight == 0 {
		return
	}

	// Scroll down
	if s.Selected >= s.offset+s.MaxHeight {
		s.offset = s.Selected - s.MaxHeight + 1
	}

	// Scroll up
	if s.Selected < s.offset {
		s.offset = s.Selected
	}
}

// View renders the selector
func (s *Selector) View() string {
	var b strings.Builder

	// Render title
	titleStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Bold(true)

	b.WriteString(titleStyle.Render(s.Title))
	b.WriteString("\n\n")

	// Determine visible range
	startIdx := s.offset
	endIdx := len(s.Items)

	if s.MaxHeight > 0 && endIdx-startIdx > s.MaxHeight {
		endIdx = startIdx + s.MaxHeight
	}

	// Render items
	for i := startIdx; i < endIdx; i++ {
		item := s.Items[i]
		isSelected := i == s.Selected

		var itemStr string
		if isSelected {
			// Selected item (highlighted)
			itemStr = ui.StyleSelected.Render(ui.SymbolPointer + " " + item.Label)
		} else {
			// Unselected item
			itemStr = ui.StyleUnselected.Render("  " + item.Label)
		}

		b.WriteString(itemStr)

		// Show description for selected item
		if isSelected && item.Description != "" {
			b.WriteString("\n  ")
			descStyle := lipgloss.NewStyle().
				Foreground(ui.ColorMediumGray).
				Italic(true)
			b.WriteString(descStyle.Render(item.Description))
		}

		b.WriteString("\n")
	}

	// Show scroll indicators if needed
	if s.MaxHeight > 0 && len(s.Items) > s.MaxHeight {
		b.WriteString("\n")
		scrollInfo := lipgloss.NewStyle().
			Foreground(ui.ColorMediumGray)

		if s.offset > 0 {
			b.WriteString(scrollInfo.Render(ui.SymbolArrowUp + " More above"))
		}
		if endIdx < len(s.Items) {
			if s.offset > 0 {
				b.WriteString("  ")
			}
			b.WriteString(scrollInfo.Render(ui.SymbolArrowDown + " More below"))
		}
	}

	// Show navigation help if focused
	if s.Focused && !s.Confirmed {
		b.WriteString("\n\n")
		helpStyle := lipgloss.NewStyle().
			Foreground(ui.ColorMediumGray).
			Italic(true)
		b.WriteString(helpStyle.Render("↑/↓ or j/k: Navigate  •  Enter: Select"))
	}

	return b.String()
}

// SelectedValue returns the value of the currently selected item
func (s *Selector) SelectedValue() interface{} {
	if s.Selected >= 0 && s.Selected < len(s.Items) {
		return s.Items[s.Selected].Value
	}
	return nil
}

// SelectedItem returns the currently selected item
func (s *Selector) SelectedItem() *SelectorItem {
	if s.Selected >= 0 && s.Selected < len(s.Items) {
		return &s.Items[s.Selected]
	}
	return nil
}

// Reset resets the selector to initial state
func (s *Selector) Reset() {
	s.Selected = 0
	s.Confirmed = false
	s.offset = 0
}
