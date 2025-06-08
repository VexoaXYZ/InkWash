package services

import (
	"context"
	"os"

	"github.com/vexoa/inkwash/internal/domain"
)

// ServerService defines the interface for server management
type ServerService interface {
	// CreateServer creates a new FiveM server
	CreateServer(ctx context.Context, name, path, template string) (*domain.Server, error)
	
	// GetServer retrieves a server by ID
	GetServer(ctx context.Context, serverID string) (*domain.Server, error)
	
	// ListServers lists all servers
	ListServers(ctx context.Context) ([]*domain.Server, error)
	
	// UpdateServer updates server configuration
	UpdateServer(ctx context.Context, server *domain.Server) error
	
	// DeleteServer deletes a server
	DeleteServer(ctx context.Context, serverID string) error
	
	// StartServer starts a server
	StartServer(ctx context.Context, serverID string) error
	
	// StopServer stops a server
	StopServer(ctx context.Context, serverID string) error
	
	// GetServerStatus gets the current status of a server
	GetServerStatus(ctx context.Context, serverID string) (domain.ServerStatus, error)
}

// ResourceService defines the interface for resource management
type ResourceService interface {
	// InstallResource installs a resource to a server
	InstallResource(ctx context.Context, serverID string, resource *domain.Resource) error
	
	// RemoveResource removes a resource from a server
	RemoveResource(ctx context.Context, serverID string, resourceName string) error
	
	// UpdateResource updates a resource
	UpdateResource(ctx context.Context, serverID string, resource *domain.Resource) error
	
	// ListResources lists all resources for a server
	ListResources(ctx context.Context, serverID string) ([]*domain.Resource, error)
	
	// GetResource gets a specific resource
	GetResource(ctx context.Context, serverID string, resourceName string) (*domain.Resource, error)
	
	// EnableResource enables a resource
	EnableResource(ctx context.Context, serverID string, resourceName string) error
	
	// DisableResource disables a resource
	DisableResource(ctx context.Context, serverID string, resourceName string) error
	
	// ValidateResource validates resource configuration
	ValidateResource(ctx context.Context, resource *domain.Resource) error
}

// ArtifactService defines the interface for artifact management
type ArtifactService interface {
	// GetLatestArtifact gets the latest artifact for a platform
	GetLatestArtifact(ctx context.Context, platform domain.ArtifactPlatform, channel domain.ArtifactChannel) (*domain.Artifact, error)
	
	// DownloadArtifact downloads an artifact
	DownloadArtifact(ctx context.Context, artifact *domain.Artifact, progress ProgressCallback) error
	
	// ExtractArtifact extracts an artifact to a directory
	ExtractArtifact(ctx context.Context, artifact *domain.Artifact, destPath string) error
	
	// ListCachedArtifacts lists all cached artifacts
	ListCachedArtifacts(ctx context.Context) ([]*domain.Artifact, error)
	
	// CleanCache cleans old artifacts from cache
	CleanCache(ctx context.Context, keepLatest int) error
	
	// VerifyArtifact verifies artifact integrity
	VerifyArtifact(ctx context.Context, artifact *domain.Artifact) error
}

// TemplateService defines the interface for template management
type TemplateService interface {
	// GetTemplate gets a template by name
	GetTemplate(ctx context.Context, name string) (*domain.Template, error)
	
	// ListTemplates lists all available templates
	ListTemplates(ctx context.Context) ([]*domain.Template, error)
	
	// ApplyTemplate applies a template to a server
	ApplyTemplate(ctx context.Context, serverID string, templateName string) error
	
	// CreateTemplate creates a custom template
	CreateTemplate(ctx context.Context, template *domain.Template) error
	
	// UpdateTemplate updates a template
	UpdateTemplate(ctx context.Context, template *domain.Template) error
	
	// DeleteTemplate deletes a custom template
	DeleteTemplate(ctx context.Context, templateName string) error
	
	// ExportTemplate exports a server configuration as a template
	ExportTemplate(ctx context.Context, serverID string, templateName string) (*domain.Template, error)
}

// ProgressCallback is a function called during long operations to report progress
type ProgressCallback func(current, total int64, message string)

// FileService defines the interface for file operations
type FileService interface {
	// ReadFile reads a file
	ReadFile(path string) ([]byte, error)
	
	// WriteFile writes data to a file
	WriteFile(path string, data []byte, perm os.FileMode) error
	
	// CopyFile copies a file
	CopyFile(src, dst string) error
	
	// MoveFile moves a file
	MoveFile(src, dst string) error
	
	// DeleteFile deletes a file
	DeleteFile(path string) error
	
	// CreateDirectory creates a directory
	CreateDirectory(path string, perm os.FileMode) error
	
	// ListDirectory lists directory contents
	ListDirectory(path string) ([]string, error)
	
	// FileExists checks if a file exists
	FileExists(path string) bool
	
	// GetFileInfo gets file information
	GetFileInfo(path string) (os.FileInfo, error)
}

// DownloadService defines the interface for download operations
type DownloadService interface {
	// Download downloads a file from URL
	Download(ctx context.Context, url, destPath string, progress ProgressCallback) error
	
	// DownloadWithResume downloads with resume support
	DownloadWithResume(ctx context.Context, url, destPath string, progress ProgressCallback) error
	
	// GetContentLength gets the content length of a URL
	GetContentLength(ctx context.Context, url string) (int64, error)
}