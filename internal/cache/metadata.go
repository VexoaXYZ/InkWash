package cache

import "time"

// Metadata represents cache metadata
type Metadata struct {
	Version   int            `json:"version"`
	Builds    []CachedBuild  `json:"builds"`
	MaxBuilds int            `json:"max_builds"`
	TotalSize int64          `json:"total_size"`
}

// CachedBuild represents a cached build entry
type CachedBuild struct {
	Number      int       `json:"number"`
	Hash        string    `json:"hash"`
	Downloaded  time.Time `json:"downloaded"`
	Size        int64     `json:"size"`
	Recommended bool      `json:"recommended"`
	Optional    bool      `json:"optional"`
	LastUsed    time.Time `json:"last_used"`
}

// CacheStats represents cache statistics
type CacheStats struct {
	TotalBuilds int
	TotalSize   int64
	MaxBuilds   int
}
