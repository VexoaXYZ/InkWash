package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/VexoaXYZ/inkwash/internal/cache"
	"github.com/VexoaXYZ/inkwash/internal/download"
	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/pkg/types"
)

// InstallProgress represents installation progress
type InstallProgress struct {
	Step            string
	Progress        float64
	DownloadSpeed   float64
	DownloadETA     time.Duration
	CurrentFile     string
	TotalSteps      int
	CompletedSteps  int
}

// ProgressCallback is called during installation
type ProgressCallback func(InstallProgress)

// Installer orchestrates server installation
type Installer struct {
	artifactClient *download.ArtifactClient
	downloader     *download.Downloader
	extractor      *download.Extractor
	cache          *cache.BinaryCache
	registry       *registry.Registry
	configGen      *ConfigGenerator
}

// NewInstaller creates a new installer
func NewInstaller(cache *cache.BinaryCache, registry *registry.Registry) *Installer {
	return &Installer{
		artifactClient: download.NewArtifactClient(),
		downloader:     download.NewDownloader(3),
		extractor:      download.NewExtractor(),
		cache:          cache,
		registry:       registry,
		configGen:      NewConfigGenerator(),
	}
}

// slugifyServerName converts a server name to a safe folder name
// Example: "Vexoa Test Server" -> "vexoa-test-server"
func slugifyServerName(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove any characters that aren't alphanumeric or hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// ensureUniqueFolderName ensures the folder name doesn't already exist
// If it does, appends a number to make it unique
func ensureUniqueFolderName(basePath, folderName string) string {
	targetPath := filepath.Join(basePath, folderName)

	// If it doesn't exist, we're good
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return folderName
	}

	// Otherwise, try appending numbers until we find a unique name
	counter := 1
	for {
		uniqueName := fmt.Sprintf("%s-%d", folderName, counter)
		targetPath = filepath.Join(basePath, uniqueName)

		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			return uniqueName
		}

		counter++
		if counter > 100 {
			// Prevent infinite loop, use timestamp
			return fmt.Sprintf("%s-%d", folderName, time.Now().Unix())
		}
	}
}

// Install installs a new FiveM server
func (inst *Installer) Install(
	serverName string,
	installPath string,
	buildNumber int,
	licenseKey string,
	port int,
	onProgress ProgressCallback,
) error {
	totalSteps := 8

	// Step 1: Validate inputs
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Validating configuration",
		Progress:       0,
		TotalSteps:     totalSteps,
		CompletedSteps: 0,
	})

	if err := inst.validateInputs(serverName, installPath); err != nil {
		return err
	}

	// Convert server name to slug for folder name
	// This ensures filesystem safety: "Vexoa Test Server" -> "vexoa-test-server"
	folderSlug := slugifyServerName(serverName)
	if folderSlug == "" {
		folderSlug = "fivem-server" // Fallback for invalid names
	}

	// Ensure the folder name is unique
	folderSlug = ensureUniqueFolderName(installPath, folderSlug)

	// Step 2: Create directory structure
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Creating directories",
		Progress:       0.14,
		TotalSteps:     totalSteps,
		CompletedSteps: 1,
	})

	serverPath := filepath.Join(installPath, folderSlug)
	binaryPath := filepath.Join(serverPath, "bin")

	if err := inst.createDirectories(serverPath, binaryPath); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Step 3: Get or download FXServer build
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Checking cache for FXServer build",
		Progress:       0.28,
		TotalSteps:     totalSteps,
		CompletedSteps: 2,
	})

	targetBuild, err := inst.installBinary(buildNumber, binaryPath, onProgress)
	if err != nil {
		return fmt.Errorf("failed to install FXServer: %w", err)
	}

	// Step 4: Clone server-data repository
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Cloning cfx-server-data",
		Progress:       0.57,
		TotalSteps:     totalSteps,
		CompletedSteps: 4,
	})

	if err := inst.cloneServerData(serverPath); err != nil {
		return fmt.Errorf("failed to clone server-data: %w", err)
	}

	// Step 5: Create metadata.json
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Creating server metadata",
		Progress:       0.625,
		TotalSteps:     totalSteps,
		CompletedSteps: 5,
	})

	metadataManager := NewMetadataManager()
	metadata := types.NewServerMetadata(*targetBuild)
	if err := metadataManager.Save(serverPath, metadata); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	// Step 6: Generate server.cfg
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Generating server.cfg",
		Progress:       0.75,
		TotalSteps:     totalSteps,
		CompletedSteps: 6,
	})

	server := &types.Server{
		Name:    serverName,
		Path:    serverPath,
		Port:    port,
		Created: time.Now(),
	}

	if err := inst.configGen.GenerateServerConfig(server, licenseKey); err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Step 7: Create launch script
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Creating launch script",
		Progress:       0.875,
		TotalSteps:     totalSteps,
		CompletedSteps: 7,
	})

	if err := inst.configGen.GenerateLaunchScript(server); err != nil {
		return fmt.Errorf("failed to create launch script: %w", err)
	}

	// Step 8: Register server
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Registering server",
		Progress:       1.0,
		TotalSteps:     totalSteps,
		CompletedSteps: 8,
	})

	if err := inst.registry.Add(*server); err != nil {
		return fmt.Errorf("failed to register server: %w", err)
	}

	return nil
}

