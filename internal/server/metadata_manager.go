package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/VexoaXYZ/inkwash/pkg/types"
)

const metadataFilename = "metadata.json"

// MetadataManager handles reading/writing server metadata
type MetadataManager struct{}

// NewMetadataManager creates a new metadata manager
func NewMetadataManager() *MetadataManager {
	return &MetadataManager{}
}

// GetMetadataPath returns the path to a server's metadata.json
func (mm *MetadataManager) GetMetadataPath(serverPath string) string {
	return filepath.Join(serverPath, metadataFilename)
}

// Load loads metadata from a server's metadata.json
func (mm *MetadataManager) Load(serverPath string) (*types.ServerMetadata, error) {
	metadataPath := mm.GetMetadataPath(serverPath)

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("metadata.json not found at %s", metadataPath)
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata types.ServerMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &metadata, nil
}

// Save writes metadata to a server's metadata.json
func (mm *MetadataManager) Save(serverPath string, metadata *types.ServerMetadata) error {
	metadataPath := mm.GetMetadataPath(serverPath)

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// Exists checks if metadata.json exists
func (mm *MetadataManager) Exists(serverPath string) bool {
	_, err := os.Stat(mm.GetMetadataPath(serverPath))
	return err == nil
}

// RecordStart updates metadata when server starts
func (mm *MetadataManager) RecordStart(serverPath string) error {
	metadata, err := mm.Load(serverPath)
	if err != nil {
		return err
	}

	now := time.Now()
	metadata.Lifecycle.LastStarted = &now
	metadata.Stats.RestartCount++

	return mm.Save(serverPath, metadata)
}

// RecordStop updates metadata when server stops
func (mm *MetadataManager) RecordStop(serverPath string, startTime time.Time) error {
	metadata, err := mm.Load(serverPath)
	if err != nil {
		return err
	}

	now := time.Now()
	metadata.Lifecycle.LastStopped = &now
	metadata.Stats.TotalUptime += now.Sub(startTime)

	return mm.Save(serverPath, metadata)
}
