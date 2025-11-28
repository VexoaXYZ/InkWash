package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/VexoaXYZ/inkwash/pkg/types"
)

// Registry manages server instances
type Registry struct {
	configPath string
	data       *RegistryData
	mu         sync.RWMutex
}

// RegistryData represents the registry file structure
type RegistryData struct {
	Version int             `json:"version"`
	Servers []types.Server  `json:"servers"`
}

// NewRegistry creates a new registry
func NewRegistry(configPath string) (*Registry, error) {
	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	r := &Registry{
		configPath: configPath,
	}

	// Load or create registry
	if err := r.load(); err != nil {
		return nil, err
	}

	return r, nil
}

// Add adds a new server to the registry
func (r *Registry) Add(server types.Server) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if server already exists
	for _, s := range r.data.Servers {
		if s.Name == server.Name {
			return fmt.Errorf("server '%s' already exists", server.Name)
		}
	}

	r.data.Servers = append(r.data.Servers, server)
	return r.save()
}

// Remove removes a server from the registry
func (r *Registry) Remove(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, server := range r.data.Servers {
		if server.Name == name {
			r.data.Servers = append(r.data.Servers[:i], r.data.Servers[i+1:]...)
			return r.save()
		}
	}

	return fmt.Errorf("server '%s' not found", name)
}

// Get retrieves a server by name
func (r *Registry) Get(name string) (*types.Server, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i, server := range r.data.Servers {
		if server.Name == name {
			return &r.data.Servers[i], nil
		}
	}

	return nil, fmt.Errorf("server '%s' not found", name)
}

// List returns all servers
func (r *Registry) List() []types.Server {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modifications
	servers := make([]types.Server, len(r.data.Servers))
	copy(servers, r.data.Servers)
	return servers
}

// Update updates a server in the registry
func (r *Registry) Update(server types.Server) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, s := range r.data.Servers {
		if s.Name == server.Name {
			r.data.Servers[i] = server
			return r.save()
		}
	}

	return fmt.Errorf("server '%s' not found", server.Name)
}

// UpdatePID updates a server's PID
func (r *Registry) UpdatePID(name string, pid int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, server := range r.data.Servers {
		if server.Name == name {
			r.data.Servers[i].PID = pid
			return r.save()
		}
	}

	return fmt.Errorf("server '%s' not found", name)
}

// Exists checks if a server exists
func (r *Registry) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, server := range r.data.Servers {
		if server.Name == name {
			return true
		}
	}

	return false
}

// Count returns the number of servers
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.data.Servers)
}

// GetRunning returns all running servers
func (r *Registry) GetRunning() []types.Server {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var running []types.Server
	for _, server := range r.data.Servers {
		if server.IsRunning() {
			running = append(running, server)
		}
	}

	return running
}

// GetStopped returns all stopped servers
func (r *Registry) GetStopped() []types.Server {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var stopped []types.Server
	for _, server := range r.data.Servers {
		if !server.IsRunning() {
			stopped = append(stopped, server)
		}
	}

	return stopped
}

// load loads the registry from disk
func (r *Registry) load() error {
	// If registry doesn't exist, create empty
	if _, err := os.Stat(r.configPath); os.IsNotExist(err) {
		r.data = &RegistryData{
			Version: 1,
			Servers: []types.Server{},
		}
		return r.save()
	}

	// Load existing registry
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return fmt.Errorf("failed to read registry: %w", err)
	}

	var registryData RegistryData
	if err := json.Unmarshal(data, &registryData); err != nil {
		return fmt.Errorf("failed to parse registry: %w", err)
	}

	r.data = &registryData
	return nil
}

// save saves the registry to disk
func (r *Registry) save() error {
	data, err := json.MarshalIndent(r.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(r.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry: %w", err)
	}

	return nil
}

// Reload reloads the registry from disk
func (r *Registry) Reload() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.load()
}
