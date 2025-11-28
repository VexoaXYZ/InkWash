package types

import (
	"path/filepath"
	"time"
)

// Server represents a FiveM server instance
type Server struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	// BinaryPath removed - now calculated as {Path}/bin
	// Build removed - now in metadata.json
	// BuildHash removed - now in metadata.json
	KeyID       string    `json:"key_id"`
	Port        int       `json:"port"`
	Created     time.Time `json:"created"`
	LastStarted time.Time `json:"last_started"`
	PID         int       `json:"pid"`
	AutoStart   bool      `json:"auto_start"`
}

// GetBinaryPath returns the path to the server's bin directory
func (s *Server) GetBinaryPath() string {
	return filepath.Join(s.Path, "bin")
}

// GetBinaryExecutable returns the platform-specific executable path
func (s *Server) GetBinaryExecutable() string {
	return filepath.Join(s.GetBinaryPath(), "FXServer.exe")
}

// IsRunning returns true if the server is currently running
func (s *Server) IsRunning() bool {
	return s.PID > 0
}

// Status returns a human-readable status string
func (s *Server) Status() string {
	if s.IsRunning() {
		return "Running"
	}
	return "Stopped"
}
