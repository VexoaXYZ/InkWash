package convert

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ConversionStatus represents the status of a mod conversion
type ConversionStatus struct {
	Progress int    `json:"progress"`
	Status   int    `json:"status"`
	File     string `json:"file"`
	Message  string `json:"message"`
	Name     string `json:"name"`
}

// ConvertResponse represents the initial conversion response
type ConvertResponse struct {
	Message string `json:"message"` // UUID
	Status  int    `json:"status"`
}

// Client handles GTA5 mod conversion via convert.cfx.rs
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new conversion client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://convert.cfx.rs",
	}
}

// StartConversion initiates a mod conversion
func (c *Client) StartConversion(modURL string) (string, error) {
	// Validate URL is from gta5-mods.com
	if !strings.Contains(modURL, "gta5-mods.com") {
		return "", fmt.Errorf("URL must be from gta5-mods.com")
	}

	// Prepare form data
	data := url.Values{}
	data.Set("url", modURL)
	data.Set("lang", "en")

	// Make POST request
	resp, err := c.httpClient.PostForm(c.baseURL+"/api/convert", data)
	if err != nil {
		return "", fmt.Errorf("failed to start conversion: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var result ConvertResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Status != 200 {
		return "", fmt.Errorf("conversion failed with status %d", result.Status)
	}

	return result.Message, nil
}

// QueryProgress checks the progress of a conversion
func (c *Client) QueryProgress(uuid string) (*ConversionStatus, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("uuid", uuid)
	data.Set("lang", "en")

	// Make POST request
	resp, err := c.httpClient.PostForm(c.baseURL+"/api/query", data)
	if err != nil {
		return nil, fmt.Errorf("failed to query progress: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var status ConversionStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &status, nil
}

// GetDownloadURL returns the download URL for a converted file
func (c *Client) GetDownloadURL(file string) string {
	return c.baseURL + "/" + file
}

// DownloadFile downloads a converted file to the specified path
func (c *Client) DownloadFile(fileURL, destPath string) error {
	resp, err := c.httpClient.Get(fileURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create destination file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy content
	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
