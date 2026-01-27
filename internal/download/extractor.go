package download

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/ulikunitz/xz"
)

// Extractor handles archive extraction
type Extractor struct{}

// NewExtractor creates a new extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Extract extracts an archive to the destination directory
func (e *Extractor) Extract(archivePath, destPath string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Determine archive type from extension
	if strings.HasSuffix(archivePath, ".7z") {
		return e.extract7z(archivePath, destPath)
	} else if strings.HasSuffix(archivePath, ".tar.xz") {
		return e.extractTarXz(archivePath, destPath)
	} else if strings.HasSuffix(archivePath, ".tar.gz") {
		return e.extractTarGz(archivePath, destPath)
	} else if strings.HasSuffix(archivePath, ".zip") {
		return e.extractZip(archivePath, destPath)
	}

	return fmt.Errorf("unsupported archive format: %s", archivePath)
}

// extract7z extracts a 7z archive (Windows)
func (e *Extractor) extract7z(src, dest string) error {
	r, err := sevenzip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open 7z archive: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)

		// Security check: prevent path traversal
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path, err)
			}
			continue
		}

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Extract file
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in archive: %w", err)
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create output file %s: %w", path, err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", path, err)
		}
	}

	return nil
}

// extractTarXz extracts a tar.xz archive (Linux)
func (e *Extractor) extractTarXz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer f.Close()

	// Create XZ reader
	xzReader, err := xz.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to create xz reader: %w", err)
	}

	// Create tar reader
	tarReader := tar.NewReader(xzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		path := filepath.Join(dest, header.Name)

		// Security check: prevent path traversal
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)) {
			return fmt.Errorf("illegal file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path, err)
			}

		case tar.TypeReg:
			// Create parent directory
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create output file %s: %w", path, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to extract file %s: %w", path, err)
			}

			outFile.Close()

		case tar.TypeSymlink:
			// Handle symlinks (important for Linux)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Remove existing file/symlink if it exists
			os.Remove(path)

			if err := os.Symlink(header.Linkname, path); err != nil {
				return fmt.Errorf("failed to create symlink %s: %w", path, err)
			}
		}
	}

	return nil
}

// extractTarGz extracts a tar.gz archive (fallback/utility)
func (e *Extractor) extractTarGz(src, dest string) error {
	// Similar to extractTarXz but with gzip instead of xz
	// Not needed for FiveM but useful for future
	return fmt.Errorf("tar.gz extraction not implemented yet")
}

// extractZip extracts a zip archive
func (e *Extractor) extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)

		// Security check: prevent path traversal
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path, err)
			}
			continue
		}

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Extract file
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in archive: %w", err)
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create output file %s: %w", path, err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", path, err)
		}
	}

	return nil
}

// ExtractWithProgress extracts an archive with progress callback
func (e *Extractor) ExtractWithProgress(archivePath, destPath string, onProgress func(current, total int)) error {
	// For now, just extract without progress
	// TODO: Implement progress tracking by counting files
	return e.Extract(archivePath, destPath)
}

// GetArchiveFileCount returns the number of files in an archive
func (e *Extractor) GetArchiveFileCount(archivePath string) (int, error) {
	if strings.HasSuffix(archivePath, ".7z") {
		r, err := sevenzip.OpenReader(archivePath)
		if err != nil {
			return 0, err
		}
		defer r.Close()
		return len(r.File), nil
	}

	if strings.HasSuffix(archivePath, ".tar.xz") {
		f, err := os.Open(archivePath)
		if err != nil {
			return 0, err
		}
		defer f.Close()

		xzReader, err := xz.NewReader(f)
		if err != nil {
			return 0, err
		}

		tarReader := tar.NewReader(xzReader)
		count := 0

		for {
			_, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return 0, err
			}
			count++
		}

		return count, nil
	}

	return 0, fmt.Errorf("unsupported archive format")
}

// GetPlatformArchiveExtension returns the archive extension for the current platform
func GetPlatformArchiveExtension() string {
	if runtime.GOOS == "windows" {
		return ".7z"
	}
	return ".tar.xz"
}
