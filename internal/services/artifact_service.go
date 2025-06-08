package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/vexoa/inkwash/internal/domain"
)

// artifactServiceImpl implements ArtifactService
type artifactServiceImpl struct {
	cacheDir        string
	downloadService DownloadService
	fileService     FileService
	httpClient      *http.Client
}

// NewArtifactService creates a new artifact service
func NewArtifactService(cacheDir string, downloadService DownloadService, fileService FileService) ArtifactService {
	return &artifactServiceImpl{
		cacheDir:        cacheDir,
		downloadService: downloadService,
		fileService:     fileService,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetLatestArtifact gets the latest artifact for a platform
func (s *artifactServiceImpl) GetLatestArtifact(ctx context.Context, platform domain.ArtifactPlatform, channel domain.ArtifactChannel) (*domain.Artifact, error) {
	// Get the base URL for the platform
	baseURL, err := s.getBaseURL(platform)
	if err != nil {
		return nil, err
	}

	// Fetch the latest build number
	buildNumber, buildHash, err := s.getLatestBuild(ctx, baseURL)
	if err != nil {
		return nil, err
	}

	// Create artifact with discovered build info
	artifact := domain.NewArtifact("latest", buildNumber, platform, channel)
	
	// Construct the download URL
	downloadURL, err := s.constructDownloadURL(baseURL, buildNumber, buildHash, platform)
	if err != nil {
		return nil, err
	}
	
	artifact.DownloadURL = downloadURL
	return artifact, nil
}

// DownloadArtifact downloads an artifact
func (s *artifactServiceImpl) DownloadArtifact(ctx context.Context, artifact *domain.Artifact, progress ProgressCallback) error {
	if artifact.IsDownloaded() {
		return nil // Already downloaded
	}

	// Create cache directory
	if err := s.fileService.CreateDirectory(s.cacheDir, 0755); err != nil {
		return domain.ErrFilesystemOperation("create_cache_dir", s.cacheDir, err)
	}

	// Generate cache path
	filename := fmt.Sprintf("fivem_%s_%s_%s", artifact.Version, artifact.BuildNumber, artifact.Platform)
	if artifact.Platform == domain.ArtifactPlatformLinux {
		filename += ".tar.xz"
	} else {
		filename += ".zip"
	}
	cachePath := filepath.Join(s.cacheDir, filename)

	// Download the artifact
	if err := s.downloadService.Download(ctx, artifact.GetDownloadURL(), cachePath, progress); err != nil {
		return domain.ErrDownloadFailed(artifact.GetDownloadURL(), err)
	}

	// Calculate checksum
	checksum, err := s.calculateChecksum(cachePath)
	if err != nil {
		return domain.ErrFilesystemOperation("calculate_checksum", cachePath, err)
	}

	// Update artifact
	artifact.Checksum = checksum
	artifact.MarkAsDownloaded(cachePath)

	// Get file size
	info, err := s.fileService.GetFileInfo(cachePath)
	if err == nil {
		artifact.Size = info.Size()
	}

	return nil
}

// ExtractArtifact extracts an artifact to a directory
func (s *artifactServiceImpl) ExtractArtifact(ctx context.Context, artifact *domain.Artifact, destPath string) error {
	if !artifact.IsDownloaded() {
		return domain.NewError(domain.ErrorTypeValidation, "artifact not downloaded")
	}

	// Create destination directory
	if err := s.fileService.CreateDirectory(destPath, 0755); err != nil {
		return domain.ErrFilesystemOperation("create_dest_dir", destPath, err)
	}

	// Extract based on platform
	switch artifact.Platform {
	case domain.ArtifactPlatformLinux:
		return s.extractTarXz(artifact.CachePath, destPath)
	case domain.ArtifactPlatformWindows:
		return s.extractZip(artifact.CachePath, destPath)
	default:
		return domain.NewError(domain.ErrorTypeValidation, "unsupported platform")
	}
}

// ListCachedArtifacts lists all cached artifacts
func (s *artifactServiceImpl) ListCachedArtifacts(ctx context.Context) ([]*domain.Artifact, error) {
	files, err := s.fileService.ListDirectory(s.cacheDir)
	if err != nil {
		return nil, domain.ErrFilesystemOperation("list_cache", s.cacheDir, err)
	}

	var artifacts []*domain.Artifact
	for _, file := range files {
		// Parse artifact info from filename
		// This is simplified - in practice, you'd store metadata separately
		if filepath.Ext(file) == ".xz" || filepath.Ext(file) == ".zip" {
			// Create a basic artifact entry
			artifact := &domain.Artifact{
				CachePath: filepath.Join(s.cacheDir, file),
			}
			
			// Get file info
			info, err := s.fileService.GetFileInfo(artifact.CachePath)
			if err == nil {
				artifact.Size = info.Size()
				downloadTime := info.ModTime()
				artifact.DownloadedAt = &downloadTime
			}
			
			artifacts = append(artifacts, artifact)
		}
	}

	return artifacts, nil
}

// CleanCache cleans old artifacts from cache
func (s *artifactServiceImpl) CleanCache(ctx context.Context, keepLatest int) error {
	artifacts, err := s.ListCachedArtifacts(ctx)
	if err != nil {
		return err
	}

	if len(artifacts) <= keepLatest {
		return nil // Nothing to clean
	}

	// Sort by download time (newest first)
	// For simplicity, we'll just delete the oldest files
	toDelete := artifacts[keepLatest:]
	
	for _, artifact := range toDelete {
		if err := s.fileService.DeleteFile(artifact.CachePath); err != nil {
			return domain.ErrFilesystemOperation("delete_cache_file", artifact.CachePath, err)
		}
	}

	return nil
}

// VerifyArtifact verifies artifact integrity
func (s *artifactServiceImpl) VerifyArtifact(ctx context.Context, artifact *domain.Artifact) error {
	if !artifact.IsDownloaded() {
		return domain.NewError(domain.ErrorTypeValidation, "artifact not downloaded")
	}

	// Calculate current checksum
	currentChecksum, err := s.calculateChecksum(artifact.CachePath)
	if err != nil {
		return domain.ErrFilesystemOperation("calculate_checksum", artifact.CachePath, err)
	}

	// Compare with stored checksum
	if artifact.Checksum != "" && currentChecksum != artifact.Checksum {
		return domain.NewError(domain.ErrorTypeValidation, "artifact checksum mismatch")
	}

	return nil
}

// calculateChecksum calculates SHA256 checksum of a file
func (s *artifactServiceImpl) calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// extractTarXz extracts a tar.xz file using system tar command
func (s *artifactServiceImpl) extractTarXz(srcPath, destPath string) error {
	// Use system tar command for now - this requires tar with xz support
	// In production, you'd want to use pure Go libraries
	cmd := fmt.Sprintf("tar -xf %s -C %s", srcPath, destPath)
	
	// For basic implementation, we'll use os/exec
	// This is not ideal but works for demonstration
	return s.executeCommand(cmd)
}

// extractZip extracts a zip file (simplified - would use archive/zip)
func (s *artifactServiceImpl) extractZip(srcPath, destPath string) error {
	// This is a placeholder - in a real implementation, you'd use
	// archive/zip package
	return fmt.Errorf("zip extraction not implemented yet")
}

// getBaseURL returns the base URL for a platform
func (s *artifactServiceImpl) getBaseURL(platform domain.ArtifactPlatform) (string, error) {
	switch platform {
	case domain.ArtifactPlatformLinux:
		return "https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/", nil
	case domain.ArtifactPlatformWindows:
		return "https://runtime.fivem.net/artifacts/fivem/build_server_windows/master/", nil
	default:
		return "", domain.NewError(domain.ErrorTypeValidation, "unsupported platform")
	}
}

// getLatestBuild fetches the latest build number and hash from the FiveM artifacts page
func (s *artifactServiceImpl) getLatestBuild(ctx context.Context, baseURL string) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return "", "", domain.ErrDownloadFailed(baseURL, err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", "", domain.ErrDownloadFailed(baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", domain.ErrDownloadFailed(baseURL, fmt.Errorf("HTTP %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", domain.ErrDownloadFailed(baseURL, err)
	}

	return s.parseBuildFromHTML(string(body))
}

// parseBuildFromHTML parses the HTML to extract build numbers and finds the latest one
func (s *artifactServiceImpl) parseBuildFromHTML(html string) (string, string, error) {
	// Look for the LATEST RECOMMENDED build first
	// Example: <a href= "./7290-a654bcc2adfa27c4e020fc915a1a6343c3b4f921/fx.tar.xz" class="button is-link is-primary">
	recommendedRegex := regexp.MustCompile(`href= "\./(\d+)-([a-f0-9]+)/[^"]*" class="button is-link is-primary"`)
	if matches := recommendedRegex.FindStringSubmatch(html); len(matches) >= 3 {
		return matches[1], matches[2], nil
	}

	// Fallback: Look for any build in the panel blocks
	// Example: <a class="panel-block" href="./15744-8682969ff3e99a09330b5fda5c9947f443455cac/fx.tar.xz"
	buildRegex := regexp.MustCompile(`href="\./(\d+)-([a-f0-9]+)/[^"]*"`)
	matches := buildRegex.FindAllStringSubmatch(html, -1)

	if len(matches) == 0 {
		return "", "", domain.NewError(domain.ErrorTypeNotFound, "no builds found in artifacts page")
	}

	// Find the highest build number
	var latestBuild int
	var latestHash string

	for _, match := range matches {
		if len(match) >= 3 {
			buildNum, err := strconv.Atoi(match[1])
			if err != nil {
				continue // Skip invalid build numbers
			}

			if buildNum > latestBuild {
				latestBuild = buildNum
				latestHash = match[2]
			}
		}
	}

	if latestBuild == 0 {
		return "", "", domain.NewError(domain.ErrorTypeNotFound, "no valid builds found")
	}

	return strconv.Itoa(latestBuild), latestHash, nil
}

// constructDownloadURL constructs the full download URL for an artifact
func (s *artifactServiceImpl) constructDownloadURL(baseURL, buildNumber, buildHash string, platform domain.ArtifactPlatform) (string, error) {
	buildDir := fmt.Sprintf("%s-%s", buildNumber, buildHash)
	
	switch platform {
	case domain.ArtifactPlatformLinux:
		return fmt.Sprintf("%s%s/fx.tar.xz", baseURL, buildDir), nil
	case domain.ArtifactPlatformWindows:
		return fmt.Sprintf("%s%s/server.zip", baseURL, buildDir), nil
	default:
		return "", domain.NewError(domain.ErrorTypeValidation, "unsupported platform")
	}
}

// executeCommand executes a shell command
func (s *artifactServiceImpl) executeCommand(cmdStr string) error {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	if err := cmd.Run(); err != nil {
		return domain.ErrFilesystemOperation("execute_command", cmdStr, err)
	}

	return nil
}