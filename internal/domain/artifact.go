package domain

import (
	"fmt"
	"runtime"
	"time"
)

// ArtifactPlatform represents the target platform for a FiveM artifact
type ArtifactPlatform string

const (
	ArtifactPlatformLinux   ArtifactPlatform = "linux"
	ArtifactPlatformWindows ArtifactPlatform = "windows"
)

// ArtifactChannel represents the release channel
type ArtifactChannel string

const (
	ArtifactChannelRecommended ArtifactChannel = "recommended"
	ArtifactChannelLatest      ArtifactChannel = "latest"
	ArtifactChannelOptional    ArtifactChannel = "optional"
)

// Artifact represents a FiveM server artifact/build
type Artifact struct {
	Version      string           `json:"version"`
	BuildNumber  string           `json:"build_number"`
	Platform     ArtifactPlatform `json:"platform"`
	Channel      ArtifactChannel  `json:"channel"`
	DownloadURL  string           `json:"download_url"`
	Checksum     string           `json:"checksum"`
	Size         int64            `json:"size"`
	ReleaseDate  time.Time        `json:"release_date"`
	DownloadedAt *time.Time       `json:"downloaded_at,omitempty"`
	CachePath    string           `json:"cache_path,omitempty"`
}

// NewArtifact creates a new artifact instance
func NewArtifact(version, buildNumber string, platform ArtifactPlatform, channel ArtifactChannel) *Artifact {
	return &Artifact{
		Version:     version,
		BuildNumber: buildNumber,
		Platform:    platform,
		Channel:     channel,
		ReleaseDate: time.Now(),
	}
}

// GetDownloadURL constructs the download URL for the artifact
func (a *Artifact) GetDownloadURL() string {
	if a.DownloadURL != "" {
		return a.DownloadURL
	}
	
	// Construct FiveM artifact URL
	baseURL := "https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master"
	return fmt.Sprintf("%s/%s/fx.tar.xz", baseURL, a.BuildNumber)
}

// IsDownloaded checks if the artifact has been downloaded
func (a *Artifact) IsDownloaded() bool {
	return a.DownloadedAt != nil && a.CachePath != ""
}

// MarkAsDownloaded marks the artifact as downloaded
func (a *Artifact) MarkAsDownloaded(cachePath string) {
	now := time.Now()
	a.DownloadedAt = &now
	a.CachePath = cachePath
}

// GetPlatformString returns the platform as a string
func (a *Artifact) GetPlatformString() string {
	return string(a.Platform)
}

// GetCurrentPlatform returns the current system's platform as ArtifactPlatform
func GetCurrentPlatform() ArtifactPlatform {
	switch runtime.GOOS {
	case "linux":
		return ArtifactPlatformLinux
	case "windows":
		return ArtifactPlatformWindows
	default:
		return ArtifactPlatformLinux // Default to Linux
	}
}