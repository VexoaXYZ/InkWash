package domain

import "time"

// ResourceType represents the type of FiveM resource
type ResourceType string

const (
	ResourceTypeGamemode  ResourceType = "gamemode"
	ResourceTypeScript    ResourceType = "script"
	ResourceTypeMap       ResourceType = "map"
	ResourceTypeVehicle   ResourceType = "vehicle"
	ResourceTypeWeapon    ResourceType = "weapon"
	ResourceTypeUI        ResourceType = "ui"
	ResourceTypeLibrary   ResourceType = "library"
	ResourceTypeUnknown   ResourceType = "unknown"
)

// ResourceStatus represents the current state of a resource
type ResourceStatus string

const (
	ResourceStatusInstalled   ResourceStatus = "installed"
	ResourceStatusDownloading ResourceStatus = "downloading"
	ResourceStatusUpdating    ResourceStatus = "updating"
	ResourceStatusError       ResourceStatus = "error"
	ResourceStatusDisabled    ResourceStatus = "disabled"
)

// Resource represents a FiveM resource
type Resource struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Type         ResourceType      `json:"type"`
	Status       ResourceStatus    `json:"status"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	Repository   string            `json:"repository"`
	Dependencies []string          `json:"dependencies"`
	Config       map[string]string `json:"config"`
	Path         string            `json:"path"`
	Enabled      bool              `json:"enabled"`
	InstalledAt  time.Time         `json:"installed_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// NewResource creates a new resource instance
func NewResource(name string, resourceType ResourceType) *Resource {
	return &Resource{
		Name:         name,
		Type:         resourceType,
		Status:       ResourceStatusInstalled,
		Dependencies: []string{},
		Config:       make(map[string]string),
		Enabled:      true,
		InstalledAt:  time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Enable enables the resource
func (r *Resource) Enable() {
	r.Enabled = true
	r.UpdatedAt = time.Now()
}

// Disable disables the resource
func (r *Resource) Disable() {
	r.Enabled = false
	r.UpdatedAt = time.Now()
}

// HasDependency checks if the resource has a specific dependency
func (r *Resource) HasDependency(dep string) bool {
	for _, d := range r.Dependencies {
		if d == dep {
			return true
		}
	}
	return false
}

// AddDependency adds a dependency to the resource
func (r *Resource) AddDependency(dep string) {
	if !r.HasDependency(dep) {
		r.Dependencies = append(r.Dependencies, dep)
		r.UpdatedAt = time.Now()
	}
}

// SetConfig sets a configuration value for the resource
func (r *Resource) SetConfig(key, value string) {
	r.Config[key] = value
	r.UpdatedAt = time.Now()
}

// GetConfig retrieves a configuration value
func (r *Resource) GetConfig(key string) (string, bool) {
	value, exists := r.Config[key]
	return value, exists
}