package services

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/vexoa/inkwash/internal/domain"
)

// serverServiceImpl implements ServerService
type serverServiceImpl struct {
	artifactService ArtifactService
	templateService TemplateService
	fileService     FileService
	serversDir      string
}

// NewServerService creates a new server service
func NewServerService(artifactService ArtifactService, templateService TemplateService, fileService FileService, serversDir string) ServerService {
	return &serverServiceImpl{
		artifactService: artifactService,
		templateService: templateService,
		fileService:     fileService,
		serversDir:      serversDir,
	}
}

// CreateServer creates a new FiveM server
func (s *serverServiceImpl) CreateServer(ctx context.Context, name, path, templateName string) (*domain.Server, error) {
	// Validate inputs
	if name == "" {
		return nil, domain.ErrInvalidServerConfig("server name cannot be empty")
	}
	if path == "" {
		return nil, domain.ErrInvalidServerConfig("server path cannot be empty")
	}

	// Check if server already exists
	if s.fileService.FileExists(path) {
		return nil, domain.ErrServerAlreadyExists(name)
	}

	// Create server instance
	server := domain.NewServer(name, path, templateName)

	// Get the template
	template, err := s.templateService.GetTemplate(ctx, templateName)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Get platform-appropriate artifact
	platform := domain.GetCurrentPlatform()
	artifact, err := s.artifactService.GetLatestArtifact(ctx, platform, domain.ArtifactChannelRecommended)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifact: %w", err)
	}

	server.Artifact = artifact

	// Create server directory
	if err := s.fileService.CreateDirectory(path, 0755); err != nil {
		return nil, err
	}

	// Download artifact if not cached
	progressCallback := func(current, total int64, message string) {
		fmt.Printf("\r%s: %d/%d bytes (%.1f%%)", message, current, total, float64(current)/float64(total)*100)
	}

	if !artifact.IsDownloaded() {
		fmt.Println("📥 Downloading FiveM artifacts...")
		if err := s.artifactService.DownloadArtifact(ctx, artifact, progressCallback); err != nil {
			s.fileService.DeleteFile(path) // Cleanup on failure
			return nil, fmt.Errorf("failed to download artifact: %w", err)
		}
		fmt.Println() // New line after progress
	}

	// Extract artifact to server directory
	fmt.Println("📦 Extracting server files...")
	if err := s.artifactService.ExtractArtifact(ctx, artifact, path); err != nil {
		s.fileService.DeleteFile(path) // Cleanup on failure
		return nil, fmt.Errorf("failed to extract artifact: %w", err)
	}

	// Apply template
	fmt.Println("🎨 Applying server template...")
	if err := s.templateService.ApplyTemplate(ctx, server.ID, templateName); err != nil {
		return nil, fmt.Errorf("failed to apply template: %w", err)
	}

	// Clean unnecessary files
	fmt.Println("🧹 Cleaning unnecessary files...")
	if err := s.cleanServerFiles(path); err != nil {
		return nil, fmt.Errorf("failed to clean server files: %w", err)
	}

	// Generate server.cfg
	fmt.Println("⚙️ Generating server configuration...")
	if err := s.generateServerConfig(server, template); err != nil {
		return nil, fmt.Errorf("failed to generate server config: %w", err)
	}

	// Save server metadata
	if err := s.saveServerMetadata(server); err != nil {
		return nil, fmt.Errorf("failed to save server metadata: %w", err)
	}

	fmt.Printf("✅ Successfully created FiveM server '%s' at %s\n", name, path)
	return server, nil
}

// GetServer retrieves a server by ID
func (s *serverServiceImpl) GetServer(ctx context.Context, serverID string) (*domain.Server, error) {
	metadataPath := filepath.Join(s.serversDir, serverID, ".inkwash", "server.json")
	
	data, err := s.fileService.ReadFile(metadataPath)
	if err != nil {
		return nil, domain.ErrServerNotFound(serverID)
	}

	var server domain.Server
	if err := json.Unmarshal(data, &server); err != nil {
		return nil, domain.NewError(domain.ErrorTypeInternal, "failed to parse server metadata").WithCause(err)
	}

	return &server, nil
}

