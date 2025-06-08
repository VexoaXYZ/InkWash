package services

import (
	"os"
	"path/filepath"
)

// Container holds all services for dependency injection
type Container struct {
	FileService     FileService
	DownloadService DownloadService
	ArtifactService ArtifactService
	TemplateService TemplateService
	ServerService   ServerService
}

// NewContainer creates a new service container
func NewContainer() *Container {
	// Create base directories
	homeDir, _ := os.UserHomeDir()
	inkwashDir := filepath.Join(homeDir, ".inkwash")
	cacheDir := filepath.Join(inkwashDir, "cache")
	templatesDir := filepath.Join(inkwashDir, "templates")
	serversDir := filepath.Join(inkwashDir, "servers")

	// Create services in dependency order
	fileService := NewFileService()
	downloadService := NewDownloadService()
	artifactService := NewArtifactService(cacheDir, downloadService, fileService)
	templateService := NewTemplateService(fileService, templatesDir)
	serverService := NewServerService(artifactService, templateService, fileService, serversDir)

	return &Container{
		FileService:     fileService,
		DownloadService: downloadService,
		ArtifactService: artifactService,
		TemplateService: templateService,
		ServerService:   serverService,
	}
}