package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/vexoa/inkwash/internal/utils"
)

const (
	githubAPIURL = "https://api.github.com/repos/VexoaXYZ/InkWash/releases/latest"
	GithubRepo   = "VexoaXYZ/InkWash"
)

// Release represents a GitHub release
type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []Asset `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	Available   bool
	CurrentVersion string
	LatestVersion  string
	DownloadURL    string
	AssetName      string
	ReleaseNotes   string
}

// Updater handles checking and applying updates
type Updater struct {
	currentVersion string
	httpClient     *http.Client
}

// NewUpdater creates a new updater instance
func NewUpdater(currentVersion string) *Updater {
	return &Updater{
		currentVersion: currentVersion,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckForUpdate checks if a new version is available
func (u *Updater) CheckForUpdate() (*UpdateInfo, error) {
	release, err := u.fetchLatestRelease()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}

	info := &UpdateInfo{
		CurrentVersion: u.currentVersion,
		LatestVersion:  strings.TrimPrefix(release.TagName, "v"),
		ReleaseNotes:   release.Body,
	}

	// Compare versions
	if !u.isNewerVersion(info.LatestVersion) {
		info.Available = false
		return info, nil
	}

	// Find the appropriate asset for this platform
	assetName := u.getAssetName()
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			info.Available = true
			info.DownloadURL = asset.BrowserDownloadURL
			info.AssetName = asset.Name
			return info, nil
		}
	}

	return nil, fmt.Errorf("no compatible binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
}

// Update performs the update
func (u *Updater) Update(info *UpdateInfo) error {
	if !info.Available {
		return fmt.Errorf("no update available")
	}

	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Create backup of current binary
	backupPath := execPath + ".backup"
	if err := u.createBackup(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Download new binary to temp location
	tempPath := execPath + ".tmp"
	if err := u.downloadBinary(info.DownloadURL, tempPath); err != nil {
		os.Remove(tempPath)
		os.Remove(backupPath)
		return fmt.Errorf("failed to download update: %w", err)
	}

	// Make the new binary executable
	if err := os.Chmod(tempPath, 0755); err != nil {
		os.Remove(tempPath)
		os.Remove(backupPath)
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Replace the current binary
	if err := u.replaceBinary(execPath, tempPath); err != nil {
		// Attempt to restore from backup
		if restoreErr := u.restoreBackup(backupPath, execPath); restoreErr != nil {
			return fmt.Errorf("failed to replace binary and restore backup: update error: %w, restore error: %v", err, restoreErr)
		}
		return fmt.Errorf("failed to replace binary (backup restored): %w", err)
	}

	// Save backup info for potential manual rollback
	if err := u.saveBackupInfo(backupPath, info.CurrentVersion); err != nil {
		// Non-critical error, just log it
		fmt.Fprintf(os.Stderr, "Warning: failed to save backup info: %v\n", err)
	}

	return nil
}

// fetchLatestRelease gets the latest release information from GitHub
func (u *Updater) fetchLatestRelease() (*Release, error) {
	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return nil, err
	}

	// GitHub recommends setting this header
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// isNewerVersion compares version strings
func (u *Updater) isNewerVersion(latestVersion string) bool {
	current := strings.TrimPrefix(u.currentVersion, "v")
	latest := strings.TrimPrefix(latestVersion, "v")

	// Simple version comparison (works for semantic versioning)
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	for i := 0; i < len(currentParts) && i < len(latestParts); i++ {
		var currentNum, latestNum int
		fmt.Sscanf(currentParts[i], "%d", &currentNum)
		fmt.Sscanf(latestParts[i], "%d", &latestNum)

		if latestNum > currentNum {
			return true
		} else if latestNum < currentNum {
			return false
		}
	}

	return len(latestParts) > len(currentParts)
}

// getAssetName returns the expected asset name for the current platform
func (u *Updater) getAssetName() string {
	name := fmt.Sprintf("inkwash-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

// createBackup creates a backup of the current binary
func (u *Updater) createBackup(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Preserve permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// downloadBinary downloads the new binary
func (u *Updater) downloadBinary(url, dst string) error {
	// Use our existing download utility
	return utils.DownloadFile(url, dst, fmt.Sprintf("Downloading update"))
}

// replaceBinary replaces the current binary with the new one
func (u *Updater) replaceBinary(current, new string) error {
	// On Windows, we need to rename the current binary first
	if runtime.GOOS == "windows" {
		if err := os.Rename(current, current+".old"); err != nil {
			return err
		}
		if err := os.Rename(new, current); err != nil {
			// Try to restore
			os.Rename(current+".old", current)
			return err
		}
		// Clean up old binary
		os.Remove(current + ".old")
	} else {
		// On Unix systems, we can directly replace
		if err := os.Rename(new, current); err != nil {
			return err
		}
	}

	return nil
}

// restoreBackup attempts to restore from backup
func (u *Updater) restoreBackup(backup, dst string) error {
	return os.Rename(backup, dst)
}

// GetUpdateCheckPath returns the path to store last update check time
func GetUpdateCheckPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".inkwash", "last_update_check")
}

// ShouldCheckForUpdate determines if we should check for updates
func ShouldCheckForUpdate() bool {
	checkPath := GetUpdateCheckPath()
	
	info, err := os.Stat(checkPath)
	if err != nil {
		// File doesn't exist, should check
		return true
	}

	// Check if last check was more than 24 hours ago
	return time.Since(info.ModTime()) > 24*time.Hour
}

// SaveUpdateCheckTime saves the current time as the last update check
func SaveUpdateCheckTime() error {
	checkPath := GetUpdateCheckPath()
	dir := filepath.Dir(checkPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(checkPath, []byte(time.Now().Format(time.RFC3339)), 0644)
}

// saveBackupInfo saves information about the backup
func (u *Updater) saveBackupInfo(backupPath, version string) error {
	homeDir, _ := os.UserHomeDir()
	infoPath := filepath.Join(homeDir, ".inkwash", "backup_info.json")
	
	info := struct {
		BackupPath string    `json:"backup_path"`
		Version    string    `json:"version"`
		CreatedAt  time.Time `json:"created_at"`
	}{
		BackupPath: backupPath,
		Version:    version,
		CreatedAt:  time.Now(),
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(infoPath, data, 0644)
}

// Rollback attempts to rollback to the previous version
func Rollback() error {
	homeDir, _ := os.UserHomeDir()
	infoPath := filepath.Join(homeDir, ".inkwash", "backup_info.json")

	// Read backup info
	data, err := os.ReadFile(infoPath)
	if err != nil {
		return fmt.Errorf("no backup information found: %w", err)
	}

	var info struct {
		BackupPath string    `json:"backup_path"`
		Version    string    `json:"version"`
		CreatedAt  time.Time `json:"created_at"`
	}

	if err := json.Unmarshal(data, &info); err != nil {
		return fmt.Errorf("failed to parse backup info: %w", err)
	}

	// Check if backup exists
	if _, err := os.Stat(info.BackupPath); err != nil {
		return fmt.Errorf("backup file not found at %s: %w", info.BackupPath, err)
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Restore the backup
	updater := &Updater{}
	if err := updater.restoreBackup(info.BackupPath, execPath); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// Remove backup info file
	os.Remove(infoPath)

	fmt.Printf("Successfully rolled back to version %s\n", info.Version)
	return nil
}