package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/VexoaXYZ/inkwash/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

// ProgressBar represents a progress bar component
type ProgressBar struct {
	Width    int
	Progress float64 // 0.0 to 1.0
	Shimmer  bool
}

// NewProgressBar creates a new progress bar
func NewProgressBar(width int) *ProgressBar {
	return &ProgressBar{
		Width:   width,
		Shimmer: true,
	}
}

// SetProgress updates the progress (0.0 to 1.0)
func (p *ProgressBar) SetProgress(progress float64) {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	p.Progress = progress
}

// Render renders the progress bar
func (p *ProgressBar) Render() string {
	return p.RenderWithTier(ui.TierFull)
}

// RenderWithTier renders the progress bar based on animation tier
func (p *ProgressBar) RenderWithTier(tier ui.AnimationTier) string {
	filled := int(p.Progress * float64(p.Width))

	switch tier {
	case ui.TierMinimal:
		// Simple bar without shimmer
		bar := strings.Repeat("█", filled) + strings.Repeat("░", p.Width-filled)
		return lipgloss.NewStyle().Foreground(ui.ColorPrimary).Render(bar)

	case ui.TierBalanced:
		// Smooth fill without shimmer
		bar := strings.Repeat("█", filled) + strings.Repeat("░", p.Width-filled)
		filledPart := lipgloss.NewStyle().Foreground(ui.ColorPrimary).Render(strings.Repeat("█", filled))
		emptyPart := lipgloss.NewStyle().Foreground(ui.ColorLightGray).Render(strings.Repeat("░", p.Width-filled))
		return filledPart + emptyPart

	case ui.TierFull:
		// Full animation with shimmer
		if !p.Shimmer {
			return p.RenderWithTier(ui.TierBalanced)
		}

		// Shimmer position (cycles every 2s)
		shimmerPos := int(time.Now().UnixMilli()/20) % p.Width

		var bar strings.Builder
		for i := 0; i < p.Width; i++ {
			if i < filled {
				// Filled portion
				if abs(i-shimmerPos) < 3 {
					// Shimmer highlight
					bar.WriteString(lipgloss.NewStyle().
						Foreground(ui.ColorPrimaryGlow).Render("█"))
				} else {
					bar.WriteString(lipgloss.NewStyle().
						Foreground(ui.ColorPrimary).Render("█"))
				}
			} else {
				// Empty portion
				bar.WriteString(lipgloss.NewStyle().
					Foreground(ui.ColorLightGray).Render("░"))
			}
		}

		return bar.String()

	default:
		return p.RenderWithTier(ui.TierBalanced)
	}
}

// RenderWithLabel renders the progress bar with a percentage label
func (p *ProgressBar) RenderWithLabel() string {
	bar := p.Render()
	percentage := fmt.Sprintf("%.0f%%", p.Progress*100)

	return fmt.Sprintf("%s  %s", bar, ui.StyleTextMuted.Render(percentage))
}

// RenderWithStats renders the progress bar with speed and ETA
func (p *ProgressBar) RenderWithStats(speed string, eta string) string {
	bar := p.Render()
	percentage := fmt.Sprintf("%.0f%%", p.Progress*100)

	stats := fmt.Sprintf("%s  %s  %s",
		ui.StyleTextMuted.Render(percentage),
		ui.StyleAccent.Render(speed),
		ui.StyleTextMuted.Render("ETA "+eta),
	)

	return fmt.Sprintf("%s\n%s", bar, stats)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DownloadProgress represents a download progress component
type DownloadProgress struct {
	ProgressBar      *ProgressBar
	DownloadedBytes int64
	TotalBytes      int64
	Speed           float64 // MB/s
	ETA             time.Duration
}

// NewDownloadProgress creates a new download progress component
func NewDownloadProgress(totalBytes int64) *DownloadProgress {
	return &DownloadProgress{
		ProgressBar: NewProgressBar(40),
		TotalBytes:  totalBytes,
	}
}

// Update updates the download progress
func (d *DownloadProgress) Update(downloadedBytes int64, speed float64) {
	d.DownloadedBytes = downloadedBytes
	d.Speed = speed

	if d.TotalBytes > 0 {
		d.ProgressBar.SetProgress(float64(downloadedBytes) / float64(d.TotalBytes))
	}

	// Calculate ETA
	if speed > 0 {
		remaining := float64(d.TotalBytes - d.DownloadedBytes)
		etaSeconds := remaining / (speed * 1024 * 1024) // Convert MB/s to bytes/s
		d.ETA = time.Duration(etaSeconds) * time.Second
	}
}

// Render renders the download progress
func (d *DownloadProgress) Render() string {
	speedStr := fmt.Sprintf("%.1f MB/s", d.Speed)
	etaStr := d.ETA.Round(time.Second).String()

	return d.ProgressBar.RenderWithStats(speedStr, etaStr)
}
