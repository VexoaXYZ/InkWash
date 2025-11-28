package components

import (
	"strings"
	"time"

	"github.com/VexoaXYZ/inkwash/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TextInput represents an interactive text input field
type TextInput struct {
	Label        string
	Placeholder  string
	Value        string
	MaxLength    int
	Focused      bool
	Error        string
	Validator    func(string) error
	cursor       int
	showCursor   bool
	clearOnFocus bool // Clear value on first keypress after focus
}

// NewTextInput creates a new text input field
func NewTextInput(label, placeholder string, maxLength int) *TextInput {
	return &TextInput{
		Label:       label,
		Placeholder: placeholder,
		MaxLength:   maxLength,
		Focused:     false,
		cursor:      0,
		showCursor:  true,
	}
}

// SetValidator sets the validation function
func (t *TextInput) SetValidator(validator func(string) error) {
	t.Validator = validator
}

// Focus sets the input as focused
func (t *TextInput) Focus() {
	t.Focused = true
	// Move cursor to end of existing text
	t.cursor = len(t.Value)
	// Mark that we should clear on first keypress (for default values)
	if t.Value != "" {
		t.clearOnFocus = true
	}
}

// Blur removes focus from the input
func (t *TextInput) Blur() {
	t.Focused = false
	t.Validate()
}

// Validate runs the validator if set
func (t *TextInput) Validate() {
	if t.Validator != nil {
		if err := t.Validator(t.Value); err != nil {
			t.Error = err.Error()
		} else {
			t.Error = ""
		}
	}
}

// Update handles key input and cursor blinking
func (t *TextInput) Update(msg tea.Msg) tea.Cmd {
	if !t.Focused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyBackspace:
			// Clear default value on first keypress
			if t.clearOnFocus {
				t.Value = ""
				t.cursor = 0
				t.clearOnFocus = false
			} else if len(t.Value) > 0 && t.cursor > 0 {
				t.Value = t.Value[:t.cursor-1] + t.Value[t.cursor:]
				t.cursor--
			}

		case tea.KeyDelete:
			// Clear default value on first keypress
			if t.clearOnFocus {
				t.Value = ""
				t.cursor = 0
				t.clearOnFocus = false
			} else if t.cursor < len(t.Value) {
				t.Value = t.Value[:t.cursor] + t.Value[t.cursor+1:]
			}

		case tea.KeyLeft:
			if t.cursor > 0 {
				t.cursor--
			}

		case tea.KeyRight:
			if t.cursor < len(t.Value) {
				t.cursor++
			}

		case tea.KeyHome:
			t.cursor = 0

		case tea.KeyEnd:
			t.cursor = len(t.Value)

		case tea.KeySpace:
			// Clear default value on first keypress
			if t.clearOnFocus {
				t.Value = ""
				t.cursor = 0
				t.clearOnFocus = false
			}
			if t.MaxLength == 0 || len(t.Value) < t.MaxLength {
				t.Value = t.Value[:t.cursor] + " " + t.Value[t.cursor:]
				t.cursor++
			}

		case tea.KeyRunes:
			// Clear default value on first keypress
			if t.clearOnFocus {
				t.Value = ""
				t.cursor = 0
				t.clearOnFocus = false
			}
			if t.MaxLength == 0 || len(t.Value) < t.MaxLength {
				t.Value = t.Value[:t.cursor] + string(msg.Runes) + t.Value[t.cursor:]
				t.cursor += len(msg.Runes)
			}
		}

		// Clear error when user types
		t.Error = ""

	case CursorBlinkMsg:
		if t.Focused {
			t.showCursor = !t.showCursor
			return t.BlinkCmd()
		}
	}

	return nil
}

// View renders the text input
func (t *TextInput) View() string {
	var b strings.Builder

	// Render label
	labelStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Bold(true)

	b.WriteString(labelStyle.Render(t.Label))
	b.WriteString("\n")

	// Prepare input text
	displayText := t.Value
	if displayText == "" && !t.Focused {
		displayText = t.Placeholder
	}

	// Add cursor if focused
	if t.Focused && t.showCursor {
		if t.cursor <= len(displayText) {
			displayText = displayText[:t.cursor] + "█" + displayText[t.cursor:]
		}
	}

	// Render input box
	var inputStyle lipgloss.Style
	if t.Focused {
		inputStyle = ui.StyleInputFocused
	} else {
		inputStyle = ui.StyleInputUnfocused
	}

	// Show placeholder styling if empty
	if t.Value == "" && !t.Focused {
		b.WriteString(inputStyle.Foreground(ui.ColorMediumGray).Render(displayText))
	} else {
		b.WriteString(inputStyle.Render(displayText))
	}

	// Render error if present
	if t.Error != "" {
		b.WriteString("\n")
		errorStyle := lipgloss.NewStyle().
			Foreground(ui.ColorError)
		b.WriteString(errorStyle.Render(ui.SymbolCross + " " + t.Error))
	}

	return b.String()
}

// BlinkCmd returns a command for cursor blinking
func (t *TextInput) BlinkCmd() tea.Cmd {
	return tea.Tick(ui.CursorBlinkRate, func(_ time.Time) tea.Msg {
		return CursorBlinkMsg{}
	})
}

// Reset clears the input value
func (t *TextInput) Reset() {
	t.Value = ""
	t.cursor = 0
	t.Error = ""
}

// CursorBlinkMsg is sent to blink the cursor
type CursorBlinkMsg struct{}