// installBinary installs the FXServer binary and returns the Build info
func (inst *Installer) installBinary(buildNumber int, binaryPath string, onProgress ProgressCallback) (*types.Build, error) {
	// Fetch available builds first (needed for metadata even if cached)
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Fetching build information",
		Progress:       0.30,
		TotalSteps:     7,
		CompletedSteps: 2,
	})

	builds, err := inst.artifactClient.FetchBuilds()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch builds: %w", err)
	}

	// Find the requested build
	var targetBuild *types.Build
	for _, build := range builds {
		if build.Number == buildNumber {
			targetBuild = &build
			break
		}
	}

	if targetBuild == nil {
		return nil, fmt.Errorf("build %d not found", buildNumber)
	}

	// Check cache after getting build info
	cachedPath, err := inst.cache.Get(buildNumber)
	if err == nil {
		// Copy from cache
		inst.reportProgress(onProgress, InstallProgress{
			Step:           "Copying from cache",
			Progress:       0.35,
			CurrentFile:    fmt.Sprintf("Build %d (cached)", buildNumber),
			TotalSteps:     7,
			CompletedSteps: 2,
		})

		if err := copyDir(cachedPath, binaryPath); err != nil {
			return nil, err
		}
		return targetBuild, nil
	}

	// Download
	downloadURL := inst.artifactClient.GetDownloadURL(*targetBuild)
	tmpDir := filepath.Join(os.TempDir(), "inkwash-download")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, "server"+download.GetPlatformArchiveExtension())

	err = inst.downloader.Download(downloadURL, archivePath, func(p download.Progress) {
		downloadProgress := float64(p.DownloadedBytes) / float64(p.TotalBytes) * 0.15
		inst.reportProgress(onProgress, InstallProgress{
			Step:           "Downloading FXServer",
			Progress:       0.30 + downloadProgress,
			DownloadSpeed:  p.Speed,
			DownloadETA:    p.ETA,
			CurrentFile:    fmt.Sprintf("Build %d", buildNumber),
			TotalSteps:     7,
			CompletedSteps: 3,
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}

	// Extract
	inst.reportProgress(onProgress, InstallProgress{
		Step:           "Extracting archive",
		Progress:       0.45,
		TotalSteps:     7,
		CompletedSteps: 3,
	})

	extractPath := filepath.Join(tmpDir, "extracted")
	if err := inst.extractor.Extract(archivePath, extractPath); err != nil {
		return nil, fmt.Errorf("failed to extract: %w", err)
	}

	// Find the actual binary directory (may be nested like alpine/)
	sourcePath := findBinaryDir(extractPath)

	// Copy to destination
	if err := copyDirSkipBrokenSymlinks(sourcePath, binaryPath); err != nil {
		return nil, fmt.Errorf("failed to copy files: %w", err)
	}

	// Add to cache
	inst.cache.Add(*targetBuild, archivePath, extractPath)

	return targetBuild, nil
}

// cloneServerData clones the cfx-server-data repository
func (inst *Installer) cloneServerData(serverPath string) error {
	// Clone to temporary directory
	tmpDir := filepath.Join(os.TempDir(), "inkwash-server-data")
	os.RemoveAll(tmpDir) // Clean up any previous clone
	defer os.RemoveAll(tmpDir)

	// Clone using git (suppress progress output for clean TUI)
	cmd := exec.Command("git", "clone", "--quiet", "--depth", "1", "https://github.com/citizenfx/cfx-server-data.git", tmpDir)
	// Suppress output to avoid breaking TUI
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		// If git fails, create basic structure manually
		return inst.createBasicStructure(serverPath)
	}

	// Copy resources folder from cloned repo
	srcResources := filepath.Join(tmpDir, "resources")
	dstResources := filepath.Join(serverPath, "resources")

	if err := copyDir(srcResources, dstResources); err != nil {
		return fmt.Errorf("failed to copy resources: %w", err)
	}

	return nil
}

// createBasicStructure creates a basic server structure without git
func (inst *Installer) createBasicStructure(serverPath string) error {
	// Create basic directories
	dirs := []string{
		filepath.Join(serverPath, "resources"),
		filepath.Join(serverPath, "cache"),
		filepath.Join(serverPath, "logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// validateInputs validates installation inputs
func (inst *Installer) validateInputs(serverName, installPath string) error {
	// Check if server name is valid
	if serverName == "" {
		return fmt.Errorf("server name cannot be empty")
	}

	// Check if server already exists
	if inst.registry.Exists(serverName) {
		return fmt.Errorf("server '%s' already exists", serverName)
	}

	// Check if install path is writable
	testFile := filepath.Join(installPath, ".inkwash-test")
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return fmt.Errorf("cannot create install directory: %w", err)
	}

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("install path not writable: %w", err)
	}
	os.Remove(testFile)

	return nil
}

// createDirectories creates the directory structure
func (inst *Installer) createDirectories(serverPath, binaryPath string) error {
	dirs := []string{serverPath, binaryPath}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// reportProgress reports progress to callback
func (inst *Installer) reportProgress(callback ProgressCallback, progress InstallProgress) {
	if callback != nil {
		callback(progress)
	}
}

// Helper function to copy directory
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0755)
}

// GetPlatform returns the current platform
func GetPlatform() string {
	return runtime.GOOS
}

// findBinaryDir finds the actual binary directory within an extracted archive
// Some archives (like Linux FXServer) have a nested directory structure (e.g., alpine/)
func findBinaryDir(extractPath string) string {
	entries, err := os.ReadDir(extractPath)
	if err != nil {
		return extractPath
	}

	// If there's only one directory entry and no files, use that directory
	if len(entries) == 1 && entries[0].IsDir() {
		nestedPath := filepath.Join(extractPath, entries[0].Name())
		// Verify it contains actual files (not just another single directory)
		nestedEntries, err := os.ReadDir(nestedPath)
		if err == nil && len(nestedEntries) > 1 {
			return nestedPath
		}
	}

	return extractPath
}

// copyDirSkipBrokenSymlinks copies a directory, skipping broken symlinks
func copyDirSkipBrokenSymlinks(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip files we can't access
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		// Handle symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(path)
			if err != nil {
				// Skip broken symlinks
				return nil
			}

			// Check if the symlink target exists
			var targetPath string
			if filepath.IsAbs(linkTarget) {
				targetPath = linkTarget
			} else {
				targetPath = filepath.Join(filepath.Dir(path), linkTarget)
			}

			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				// Skip broken symlinks (target doesn't exist)
				return nil
			}

			// Create the symlink
			os.MkdirAll(filepath.Dir(dstPath), 0755)
			os.Remove(dstPath) // Remove if exists
			return os.Symlink(linkTarget, dstPath)
		}

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}
