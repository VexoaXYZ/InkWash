package download

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/VexoaXYZ/inkwash/pkg/types"
)

const (
	WindowsArtifactURL = "https://runtime.fivem.net/artifacts/fivem/build_server_windows/master/"
	LinuxArtifactURL   = "https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/"
)

// ArtifactClient handles fetching FiveM server builds
type ArtifactClient struct {
	httpClient *http.Client
}

// NewArtifactClient creates a new artifact client
func NewArtifactClient() *ArtifactClient {
	return &ArtifactClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchBuilds fetches available builds from the FiveM artifacts page
func (ac *ArtifactClient) FetchBuilds() ([]types.Build, error) {
	url := ac.getArtifactURL()

	resp, err := ac.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artifacts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML directory listing
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return ac.parseBuilds(doc)
}

// getArtifactURL returns the appropriate artifact URL for the current platform
func (ac *ArtifactClient) getArtifactURL() string {
	if runtime.GOOS == "windows" {
		return WindowsArtifactURL
	}
	return LinuxArtifactURL
}

// parseBuilds parses builds from the HTML document
func (ac *ArtifactClient) parseBuilds(doc *goquery.Document) ([]types.Build, error) {
	var builds []types.Build
	pageText := doc.Text()

	// Find recommended and optional build numbers from page text
	recommendedBuild := ac.findRecommendedBuild(pageText)
	optionalBuild := ac.findOptionalBuild(pageText)

	// Parse build entries from links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")

		// Look for build archive links: ./BUILD-HASH/server.7z or ./BUILD-HASH/fx.tar.xz
		if !strings.Contains(href, "/server.7z") && !strings.Contains(href, "/fx.tar.xz") {
			return
		}

		// Extract directory part: "./22934-1c490ee35560b652c97a4bfd5a5852cb9f033284/server.7z"
		// Remove "./" prefix and "/server.7z" or "/fx.tar.xz" suffix
		href = strings.TrimPrefix(href, "./")
		href = strings.TrimSuffix(href, "/server.7z")
		href = strings.TrimSuffix(href, "/fx.tar.xz")

		// Parse: "22934-1c490ee35560b652c97a4bfd5a5852cb9f033284"
		parts := strings.SplitN(href, "-", 2)
		if len(parts) < 2 {
			return
		}

		number, err := strconv.Atoi(parts[0])
		if err != nil {
			return
		}

		hash := parts[1]

		build := types.Build{
			Number:      number,
			Hash:        fmt.Sprintf("%d-%s", number, hash),
			Timestamp:   time.Now(), // We don't have exact timestamp from the page
			Recommended: number == recommendedBuild,
			Optional:    number == optionalBuild,
		}

		builds = append(builds, build)
	})

	if len(builds) == 0 {
		return nil, fmt.Errorf("no builds found")
	}

	return builds, nil
}

// findRecommendedBuild extracts the recommended build number from page text
func (ac *ArtifactClient) findRecommendedBuild(pageText string) int {
	// Look for pattern like "LATEST RECOMMENDED (17000)"
	start := strings.Index(pageText, "LATEST RECOMMENDED")
	if start == -1 {
		return 0
	}

	// Find opening parenthesis
	openParen := strings.Index(pageText[start:], "(")
	if openParen == -1 {
		return 0
	}

	// Find closing parenthesis
	closeParen := strings.Index(pageText[start+openParen:], ")")
	if closeParen == -1 {
		return 0
	}

	// Extract number
	numberStr := pageText[start+openParen+1 : start+openParen+closeParen]
	number, _ := strconv.Atoi(strings.TrimSpace(numberStr))

	return number
}

// findOptionalBuild extracts the optional build number from page text
func (ac *ArtifactClient) findOptionalBuild(pageText string) int {
	// Look for pattern like "LATEST OPTIONAL (7290)"
	start := strings.Index(pageText, "LATEST OPTIONAL")
	if start == -1 {
		return 0
	}

	// Find opening parenthesis
	openParen := strings.Index(pageText[start:], "(")
	if openParen == -1 {
		return 0
	}

	// Find closing parenthesis
	closeParen := strings.Index(pageText[start+openParen:], ")")
	if closeParen == -1 {
		return 0
	}

	// Extract number
	numberStr := pageText[start+openParen+1 : start+openParen+closeParen]
	number, _ := strconv.Atoi(strings.TrimSpace(numberStr))

	return number
}

// GetDownloadURL returns the download URL for a specific build
func (ac *ArtifactClient) GetDownloadURL(build types.Build) string {
	var baseURL string
	var filename string

	if runtime.GOOS == "windows" {
		baseURL = WindowsArtifactURL
		filename = "server.7z"
	} else {
		baseURL = LinuxArtifactURL
		filename = "fx.tar.xz"
	}

	return fmt.Sprintf("%s%s/%s", baseURL, build.Hash, filename)
}

// GetFileSize gets the size of a file from a URL using HEAD request
func (ac *ArtifactClient) GetFileSize(url string) (int64, error) {
	resp, err := ac.httpClient.Head(url)
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

	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse content-length: %w", err)
	}

	return size, nil
}
