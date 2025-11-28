package ui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Color palette - Monochrome Elegance
var (
	// Monochrome Foundation
	ColorPureWhite  = lipgloss.Color("#FFFFFF")
	ColorSoftWhite  = lipgloss.Color("#F5F5F5")
	ColorLightGray  = lipgloss.Color("#E5E5E5")
	ColorMediumGray = lipgloss.Color("#A0A0A0")
	ColorDarkGray   = lipgloss.Color("#404040")
	ColorDeepBlack  = lipgloss.Color("#0A0A0A")

	// Single Accent (Adaptable)
	ColorPrimary     = lipgloss.Color("#7C3AED")
	ColorPrimaryDim  = lipgloss.Color("#6D28D9")
	ColorPrimaryGlow = lipgloss.Color("#8B5CF6")

	// Semantic (Minimal Usage)
	ColorSuccess = lipgloss.Color("#10B981")
	ColorError   = lipgloss.Color("#EF4444")
	ColorWarning = lipgloss.Color("#F59E0B")
)

// Base styles
var (
	// Text styles
	StyleText = lipgloss.NewStyle().
			Foreground(ColorPureWhite)

	StyleTextMuted = lipgloss.NewStyle().
			Foreground(ColorMediumGray)

	StyleTextDim = lipgloss.NewStyle().
			Foreground(ColorDarkGray)

	// Header styles
	StyleHeader = lipgloss.NewStyle().
			Foreground(ColorPureWhite).
			Bold(true).
			Underline(true)

	StyleSubheader = lipgloss.NewStyle().
			Foreground(ColorSoftWhite).
			Bold(true)

	// Accent styles
	StyleAccent = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	StyleAccentDim = lipgloss.NewStyle().
			Foreground(ColorPrimaryDim)

	StyleAccentGlow = lipgloss.NewStyle().
			Foreground(ColorPrimaryGlow)

	// Status styles
	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning)

	// Border styles
	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorLightGray)

	StyleBorderAccent = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary)

	// Box styles
	StyleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorLightGray).
			Padding(1, 2)

	StyleBoxAccent = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	// Code/Path styles
	StyleCode = lipgloss.NewStyle().
			Foreground(ColorMediumGray).
			Italic(true)

	StylePath = lipgloss.NewStyle().
			Foreground(ColorMediumGray)

	// Help text styles
	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorMediumGray).
			Italic(true)

	// Title bar style
	StyleTitleBar = lipgloss.NewStyle().
			Foreground(ColorPureWhite).
			Background(ColorPrimary).
			Bold(true).
			Padding(0, 1)

	// Status bar style
	StyleStatusBar = lipgloss.NewStyle().
			Foreground(ColorMediumGray).
			Background(ColorDarkGray).
			Padding(0, 1)

	// Selected item style
	StyleSelected = lipgloss.NewStyle().
			Foreground(ColorPureWhite).
			Background(ColorPrimary).
			Padding(0, 1)

	// Unselected item style
	StyleUnselected = lipgloss.NewStyle().
			Foreground(ColorMediumGray).
			Padding(0, 1)

	// Focused input style
	StyleInputFocused = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1)

	// Unfocused input style
	StyleInputUnfocused = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorLightGray).
				Padding(0, 1)
)

// Symbols
const (
	SymbolRunning  = "●"
	SymbolStopped  = "○"
	SymbolPointer  = "▸"
	SymbolCheck    = "✓"
	SymbolCross    = "✗"
	SymbolDot      = "•"
	SymbolLine     = "─"
	SymbolArrowUp  = "↑"
	SymbolArrowDown = "↓"
)

// Spacing helpers
const (
	SpacingMicro  = 1 // Between related items
	SpacingSmall  = 2 // Between sections
	SpacingMedium = 3 // Between major sections
	SpacingLarge  = 5 // Before/after critical actions
)

// Animation timing
const (
	CursorBlinkRate = 500 * time.Millisecond
)

// RenderTitle renders a title with the accent color
func RenderTitle(title string) string {
	return StyleTitleBar.Render(title)
}

// RenderHeader renders a header with proper styling
func RenderHeader(text string) string {
	return StyleHeader.Render(text)
}

// RenderSubheader renders a subheader
func RenderSubheader(text string) string {
	return StyleSubheader.Render(text)
}

// RenderAccent renders text with accent color
func RenderAccent(text string) string {
	return StyleAccent.Render(text)
}

// RenderSuccess renders success text
func RenderSuccess(text string) string {
	return StyleSuccess.Render(SymbolCheck + " " + text)
}

// RenderError renders error text
func RenderError(text string) string {
	return StyleError.Render(SymbolCross + " " + text)
}

// RenderWarning renders warning text
func RenderWarning(text string) string {
	return StyleWarning.Render(text)
}

// RenderMuted renders muted text
func RenderMuted(text string) string {
	return StyleTextMuted.Render(text)
}

// RenderCode renders code/path text
func RenderCode(text string) string {
	return StyleCode.Render(text)
}

// RenderPath renders a file path
func RenderPath(path string) string {
	return StylePath.Render(path)
}

// RenderHelp renders help text
func RenderHelp(text string) string {
	return StyleHelp.Render(text)
}

// RenderBox renders content in a bordered box
func RenderBox(content string) string {
	return StyleBox.Render(content)
}

// RenderBoxAccent renders content in a bordered box with accent color
func RenderBoxAccent(content string) string {
	return StyleBoxAccent.Render(content)
}

// RenderStatusRunning renders a running status indicator
func RenderStatusRunning(text string) string {
	return StyleSuccess.Render(SymbolRunning) + " " + text
}

// RenderStatusStopped renders a stopped status indicator
func RenderStatusStopped(text string) string {
	return StyleTextMuted.Render(SymbolStopped) + " " + text
}

// NewSpacing returns a string of newlines for spacing
func NewSpacing(lines int) string {
	spacing := ""
	for i := 0; i < lines; i++ {
		spacing += "\n"
	}
	return spacing
}
