package services

import (
	"io"
	"os"
	"path/filepath"

	"github.com/vexoa/inkwash/internal/domain"
)

// fileServiceImpl implements FileService
type fileServiceImpl struct{}

// NewFileService creates a new file service
func NewFileService() FileService {
	return &fileServiceImpl{}
}

// ReadFile reads a file
func (s *fileServiceImpl) ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, domain.ErrFilesystemOperation("read_file", path, err)
	}
	return data, nil
}

// WriteFile writes data to a file
func (s *fileServiceImpl) WriteFile(path string, data []byte, perm os.FileMode) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := s.CreateDirectory(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(path, data, perm); err != nil {
		return domain.ErrFilesystemOperation("write_file", path, err)
	}
	return nil
}

// CopyFile copies a file
func (s *fileServiceImpl) CopyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return domain.ErrFilesystemOperation("open_source", src, err)
	}
	defer srcFile.Close()

	// Create destination directory
	dir := filepath.Dir(dst)
	if err := s.CreateDirectory(dir, 0755); err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return domain.ErrFilesystemOperation("create_destination", dst, err)
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return domain.ErrFilesystemOperation("copy_contents", dst, err)
	}

	// Copy permissions
	srcInfo, err := srcFile.Stat()
	if err == nil {
		os.Chmod(dst, srcInfo.Mode())
	}

	return nil
}

// MoveFile moves a file
func (s *fileServiceImpl) MoveFile(src, dst string) error {
	// Create destination directory
	dir := filepath.Dir(dst)
	if err := s.CreateDirectory(dir, 0755); err != nil {
		return err
	}

	if err := os.Rename(src, dst); err != nil {
		// If rename fails, try copy and delete
		if copyErr := s.CopyFile(src, dst); copyErr != nil {
			return domain.ErrFilesystemOperation("move_file", src, err)
		}
		if delErr := s.DeleteFile(src); delErr != nil {
			return domain.ErrFilesystemOperation("delete_source", src, delErr)
		}
	}
	return nil
}

// DeleteFile deletes a file
func (s *fileServiceImpl) DeleteFile(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return domain.ErrFilesystemOperation("delete_file", path, err)
	}
	return nil
}

// CreateDirectory creates a directory
func (s *fileServiceImpl) CreateDirectory(path string, perm os.FileMode) error {
	if err := os.MkdirAll(path, perm); err != nil {
		return domain.ErrFilesystemOperation("create_directory", path, err)
	}
	return nil
}

// ListDirectory lists directory contents
func (s *fileServiceImpl) ListDirectory(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, domain.ErrFilesystemOperation("list_directory", path, err)
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	return files, nil
}

// FileExists checks if a file exists
func (s *fileServiceImpl) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetFileInfo gets file information
func (s *fileServiceImpl) GetFileInfo(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, domain.ErrFilesystemOperation("get_file_info", path, err)
	}
	return info, nil
}