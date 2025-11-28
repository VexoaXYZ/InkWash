package ui

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/mem"
)

// AnimationTier represents the level of animations to display
type AnimationTier int

const (
	// TierMinimal - Basic terminals, slow systems
	TierMinimal AnimationTier = iota
	// TierBalanced - Default for most users
	TierBalanced
	// TierFull - Modern terminals, fast systems
	TierFull
)

// DetectAnimationTier determines the optimal animation tier based on system capabilities
func DetectAnimationTier() AnimationTier {
	// Check 1: Terminal capabilities
	if !supportsANSI256() {
		return TierMinimal
	}

	// Check 2: System performance (CPU cores, available RAM)
	if runtime.NumCPU() < 4 || !hasEnoughRAM() {
		return TierBalanced
	}

	// Check 3: Terminal emulator detection
	if !isModernTerminal() {
		return TierBalanced
	}

	return TierFull
}

// supportsANSI256 checks if the terminal supports 256-color ANSI
func supportsANSI256() bool {
	term := os.Getenv("TERM")

	// Check for common 256-color terminal types
	if strings.Contains(term, "256color") {
		return true
	}

	// Check COLORTERM for truecolor support
	colorterm := os.Getenv("COLORTERM")
	if colorterm == "truecolor" || colorterm == "24bit" {
		return true
	}

	// Fallback: check if TERM is set at all
	return term != "" && term != "dumb"
}

// hasEnoughRAM checks if system has at least 2GB available RAM
func hasEnoughRAM() bool {
	v, err := mem.VirtualMemory()
	if err != nil {
		// Assume sufficient RAM if we can't detect
		return true
	}

	// 2GB in bytes
	const minRAM = 2 * 1024 * 1024 * 1024
	return v.Available >= minRAM
}

// isModernTerminal detects if we're running in a modern terminal emulator
func isModernTerminal() bool {
	termProgram := os.Getenv("TERM_PROGRAM")

	modernTerminals := []string{
		"iTerm.app",
		"WezTerm",
		"Alacritty",
		"Windows Terminal",
		"Hyper",
		"Tabby",
	}

	for _, modern := range modernTerminals {
		if termProgram == modern {
			return true
		}
	}

	// Check for Windows Terminal via WT_SESSION
	if os.Getenv("WT_SESSION") != "" {
		return true
	}

	return false
}

// GetTerminalSize returns the width and height of the terminal
func GetTerminalSize() (width, height int) {
	// Try to get from environment first
	if w := os.Getenv("COLUMNS"); w != "" {
		width, _ = strconv.Atoi(w)
	}
	if h := os.Getenv("LINES"); h != "" {
		height, _ = strconv.Atoi(h)
	}

	// Default fallback
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 24
	}

	return width, height
}

// String returns the string representation of AnimationTier
func (t AnimationTier) String() string {
	switch t {
	case TierMinimal:
		return "minimal"
	case TierBalanced:
		return "balanced"
	case TierFull:
		return "full"
	default:
		return "unknown"
	}
}