// ListServers lists all servers
func (s *serverServiceImpl) ListServers(ctx context.Context) ([]*domain.Server, error) {
	if !s.fileService.FileExists(s.serversDir) {
		return []*domain.Server{}, nil
	}

	entries, err := s.fileService.ListDirectory(s.serversDir)
	if err != nil {
		return nil, err
	}

	var servers []*domain.Server
	for _, entry := range entries {
		server, err := s.GetServer(ctx, entry)
		if err != nil {
			continue // Skip invalid servers
		}
		servers = append(servers, server)
	}

	return servers, nil
}

// UpdateServer updates server configuration
func (s *serverServiceImpl) UpdateServer(ctx context.Context, server *domain.Server) error {
	return s.saveServerMetadata(server)
}

// DeleteServer deletes a server
func (s *serverServiceImpl) DeleteServer(ctx context.Context, serverID string) error {
	server, err := s.GetServer(ctx, serverID)
	if err != nil {
		return err
	}

	return s.fileService.DeleteFile(server.Path)
}

// StartServer starts a server
func (s *serverServiceImpl) StartServer(ctx context.Context, serverID string) error {
	// This would implement actual server starting logic
	// For now, just update the status
	server, err := s.GetServer(ctx, serverID)
	if err != nil {
		return err
	}

	server.Status = domain.ServerStatusRunning
	return s.UpdateServer(ctx, server)
}

// StopServer stops a server
func (s *serverServiceImpl) StopServer(ctx context.Context, serverID string) error {
	// This would implement actual server stopping logic
	// For now, just update the status
	server, err := s.GetServer(ctx, serverID)
	if err != nil {
		return err
	}

	server.Status = domain.ServerStatusStopped
	return s.UpdateServer(ctx, server)
}

// GetServerStatus gets the current status of a server
func (s *serverServiceImpl) GetServerStatus(ctx context.Context, serverID string) (domain.ServerStatus, error) {
	server, err := s.GetServer(ctx, serverID)
	if err != nil {
		return "", err
	}

	return server.Status, nil
}

// cleanServerFiles removes unnecessary files from the server directory
func (s *serverServiceImpl) cleanServerFiles(serverPath string) error {
	// List of files/directories to remove (similar to original logic)
	toRemove := []string{
		"citizen/system_resources/monitor",
		"citizen/system_resources/webadmin",
		"citizen/system_resources/baseevents",
	}

	for _, item := range toRemove {
		fullPath := filepath.Join(serverPath, item)
		if s.fileService.FileExists(fullPath) {
			if err := s.fileService.DeleteFile(fullPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// generateServerConfig generates the server.cfg file
func (s *serverServiceImpl) generateServerConfig(server *domain.Server, template *domain.Template) error {
	configContent := fmt.Sprintf(`# Server Configuration for %s
endpoint_add_tcp "0.0.0.0:%d"
endpoint_add_udp "0.0.0.0:%d"

sv_hostname "%s"
sv_maxclients %d

# License key (get one from https://keymaster.fivem.net)
sv_licenseKey ""

# Steam API key (optional)
steam_webApiKey ""

# Game build (optional)
sv_enforceGameBuild 2699

# Add resources
`, server.Name, server.Port, server.Port, server.Name, server.MaxPlayers)

	// Add template resources
	for _, resource := range template.Resources {
		configContent += fmt.Sprintf("start %s\n", resource)
	}

	// Add template config
	for key, value := range template.Config {
		configContent += fmt.Sprintf("%s \"%s\"\n", key, value)
	}

	// Write server.cfg
	configPath := filepath.Join(server.Path, "server.cfg")
	return s.fileService.WriteFile(configPath, []byte(configContent), 0644)
}

// saveServerMetadata saves server metadata to disk
func (s *serverServiceImpl) saveServerMetadata(server *domain.Server) error {
	metadataDir := filepath.Join(server.Path, ".inkwash")
	if err := s.fileService.CreateDirectory(metadataDir, 0755); err != nil {
		return err
	}

	metadataPath := filepath.Join(metadataDir, "server.json")
	data, err := json.MarshalIndent(server, "", "  ")
	if err != nil {
		return domain.NewError(domain.ErrorTypeInternal, "failed to marshal server metadata").WithCause(err)
	}

	return s.fileService.WriteFile(metadataPath, data, 0644)
}