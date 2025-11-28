package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Progress represents download progress
type Progress struct {
	TotalBytes      int64
	DownloadedBytes int64
	Speed           float64 // MB/s
	ETA             time.Duration
	ChunkProgress   []int64 // Bytes downloaded per chunk
}

// ProgressCallback is called periodically with download progress
type ProgressCallback func(Progress)

// Downloader handles parallel downloads
type Downloader struct {
	httpClient *http.Client
	numChunks  int
}

// NewDownloader creates a new downloader
func NewDownloader(numChunks int) *Downloader {
	if numChunks <= 0 {
		numChunks = 3
	}

	return &Downloader{
		httpClient: &http.Client{
			Timeout: 10 * time.Minute,
		},
		numChunks: numChunks,
	}
}

// Download downloads a file with parallel chunks
func (d *Downloader) Download(url, destPath string, onProgress ProgressCallback) error {
	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Get file size
	totalSize, err := d.getFileSize(url)
	if err != nil {
		return err
	}

	// Check if server supports range requests
	supportsRanges, err := d.supportsRangeRequests(url)
	if err != nil {
		return err
	}

	if !supportsRanges {
		// Fallback to single download
		return d.downloadSingle(url, destPath, totalSize, onProgress)
	}

	// Download in parallel chunks
	return d.downloadParallel(url, destPath, totalSize, onProgress)
}

// downloadParallel downloads a file in parallel chunks
func (d *Downloader) downloadParallel(url, destPath string, totalSize int64, onProgress ProgressCallback) error {
	chunkSize := totalSize / int64(d.numChunks)

	// Create progress tracker
	progress := Progress{
		TotalBytes:    totalSize,
		ChunkProgress: make([]int64, d.numChunks),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, d.numChunks)
	progressChan := make(chan struct{}, 100)

	// Start progress reporter
	stopProgress := make(chan struct{})
	go d.reportProgress(&progress, &mu, onProgress, progressChan, stopProgress)

	// Download chunks
	for i := 0; i < d.numChunks; i++ {
		wg.Add(1)
		go func(chunkID int) {
			defer wg.Done()

			start := int64(chunkID) * chunkSize
			end := start + chunkSize - 1

			// Last chunk gets any remainder
			if chunkID == d.numChunks-1 {
				end = totalSize - 1
			}

			chunkPath := fmt.Sprintf("%s.part%d", destPath, chunkID)

			if err := d.downloadChunk(url, start, end, chunkPath, chunkID, &progress, &mu, progressChan); err != nil {
				errChan <- fmt.Errorf("chunk %d failed: %w", chunkID, err)
			}
		}(i)
	}

	wg.Wait()
	close(stopProgress)
	close(errChan)

	// Check for errors
	if len(errChan) > 0 {
		return <-errChan
	}

	// Merge chunks
	return d.mergeChunks(destPath, d.numChunks)
}

// downloadChunk downloads a single chunk
func (d *Downloader) downloadChunk(url string, start, end int64, destPath string, chunkID int, progress *Progress, mu *sync.Mutex, progressChan chan struct{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set range header
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create chunk file
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Download with progress tracking
	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := file.Write(buffer[:n]); writeErr != nil {
				return writeErr
			}

			// Update progress
			mu.Lock()
			progress.ChunkProgress[chunkID] += int64(n)
			mu.Unlock()

			// Notify progress reporter (non-blocking)
			select {
			case progressChan <- struct{}{}:
			default:
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// reportProgress reports download progress periodically
func (d *Downloader) reportProgress(progress *Progress, mu *sync.Mutex, callback ProgressCallback, progressChan chan struct{}, stop chan struct{}) {
	if callback == nil {
		return
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	startTime := time.Now()
	lastBytes := int64(0)
	lastTime := startTime

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			mu.Lock()

			// Calculate total downloaded
			total := int64(0)
			for _, bytes := range progress.ChunkProgress {
				total += bytes
			}
			progress.DownloadedBytes = total

			// Calculate speed (MB/s)
			elapsed := time.Since(lastTime).Seconds()
			if elapsed > 0 {
				deltaBytes := float64(total - lastBytes)
				progress.Speed = (deltaBytes / elapsed) / 1024 / 1024
			}

			// Calculate ETA
			if progress.Speed > 0 {
				remaining := float64(progress.TotalBytes - progress.DownloadedBytes)
				etaSeconds := remaining / (progress.Speed * 1024 * 1024)
				progress.ETA = time.Duration(etaSeconds) * time.Second
			}

			lastBytes = total
			lastTime = time.Now()

			// Create copy for callback
			progressCopy := *progress

			mu.Unlock()

			callback(progressCopy)
		}
	}
}

// mergeChunks merges chunk files into the final file
func (d *Downloader) mergeChunks(destPath string, numChunks int) error {
	// Create final file
	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Merge chunks in order
	for i := 0; i < numChunks; i++ {
		chunkPath := fmt.Sprintf("%s.part%d", destPath, i)

		// Open chunk file
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("failed to open chunk %d: %w", i, err)
		}

		// Copy chunk to output
		if _, err := io.Copy(outFile, chunkFile); err != nil {
			chunkFile.Close()
			return fmt.Errorf("failed to copy chunk %d: %w", i, err)
		}

		chunkFile.Close()

		// Delete chunk file
		os.Remove(chunkPath)
	}

	return nil
}

// downloadSingle downloads a file without chunking
func (d *Downloader) downloadSingle(url, destPath string, totalSize int64, onProgress ProgressCallback) error {
	resp, err := d.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Download with progress tracking
	progress := Progress{
		TotalBytes:    totalSize,
		ChunkProgress: []int64{0},
	}

	buffer := make([]byte, 32*1024)
	startTime := time.Now()
	lastUpdate := startTime

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := file.Write(buffer[:n]); writeErr != nil {
				return writeErr
			}

			progress.ChunkProgress[0] += int64(n)
			progress.DownloadedBytes = progress.ChunkProgress[0]

			// Report progress every 100ms
			if time.Since(lastUpdate) >= 100*time.Millisecond {
				elapsed := time.Since(startTime).Seconds()
				if elapsed > 0 {
					progress.Speed = float64(progress.DownloadedBytes) / elapsed / 1024 / 1024
				}

				if progress.Speed > 0 {
					remaining := float64(progress.TotalBytes - progress.DownloadedBytes)
					etaSeconds := remaining / (progress.Speed * 1024 * 1024)
					progress.ETA = time.Duration(etaSeconds) * time.Second
				}

				if onProgress != nil {
					onProgress(progress)
				}

				lastUpdate = time.Now()
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// getFileSize gets the file size from a URL
func (d *Downloader) getFileSize(url string) (int64, error) {
	resp, err := d.httpClient.Head(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("no content-length header")
	}

	var size int64
	if _, err := fmt.Sscanf(contentLength, "%d", &size); err != nil {
		return 0, fmt.Errorf("failed to parse content-length: %w", err)
	}

	return size, nil
}

// supportsRangeRequests checks if the server supports range requests
func (d *Downloader) supportsRangeRequests(url string) (bool, error) {
	resp, err := d.httpClient.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	acceptRanges := resp.Header.Get("Accept-Ranges")
	return acceptRanges == "bytes", nil
}
