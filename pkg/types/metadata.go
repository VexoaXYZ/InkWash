package types

import "time"

// ServerMetadata represents per-server metadata stored in metadata.json
type ServerMetadata struct {
	Version   int               `json:"version"` // Schema version for future migrations
	Build     BuildMetadata     `json:"build"`
	Lifecycle LifecycleMetadata `json:"lifecycle"`
	Stats     UsageStats        `json:"stats"`
}

// BuildMetadata tracks the installed FXServer build
type BuildMetadata struct {
	Number      int       `json:"number"`       // Build number (e.g., 17000)
	Hash        string    `json:"hash"`         // Build hash
	InstalledAt time.Time `json:"installed_at"` // When binaries were installed
	Recommended bool      `json:"recommended"`  // Was this a recommended build?
	Optional    bool      `json:"optional"`     // Was this an optional build?
}

// LifecycleMetadata tracks server lifecycle events
type LifecycleMetadata struct {
	CreatedAt   time.Time  `json:"created_at"`    // When server was created
	LastStarted *time.Time `json:"last_started"`  // Last time server was started (nil if never)
	LastStopped *time.Time `json:"last_stopped"`  // Last time server was stopped
}

// UsageStats tracks server usage statistics
type UsageStats struct {
	RestartCount int           `json:"restart_count"` // Number of times started
	TotalUptime  time.Duration `json:"total_uptime"`  // Total uptime in nanoseconds
}

// NewServerMetadata creates metadata for a freshly created server
func NewServerMetadata(build Build) *ServerMetadata {
	now := time.Now()
	return &ServerMetadata{
		Version: 1,
		Build: BuildMetadata{
			Number:      build.Number,
			Hash:        build.Hash,
			InstalledAt: now,
			Recommended: build.Recommended,
			Optional:    build.Optional,
		},
		Lifecycle: LifecycleMetadata{
			CreatedAt:   now,
			LastStarted: nil,
			LastStopped: nil,
		},
		Stats: UsageStats{
			RestartCount: 0,
			TotalUptime:  0,
		},
	}
}
