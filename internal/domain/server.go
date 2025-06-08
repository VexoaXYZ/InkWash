package domain

import (
	"time"
)

// ServerStatus represents the current state of a FiveM server
type ServerStatus string

const (
	ServerStatusStopped  ServerStatus = "stopped"
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusRunning  ServerStatus = "running"
	ServerStatusStopping ServerStatus = "stopping"
	ServerStatusError    ServerStatus = "error"
)

// Server represents a FiveM server instance
type Server struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Status      ServerStatus      `json:"status"`
	Port        int               `json:"port"`
	MaxPlayers  int               `json:"max_players"`
	Template    string            `json:"template"`
	Artifact    *Artifact         `json:"artifact"`
	Resources   []Resource        `json:"resources"`
	Config      map[string]string `json:"config"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// NewServer creates a new server instance
func NewServer(name, path, template string) *Server {
	return &Server{
		ID:         generateID(),
		Name:       name,
		Path:       path,
		Status:     ServerStatusStopped,
		Port:       30120,
		MaxPlayers: 32,
		Template:   template,
		Config:     make(map[string]string),
		Resources:  []Resource{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// AddResource adds a resource to the server
func (s *Server) AddResource(resource Resource) {
	s.Resources = append(s.Resources, resource)
	s.UpdatedAt = time.Now()
}

// RemoveResource removes a resource from the server
func (s *Server) RemoveResource(resourceName string) {
	filtered := []Resource{}
	for _, r := range s.Resources {
		if r.Name != resourceName {
			filtered = append(filtered, r)
		}
	}
	s.Resources = filtered
	s.UpdatedAt = time.Now()
}

// SetConfig updates a configuration value
func (s *Server) SetConfig(key, value string) {
	s.Config[key] = value
	s.UpdatedAt = time.Now()
}

// GetConfig retrieves a configuration value
func (s *Server) GetConfig(key string) (string, bool) {
	value, exists := s.Config[key]
	return value, exists
}

// IsRunning checks if the server is running
func (s *Server) IsRunning() bool {
	return s.Status == ServerStatusRunning
}

// generateID creates a unique identifier for the server
func generateID() string {
	// Simple implementation, could be improved with UUID
	return time.Now().Format("20060102150405")
}