package types

import "time"

// Build represents a FiveM server build
type Build struct {
	Number      int       `json:"number"`
	Hash        string    `json:"hash"`
	Timestamp   time.Time `json:"timestamp"`
	Recommended bool      `json:"recommended"`
	Optional    bool      `json:"optional"`
	Size        int64     `json:"size"`
}

// Label returns a human-readable label for the build
func (b *Build) Label() string {
	if b.Recommended {
		return "Recommended"
	}
	if b.Optional {
		return "Optional"
	}
	return "Latest"
}

// DownloadURL returns the download URL for the build
func (b *Build) DownloadURL(platform string) string {
	var baseURL string
	var filename string

	switch platform {
	case "windows":
		baseURL = "https://runtime.fivem.net/artifacts/fivem/build_server_windows/master"
		filename = "server.7z"
	case "linux":
		baseURL = "https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master"
		filename = "fx.tar.xz"
	default:
		return ""
	}

	return baseURL + "/" + b.Hash + "/" + filename
}
