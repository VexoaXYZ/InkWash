package types

import "time"

// Server represents a FiveM server instance
type Server struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	BinaryPath  string    `json:"binary_path"`
	Build       int       `json:"build"`
	BuildHash   string    `json:"build_hash"`
	KeyID       string    `json:"key_id"`
	Port        int       `json:"port"`
	Created     time.Time `json:"created"`
	LastStarted time.Time `json:"last_started"`
	PID         int       `json:"pid"`
	AutoStart   bool      `json:"auto_start"`
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
