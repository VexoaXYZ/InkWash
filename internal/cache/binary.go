package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/VexoaXYZ/inkwash/pkg/types"
)

// BinaryCache manages cached FXServer builds
type BinaryCache struct {
	basePath  string
	metadata  *Metadata
	maxBuilds int
}

// NewBinaryCache creates a new binary cache
func NewBinaryCache(basePath string, maxBuilds int) (*BinaryCache, error) {
	if maxBuilds <= 0 {
		maxBuilds = 3
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	bc := &BinaryCache{
		basePath:  basePath,
		maxBuilds: maxBuilds,
	}

	// Load or create metadata
	if err := bc.loadMetadata(); err != nil {
		return nil, err
	}

	return bc, nil
}

// Has checks if a build is cached
func (bc *BinaryCache) Has(buildNumber int) bool {
	for _, build := range bc.metadata.Builds {
		if build.Number == buildNumber {
			return true
		}
	}
	return false
}

// Get returns the path to a cached build's extracted files
func (bc *BinaryCache) Get(buildNumber int) (string, error) {
	buildPath := filepath.Join(bc.basePath, strconv.Itoa(buildNumber), "extracted")

	// Check if it exists
	if _, err := os.Stat(buildPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("build %d not in cache", buildNumber)
		}
		return "", err
	}

	// Update last used time
	bc.updateLastUsed(buildNumber)

	return buildPath, nil
}

// Add adds a build to the cache
func (bc *BinaryCache) Add(build types.Build, archivePath, extractedPath string) error {
	buildDir := filepath.Join(bc.basePath, strconv.Itoa(build.Number))

	// Create build directory
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	// Copy archive
	archiveName := filepath.Base(archivePath)
	destArchive := filepath.Join(buildDir, archiveName)

	if archivePath != destArchive {
		if err := copyFile(archivePath, destArchive); err != nil {
			return fmt.Errorf("failed to copy archive: %w", err)
		}
	}

	// Move/copy extracted files
	destExtracted := filepath.Join(buildDir, "extracted")
	if extractedPath != destExtracted {
		if err := os.Rename(extractedPath, destExtracted); err != nil {
			// If rename fails (cross-device), try copy
			if err := copyDir(extractedPath, destExtracted); err != nil {
				return fmt.Errorf("failed to move extracted files: %w", err)
			}
			os.RemoveAll(extractedPath)
		}
	}

	// Get archive size
	archiveInfo, err := os.Stat(destArchive)
	if err != nil {
		return fmt.Errorf("failed to stat archive: %w", err)
	}

	// Add to metadata
	cacheBuild := CachedBuild{
		Number:      build.Number,
		Hash:        build.Hash,
		Downloaded:  time.Now(),
		Size:        archiveInfo.Size(),
		Recommended: build.Recommended,
		Optional:    build.Optional,
		LastUsed:    time.Now(),
	}

	bc.metadata.Builds = append(bc.metadata.Builds, cacheBuild)
	bc.metadata.TotalSize += archiveInfo.Size()

	// Enforce cache limits
	if err := bc.enforceLimits(); err != nil {
		return err
	}

	// Save metadata
	return bc.saveMetadata()
}

// Remove removes a build from the cache
func (bc *BinaryCache) Remove(buildNumber int) error {
	buildDir := filepath.Join(bc.basePath, strconv.Itoa(buildNumber))

	// Get build size for metadata update
	var buildSize int64
	for i, build := range bc.metadata.Builds {
		if build.Number == buildNumber {
			buildSize = build.Size
			// Remove from metadata
			bc.metadata.Builds = append(bc.metadata.Builds[:i], bc.metadata.Builds[i+1:]...)
			break
		}
	}

	// Remove directory
	if err := os.RemoveAll(buildDir); err != nil {
		return fmt.Errorf("failed to remove build directory: %w", err)
	}

	bc.metadata.TotalSize -= buildSize

	return bc.saveMetadata()
}

// List returns all cached builds
func (bc *BinaryCache) List() []CachedBuild {
	return bc.metadata.Builds
}

// Clear removes all cached builds
func (bc *BinaryCache) Clear() error {
	for _, build := range bc.metadata.Builds {
		buildDir := filepath.Join(bc.basePath, strconv.Itoa(build.Number))
		if err := os.RemoveAll(buildDir); err != nil {
			return fmt.Errorf("failed to remove build %d: %w", build.Number, err)
		}
	}

	bc.metadata.Builds = []CachedBuild{}
	bc.metadata.TotalSize = 0

	return bc.saveMetadata()
}

// GetStats returns cache statistics
func (bc *BinaryCache) GetStats() CacheStats {
	return CacheStats{
		TotalBuilds: len(bc.metadata.Builds),
		TotalSize:   bc.metadata.TotalSize,
		MaxBuilds:   bc.maxBuilds,
	}
}

// enforceLimits enforces cache size limits using LRU eviction
func (bc *BinaryCache) enforceLimits() error {
	if len(bc.metadata.Builds) <= bc.maxBuilds {
		return nil
	}

	// Sort by last used (oldest first)
	sort.Slice(bc.metadata.Builds, func(i, j int) bool {
		return bc.metadata.Builds[i].LastUsed.Before(bc.metadata.Builds[j].LastUsed)
	})

	// Remove oldest builds
	toRemove := len(bc.metadata.Builds) - bc.maxBuilds
	for i := 0; i < toRemove; i++ {
		build := bc.metadata.Builds[0]
		if err := bc.Remove(build.Number); err != nil {
			return err
		}
	}

	return nil
}

// updateLastUsed updates the last used timestamp for a build
func (bc *BinaryCache) updateLastUsed(buildNumber int) {
	for i, build := range bc.metadata.Builds {
		if build.Number == buildNumber {
			bc.metadata.Builds[i].LastUsed = time.Now()
			bc.saveMetadata()
			return
		}
	}
}

// loadMetadata loads metadata from disk
func (bc *BinaryCache) loadMetadata() error {
	metadataPath := filepath.Join(bc.basePath, "metadata.json")

	// If metadata doesn't exist, create empty
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		bc.metadata = &Metadata{
			Version:   1,
			Builds:    []CachedBuild{},
			MaxBuilds: bc.maxBuilds,
			TotalSize: 0,
		}
		return bc.saveMetadata()
	}

	// Load existing metadata
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata Metadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	bc.metadata = &metadata
	return nil
}

// saveMetadata saves metadata to disk
func (bc *BinaryCache) saveMetadata() error {
	metadataPath := filepath.Join(bc.basePath, "metadata.json")

	data, err := json.MarshalIndent(bc.metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// Helper functions

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

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
