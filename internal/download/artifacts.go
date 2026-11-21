package download

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/VexoaXYZ/inkwash/pkg/types"
)

const (
	// ServerDownloadPageURL lists every published server artifact. The old
	// runtime.fivem.net directory listings now redirect here.
	ServerDownloadPageURL = "https://docs.fivem.net/docs/server-download/"

	WindowsArtifactURL = "https://runtime.fivem.net/artifacts/fivem/build_server_windows/master/"
	LinuxArtifactURL   = "https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/"
)

var (
	// buildLinkPattern matches an artifact archive link in any form: relative
	// ("./38435-<sha>/server.7z") or absolute ("https://.../38435-<sha>/server.7z").
	buildLinkPattern = regexp.MustCompile(`(\d+)-([0-9a-f]{40})/(server\.7z|fx\.tar\.xz)`)

	nextDataPattern = regexp.MustCompile(`(?s)<script id="__NEXT_DATA__"[^>]*>(.*?)</script>`)

	recommendedTextPattern = regexp.MustCompile(`LATEST RECOMMENDED\s*\(\s*(\d+)\s*\)`)
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

// FetchBuilds fetches available builds from the FiveM server download page
func (ac *ArtifactClient) FetchBuilds() ([]types.Build, error) {
	resp, err := ac.httpClient.Get(ServerDownloadPageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artifacts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read artifacts page: %w", err)
	}

	return parseBuilds(string(body), platformKey(), archiveName())
}

// platformKey returns the key the download page uses for this OS
func platformKey() string {
	if runtime.GOOS == "windows" {
		return "windows"
	}
	return "linux"
}

// archiveName returns the artifact archive file name for this OS
func archiveName() string {
	if runtime.GOOS == "windows" {
		return "server.7z"
	}
	return "fx.tar.xz"
}

// artifactEntry is one download entry in the page's embedded data
type artifactEntry struct {
	DownloadURL string `json:"downloadURL"`
}

// downloadPageData is the subset of the download page's __NEXT_DATA__ we use
type downloadPageData struct {
	Props struct {
		PageProps struct {
			Legacy struct {
				Recommended map[string][]artifactEntry `json:"recommended"`
				Latest      map[string][]artifactEntry `json:"latest"`
			} `json:"legacy"`
		} `json:"pageProps"`
	} `json:"props"`
}

// parseBuilds extracts builds for one platform from the download page HTML.
// It reads the page's embedded build data first, and falls back to scanning
// every artifact link in the page if that data is missing or changes shape.
// The recommended build comes first, the rest follow newest first.
func parseBuilds(html, platform, archive string) ([]types.Build, error) {
	recommended := 0
	var numbers []int
	hashes := map[int]string{}

	add := func(link string) int {
		m := buildLinkPattern.FindStringSubmatch(link)
		if m == nil || m[3] != archive {
			return 0
		}
		number, err := strconv.Atoi(m[1])
		if err != nil {
			return 0
		}
		if _, seen := hashes[number]; !seen {
			hashes[number] = fmt.Sprintf("%d-%s", number, m[2])
			numbers = append(numbers, number)
		}
		return number
	}

	if m := nextDataPattern.FindStringSubmatch(html); m != nil {
		var data downloadPageData
		if err := json.Unmarshal([]byte(m[1]), &data); err == nil {
			legacy := data.Props.PageProps.Legacy
			for _, entry := range legacy.Recommended[platform] {
				if n := add(entry.DownloadURL); n != 0 && recommended == 0 {
					recommended = n
				}
			}
			for _, entry := range legacy.Latest[platform] {
				add(entry.DownloadURL)
			}
		}
	}

	if len(numbers) == 0 {
		for _, m := range buildLinkPattern.FindAllString(html, -1) {
			add(m)
		}
		if m := recommendedTextPattern.FindStringSubmatch(html); m != nil {
			recommended, _ = strconv.Atoi(m[1])
		}
	}

	if len(numbers) == 0 {
		return nil, fmt.Errorf("no builds found")
	}

	sort.Slice(numbers, func(i, j int) bool {
		if (numbers[i] == recommended) != (numbers[j] == recommended) {
			return numbers[i] == recommended
		}
		return numbers[i] > numbers[j]
	})

	builds := make([]types.Build, 0, len(numbers))
	for _, number := range numbers {
		builds = append(builds, types.Build{
			Number:      number,
			Hash:        hashes[number],
			Timestamp:   time.Now(), // the page does not publish build timestamps
			Recommended: number == recommended,
		})
	}

	return builds, nil
}

// GetDownloadURL returns the download URL for a specific build
func (ac *ArtifactClient) GetDownloadURL(build types.Build) string {
	baseURL := LinuxArtifactURL
	if runtime.GOOS == "windows" {
		baseURL = WindowsArtifactURL
	}

	return fmt.Sprintf("%s%s/%s", baseURL, build.Hash, archiveName())
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
