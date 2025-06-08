package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/vexoa/inkwash/internal/domain"
)

// downloadServiceImpl implements DownloadService
type downloadServiceImpl struct {
	client *http.Client
}

// NewDownloadService creates a new download service
func NewDownloadService() DownloadService {
	return &downloadServiceImpl{
		client: &http.Client{
			Timeout: 0, // No timeout for downloads
		},
	}
}

// Download downloads a file from URL
func (s *downloadServiceImpl) Download(ctx context.Context, url, destPath string, progress ProgressCallback) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return domain.ErrFilesystemOperation("create_directory", dir, err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.ErrDownloadFailed(url, err)
	}

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.ErrDownloadFailed(url, err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return domain.ErrDownloadFailed(url, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status))
	}

	// Create the output file
	out, err := os.Create(destPath)
	if err != nil {
		return domain.ErrFilesystemOperation("create_file", destPath, err)
	}
	defer out.Close()

	// Get content length for progress reporting
	contentLength := resp.ContentLength
	
	// Copy with progress reporting
	if progress != nil && contentLength > 0 {
		return s.copyWithProgress(resp.Body, out, contentLength, progress)
	}

	// Simple copy without progress
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return domain.ErrDownloadFailed(url, err)
	}

	return nil
}

// DownloadWithResume downloads with resume support
func (s *downloadServiceImpl) DownloadWithResume(ctx context.Context, url, destPath string, progress ProgressCallback) error {
	// Check if file already exists
	var startPos int64 = 0
	if info, err := os.Stat(destPath); err == nil {
		startPos = info.Size()
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return domain.ErrFilesystemOperation("create_directory", dir, err)
	}

	// Create the request with Range header for resume
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.ErrDownloadFailed(url, err)
	}

	if startPos > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startPos))
	}

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.ErrDownloadFailed(url, err)
	}
	defer resp.Body.Close()

	// Check status code (206 for partial content, 200 for full)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return domain.ErrDownloadFailed(url, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status))
	}

	// Open file for appending if resuming, create new if starting fresh
	var out *os.File
	if startPos > 0 && resp.StatusCode == http.StatusPartialContent {
		out, err = os.OpenFile(destPath, os.O_WRONLY|os.O_APPEND, 0644)
	} else {
		out, err = os.Create(destPath)
		startPos = 0 // Reset if server doesn't support resume
	}

	if err != nil {
		return domain.ErrFilesystemOperation("open_file", destPath, err)
	}
	defer out.Close()

	// Get total content length
	totalLength := resp.ContentLength + startPos

	// Copy with progress reporting
	if progress != nil && totalLength > 0 {
		return s.copyWithProgressAndOffset(resp.Body, out, resp.ContentLength, totalLength, startPos, progress)
	}

	// Simple copy without progress
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return domain.ErrDownloadFailed(url, err)
	}

	return nil
}

// GetContentLength gets the content length of a URL
func (s *downloadServiceImpl) GetContentLength(ctx context.Context, url string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return 0, domain.ErrDownloadFailed(url, err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, domain.ErrDownloadFailed(url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, domain.ErrDownloadFailed(url, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status))
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, nil // Unknown length
	}

	length, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, domain.ErrDownloadFailed(url, err)
	}

	return length, nil
}

// copyWithProgress copies data with progress reporting
func (s *downloadServiceImpl) copyWithProgress(src io.Reader, dst io.Writer, total int64, progress ProgressCallback) error {
	var written int64
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				if progress != nil {
					progress(written, total, "Downloading...")
				}
			}
			if ew != nil {
				return ew
			}
			if nr != nw {
				return io.ErrShortWrite
			}
		}
		if er != nil {
			if er == io.EOF {
				break
			}
			return er
		}
	}

	return nil
}

// copyWithProgressAndOffset copies data with progress reporting and offset
func (s *downloadServiceImpl) copyWithProgressAndOffset(src io.Reader, dst io.Writer, currentSize, total, offset int64, progress ProgressCallback) error {
	var written int64 = offset
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				if progress != nil {
					progress(written, total, "Downloading...")
				}
			}
			if ew != nil {
				return ew
			}
			if nr != nw {
				return io.ErrShortWrite
			}
		}
		if er != nil {
			if er == io.EOF {
				break
			}
			return er
		}
	}

	return nil
}