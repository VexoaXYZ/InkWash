package wizard

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/VexoaXYZ/inkwash/internal/convert"
	"github.com/VexoaXYZ/inkwash/internal/download"
	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/ui"
	"github.com/VexoaXYZ/inkwash/internal/ui/components"
	"github.com/VexoaXYZ/inkwash/pkg/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConvertStep represents the current step in the wizard
type ConvertStep int

const (
	ConvertStepSelectServer ConvertStep = iota
	ConvertStepCustomPath
	ConvertStepEnterURLs
	ConvertStepConverting
	ConvertStepDownloading
	ConvertStepComplete
	ConvertStepError
)

// ConversionItem tracks a single mod conversion
type ConversionItem struct {
	URL      string
	UUID     string
	Status   *convert.ConversionStatus
	Error    error
	FileName string
	Category string // e.g., "vehicles", "weapons", "scripts"
}

// ConvertWizardModel represents the state of the conversion wizard
type ConvertWizardModel struct {
	step      ConvertStep
	client    *convert.Client
	downloader *download.Downloader
	registry  *registry.Registry

	// Input components
	serverSelector *components.Selector
	urlInput       *components.TextInput
	customPathInput *components.TextInput

	// Progress components
	progressBar *components.ProgressBar
	spinner     *components.Spinner

	// State
	selectedServer *types.Server
	externalMode   string // "current" or "custom" or "" if using registered server
	customPath     string
	urls           []string
	conversions    map[string]*ConversionItem // UUID -> item
	conversionList []string                   // Ordered UUIDs
	downloads      []string                   // Files to download
	error          string
	quitting       bool
	completed      bool

	// Progress tracking
	overallProgress float64
	downloadProgress map[string]float64
	pollingActive   bool
	lastUpdate      time.Time

	// Queue management
	conversionQueue []string // URLs waiting to be converted
	activeConversions int    // Number of conversions in progress
	maxConcurrent   int      // Maximum concurrent conversions

	// UI state
	width  int
	height int
}

// NewConvertWizard creates a new conversion wizard
func NewConvertWizard(reg *registry.Registry) *ConvertWizardModel {
	tier := ui.DetectAnimationTier()

	// Create URL input for adding URLs one at a time
	urlInput := components.NewTextInput("Add GTA5 Mod URL", "https://www.gta5-mods.com/...", 500)
	urlInput.SetValidator(func(s string) error {
		if s == "" {
			return nil // Empty is okay, user might be done adding URLs
		}
		if !strings.Contains(s, "gta5-mods.com") {
			return fmt.Errorf("URL must be from gta5-mods.com")
		}
		return nil
	})

	// Create custom path input
	customPathInput := components.NewTextInput("Custom Resources Path", "", 255)
	customPathInput.SetValidator(func(s string) error {
		if s == "" {
			return fmt.Errorf("Path cannot be empty")
		}
		return nil
	})

	return &ConvertWizardModel{
		step:             ConvertStepSelectServer,
		client:           convert.NewClient(),
		downloader:       download.NewDownloader(2), // Limit concurrent downloads
		registry:         reg,
		urlInput:         urlInput,
		customPathInput:  customPathInput,
		progressBar:      components.NewProgressBar(60),
		spinner:          components.NewSpinner(tier),
		conversions:      make(map[string]*ConversionItem),
		downloadProgress: make(map[string]float64),
		maxConcurrent:    2, // Only 2 conversions at a time to respect rate limits
	}
}

// Init initializes the wizard
func (m *ConvertWizardModel) Init() tea.Cmd {
	return m.setupServerSelector()
}

// setupServerSelector creates the server selector
func (m *ConvertWizardModel) setupServerSelector() tea.Cmd {
	servers := m.registry.List()

	// Add 2 external options + registered servers
	items := make([]components.SelectorItem, len(servers)+2)

	// External server options first
	items[0] = components.SelectorItem{
		Label:       "External Server (Current Directory)",
		Description: "Download to ./resources/ in current directory",
		Value:       "external:current",
	}
	items[1] = components.SelectorItem{
		Label:       "External Server (Custom Path)",
		Description: "Specify a custom path for resources",
		Value:       "external:custom",
	}

	// Add registered servers
	for i, srv := range servers {
		items[i+2] = components.SelectorItem{
			Label:       srv.Name,
			Description: fmt.Sprintf("Port %d • %s", srv.Port, srv.Path),
			Value:       srv,
		}
	}

	m.serverSelector = components.NewSelector("Select Target Server", items)
	m.serverSelector.MaxHeight = 10
	m.serverSelector.Focus()
	return nil
}

// Update handles messages
func (m *ConvertWizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.step == ConvertStepConverting || m.step == ConvertStepDownloading {
				return m, nil // Don't quit during conversion/download
			}
			m.quitting = true
			return m, tea.Quit

		case "enter":
			// In URL input step, Enter adds current URL to list
			if m.step == ConvertStepEnterURLs {
				url := strings.TrimSpace(m.urlInput.Value)
				if url != "" && m.urlInput.Error == "" {
					// Add URL to list
					m.urls = append(m.urls, url)
					// Clear input for next URL
					m.urlInput.Clear()
					return m, nil
				} else if url == "" && len(m.urls) > 0 {
					// Empty input and we have URLs, proceed
					return m.handleEnter()
				}
				return m, nil
			}
			return m.handleEnter()

		case "ctrl+enter":
			// Ctrl+Enter submits in URL input step
			if m.step == ConvertStepEnterURLs {
				return m.handleEnter()
			}
		}

	case conversionStartedMsg:
		item := m.conversions[msg.uuid]
		if item != nil {
			item.UUID = msg.uuid
		}
		return m, nil

	case pollTickMsg:
		// Throttle updates to prevent excessive scrolling
		if time.Since(m.lastUpdate) < 500*time.Millisecond {
			return m, pollTickCmd()
		}
		m.lastUpdate = time.Now()

		// Check conversion progress
		if m.step == ConvertStepConverting && m.pollingActive {
			// Start new conversions from queue if under the limit
			for len(m.conversionQueue) > 0 && m.activeConversions < m.maxConcurrent {
				url := m.conversionQueue[0]
				m.conversionQueue = m.conversionQueue[1:]
				m.activeConversions++

				// Start conversion in background
				go func(u string) {
					uuid, err := m.client.StartConversion(u)
					if err != nil {
						if item := m.conversions[u]; item != nil {
							item.Error = err
						}
						m.activeConversions--
						return
					}

					if item := m.conversions[u]; item != nil {
						item.UUID = uuid
					}
				}(url)

				// Add a small delay between conversion starts
				time.Sleep(200 * time.Millisecond)
			}

			// Poll active conversions for progress
			allComplete := true
			for _, item := range m.conversions {
				if item.Error != nil {
					// Skip failed items
					continue
				}

				if item.UUID != "" && (item.Status == nil || item.Status.Progress < 100) {
					status, err := m.client.QueryProgress(item.UUID)
					if err == nil {
						item.Status = status
						if status.Progress >= 100 {
							item.FileName = status.File
							m.activeConversions--
						}
					}
				}

				if item.Status == nil || item.Status.Progress < 100 {
					allComplete = false
				}
			}

			m.updateConversionProgress()

			// Check if all done (queue empty and all conversions complete)
			if len(m.conversionQueue) == 0 && allComplete && m.activeConversions == 0 {
				m.pollingActive = false
				m.step = ConvertStepDownloading
				return m, downloadFilesCmd(m)
			}
			return m, pollTickCmd()
		}
		return m, nil

	case conversionCompleteMsg:
		m.step = ConvertStepDownloading
		return m, downloadFilesCmd(m)

	case downloadProgressMsg:
		m.downloadProgress[msg.file] = msg.progress
		m.updateDownloadProgress()
		return m, nil

	case downloadCompleteMsg:
		m.step = ConvertStepComplete
		m.completed = true
		return m, nil

	case wizardErrorMsg:
		m.error = string(msg)
		m.step = ConvertStepError
		return m, nil

	case components.SpinnerTickMsg:
		m.spinner.Tick()
		return m, m.spinner.TickCmd()

	case components.CursorBlinkMsg:
		if m.step == ConvertStepCustomPath {
			cmd := m.customPathInput.Update(msg)
			return m, cmd
		}
		if m.step == ConvertStepEnterURLs {
			cmd := m.urlInput.Update(msg)
			return m, cmd
		}
	}

	// Update active component
	switch m.step {
	case ConvertStepSelectServer:
		if m.serverSelector != nil {
			cmd := m.serverSelector.Update(msg)
			cmds = append(cmds, cmd)
		}

	case ConvertStepCustomPath:
		cmd := m.customPathInput.Update(msg)
		cmds = append(cmds, cmd)

	case ConvertStepEnterURLs:
		cmd := m.urlInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleEnter processes Enter key for current step
func (m *ConvertWizardModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case ConvertStepSelectServer:
		if m.serverSelector != nil {
			m.serverSelector.Update(tea.KeyMsg{Type: tea.KeyEnter})

			if m.serverSelector.Confirmed {
				value := m.serverSelector.SelectedValue()

				// Check if it's a registered server
				if srv, ok := value.(types.Server); ok {
					m.selectedServer = &srv
					m.step = ConvertStepEnterURLs
					m.urlInput.Focus()
					return m, m.urlInput.BlinkCmd()
				}

				// Check if it's external server
				if strVal, ok := value.(string); ok {
					if strVal == "external:current" {
						m.externalMode = "current"
						m.step = ConvertStepEnterURLs
						m.urlInput.Focus()
						return m, m.urlInput.BlinkCmd()
					} else if strVal == "external:custom" {
						m.externalMode = "custom"
						m.step = ConvertStepCustomPath
						m.customPathInput.Focus()
						return m, m.customPathInput.BlinkCmd()
					}
				}
			}
		}
		return m, nil

	case ConvertStepCustomPath:
		m.customPathInput.Blur()
		if m.customPathInput.Error != "" {
			return m, nil
		}
		m.customPath = filepath.Clean(m.customPathInput.Value)
		m.step = ConvertStepEnterURLs
		m.urlInput.Focus()
		return m, m.urlInput.BlinkCmd()

	case ConvertStepEnterURLs:
		m.urlInput.Blur()

		// Check if we have any URLs
		if len(m.urls) == 0 {
			return m, nil // Stay on this step
		}

		// Initialize conversion items and queue
		m.conversionQueue = make([]string, len(m.urls))
		copy(m.conversionQueue, m.urls)

		for _, url := range m.urls {
			m.conversions[url] = &ConversionItem{
				URL:      url,
				Category: extractCategory(url),
			}
		}

		m.step = ConvertStepConverting
		m.pollingActive = true
		m.activeConversions = 0
		m.lastUpdate = time.Now()
		return m, tea.Batch(
			m.spinner.TickCmd(),
			pollTickCmd(),
		)

	case ConvertStepComplete, ConvertStepError:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// updateConversionProgress calculates overall conversion progress
func (m *ConvertWizardModel) updateConversionProgress() {
	if len(m.conversions) == 0 {
		return
	}

	total := 0
	for _, item := range m.conversions {
		if item.Status != nil {
			total += item.Status.Progress
		}
	}

	m.overallProgress = float64(total) / float64(len(m.conversions)*100)
	m.progressBar.SetProgress(m.overallProgress)
}

// updateDownloadProgress calculates overall download progress
func (m *ConvertWizardModel) updateDownloadProgress() {
	if len(m.downloads) == 0 {
		return
	}

	total := 0.0
	for _, progress := range m.downloadProgress {
		total += progress
	}

	m.overallProgress = total / float64(len(m.downloads))
	m.progressBar.SetProgress(m.overallProgress)
}

// View renders the wizard
func (m *ConvertWizardModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Background(ui.ColorPrimary).
		Bold(true).
		Padding(0, 2).
		Width(m.width)

	b.WriteString(titleStyle.Render("Convert GTA5 Mods to FiveM Resources"))
	b.WriteString("\n\n")

	// Render current step
	switch m.step {
	case ConvertStepSelectServer:
		if m.serverSelector != nil {
			b.WriteString(m.serverSelector.View())
		}

	case ConvertStepCustomPath:
		b.WriteString(m.customPathInput.View())
		b.WriteString("\n\n")
		helpStyle := lipgloss.NewStyle().
			Foreground(ui.ColorMediumGray).
			Italic(true)
		b.WriteString(helpStyle.Render("Enter the path to your server's resources folder"))

	case ConvertStepEnterURLs:
		// Show list of added URLs
		if len(m.urls) > 0 {
			headerStyle := lipgloss.NewStyle().
				Foreground(ui.ColorPureWhite).
				Bold(true)

			b.WriteString(headerStyle.Render(fmt.Sprintf("Added URLs (%d):", len(m.urls))))
			b.WriteString("\n\n")

			listStyle := lipgloss.NewStyle().
				Foreground(ui.ColorMediumGray)

			for i, url := range m.urls {
				b.WriteString(listStyle.Render(fmt.Sprintf("  %d. %s", i+1, url)))
				b.WriteString("\n")
			}
			b.WriteString("\n")
		}

		// Show input for adding more URLs
		b.WriteString(m.urlInput.View())
		b.WriteString("\n\n")

		helpStyle := lipgloss.NewStyle().
			Foreground(ui.ColorMediumGray).
			Italic(true)

		if len(m.urls) > 0 {
			b.WriteString(helpStyle.Render("Enter: Add URL  •  Enter (empty): Continue  •  Ctrl+Enter: Continue"))
		} else {
			b.WriteString(helpStyle.Render("Enter: Add URL to list  •  Ctrl+Enter: Continue"))
		}

	case ConvertStepConverting:
		b.WriteString(m.renderConverting())

	case ConvertStepDownloading:
		b.WriteString(m.renderDownloading())

	case ConvertStepComplete:
		b.WriteString(m.renderComplete())

	case ConvertStepError:
		b.WriteString(m.renderError())
	}

	// Help text
	if m.step != ConvertStepConverting && m.step != ConvertStepDownloading && m.step != ConvertStepComplete && m.step != ConvertStepError {
		b.WriteString("\n\n")
		helpStyle := lipgloss.NewStyle().
			Foreground(ui.ColorMediumGray).
			Italic(true)
		b.WriteString(helpStyle.Render("Esc: Cancel  •  Enter: Continue"))
	}

	return b.String()
}

// renderConverting renders the conversion progress
func (m *ConvertWizardModel) renderConverting() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Bold(true)

	b.WriteString(headerStyle.Render(fmt.Sprintf("Converting %d Mod(s)", len(m.conversions))))
	b.WriteString("\n\n")

	// Overall progress with queue info
	completedCount := 0
	for _, item := range m.conversions {
		if item.Status != nil && item.Status.Progress >= 100 {
			completedCount++
		}
	}

	progressStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray)

	b.WriteString(progressStyle.Render(fmt.Sprintf("Progress: %d/%d completed  •  %d queued  •  %d/%d active",
		completedCount, len(m.conversions), len(m.conversionQueue), m.activeConversions, m.maxConcurrent)))
	b.WriteString("\n\n")

	// Individual mod statuses (ordered by URL list to maintain consistency)
	i := 1
	for _, url := range m.urls {
		item := m.conversions[url]
		if item == nil {
			continue
		}

		modName := extractModName(url)

		var icon, statusText string
		var statusColor lipgloss.Color

		if item.Error != nil {
			icon = ui.SymbolCross
			statusText = fmt.Sprintf("Failed: %s", item.Error)
			statusColor = ui.ColorError
		} else if item.Status != nil {
			if item.Status.Progress >= 100 {
				icon = ui.SymbolCheck
				statusText = "Complete"
				statusColor = ui.ColorSuccess
			} else if item.Status.Progress > 0 {
				icon = m.spinner.View()
				statusText = fmt.Sprintf("%d%% - Converting", item.Status.Progress)
				statusColor = ui.ColorPrimary
			} else if item.UUID != "" {
				icon = m.spinner.View()
				statusText = "Starting conversion..."
				statusColor = ui.ColorPrimary
			} else {
				icon = "⏳"
				statusText = "Queued"
				statusColor = ui.ColorMediumGray
			}
		} else if item.UUID != "" {
			icon = m.spinner.View()
			statusText = "Starting..."
			statusColor = ui.ColorPrimary
		} else {
			icon = "⏳"
			statusText = "Queued"
			statusColor = ui.ColorMediumGray
		}

		nameStyle := lipgloss.NewStyle().
			Foreground(ui.ColorPureWhite).
			Bold(true)

		statusStyle := lipgloss.NewStyle().
			Foreground(statusColor)

		b.WriteString(fmt.Sprintf("  %d. %s ", i, nameStyle.Render(modName)))
		b.WriteString(statusStyle.Render(fmt.Sprintf("%s %s", icon, statusText)))
		b.WriteString("\n")
		i++
	}

	return b.String()
}

// renderDownloading renders the download progress
func (m *ConvertWizardModel) renderDownloading() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Bold(true)

	b.WriteString(headerStyle.Render(fmt.Sprintf("Downloading %d Resource(s)", len(m.conversions))))
	b.WriteString("\n\n")

	// Overall progress
	completedCount := 0
	for _, item := range m.conversions {
		if item.FileName != "" {
			progress, exists := m.downloadProgress[item.FileName]
			if exists && progress >= 1.0 {
				completedCount++
			}
		}
	}

	progressStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray)

	b.WriteString(progressStyle.Render(fmt.Sprintf("Progress: %d/%d downloaded", completedCount, len(m.conversions))))
	b.WriteString("\n\n")

	// Individual download statuses (ordered by URL list to maintain consistency)
	i := 1
	for _, url := range m.urls {
		item := m.conversions[url]
		if item == nil {
			continue
		}

		modName := extractModName(url)

		var icon, statusText string
		var statusColor lipgloss.Color

		if item.Error != nil {
			icon = ui.SymbolCross
			statusText = "Skipped (conversion failed)"
			statusColor = ui.ColorError
		} else if item.FileName == "" {
			icon = "⏳"
			statusText = "Waiting for conversion..."
			statusColor = ui.ColorMediumGray
		} else {
			progress, exists := m.downloadProgress[item.FileName]
			if !exists {
				icon = "⏳"
				statusText = "Queued"
				statusColor = ui.ColorMediumGray
			} else if progress >= 1.0 {
				icon = ui.SymbolCheck
				statusText = "Complete"
				statusColor = ui.ColorSuccess
			} else {
				icon = m.spinner.View()
				statusText = fmt.Sprintf("%.0f%% - Downloading", progress*100)
				statusColor = ui.ColorPrimary
			}
		}

		nameStyle := lipgloss.NewStyle().
			Foreground(ui.ColorPureWhite).
			Bold(true)

		statusStyle := lipgloss.NewStyle().
			Foreground(statusColor)

		b.WriteString(fmt.Sprintf("  %d. %s ", i, nameStyle.Render(modName)))
		b.WriteString(statusStyle.Render(fmt.Sprintf("%s %s", icon, statusText)))
		b.WriteString("\n")
		i++
	}

	return b.String()
}

// renderComplete renders the completion screen
func (m *ConvertWizardModel) renderComplete() string {
	var b strings.Builder

	// Success banner
	successBanner := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Background(ui.ColorSuccess).
		Bold(true).
		Padding(0, 2).
		MarginBottom(1)

	b.WriteString(successBanner.Render(ui.SymbolCheck + " Conversion Complete"))
	b.WriteString("\n\n")

	// Server info
	nameStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPrimary).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray)

	// Show server or external path info
	if m.externalMode != "" {
		var pathStr string
		if m.externalMode == "current" {
			currentDir, _ := os.Getwd()
			pathStr = filepath.Join(currentDir, "resources")
		} else {
			pathStr = m.customPath
		}
		b.WriteString(labelStyle.Render("Resources Path: "))
		b.WriteString(nameStyle.Render(pathStr))
		b.WriteString("\n\n")
	} else {
		b.WriteString(labelStyle.Render("Server: "))
		b.WriteString(nameStyle.Render(m.selectedServer.Name))
		b.WriteString("\n")

		b.WriteString(labelStyle.Render("Resources Path: "))
		b.WriteString(nameStyle.Render(filepath.Join(m.selectedServer.Path, "resources")))
		b.WriteString("\n\n")
	}

	// Divider
	dividerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray)

	b.WriteString(dividerStyle.Render("────────────────────────────────────────"))
	b.WriteString("\n\n")

	// Summary
	headerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Bold(true)

	b.WriteString(headerStyle.Render(fmt.Sprintf("Converted %d mod(s)", len(m.conversions))))
	b.WriteString("\n\n")

	infoStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray).
		Italic(true)

	b.WriteString(infoStyle.Render("Resources have been extracted and are ready to use!"))
	b.WriteString("\n\n")

	// Divider
	b.WriteString(dividerStyle.Render("────────────────────────────────────────"))
	b.WriteString("\n\n")

	// Exit prompt
	helpStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray).
		Italic(true)

	b.WriteString(helpStyle.Render("Press Enter or Esc to exit"))

	return b.String()
}

// renderError renders the error screen
func (m *ConvertWizardModel) renderError() string {
	var b strings.Builder

	// Error banner
	errorBanner := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Background(ui.ColorError).
		Bold(true).
		Padding(0, 2).
		MarginBottom(1)

	b.WriteString(errorBanner.Render(ui.SymbolCross + " Conversion Failed"))
	b.WriteString("\n\n")

	// Error message
	errorMsgStyle := lipgloss.NewStyle().
		Foreground(ui.ColorError).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(1, 2)

	b.WriteString(errorMsgStyle.Render(m.error))
	b.WriteString("\n\n")

	// Divider
	dividerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray)

	b.WriteString(dividerStyle.Render("────────────────────────────────────────"))
	b.WriteString("\n\n")

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray).
		Italic(true)

	b.WriteString(helpStyle.Render("Press Enter or Esc to exit"))

	return b.String()
}

// Completed returns whether the wizard completed successfully
func (m *ConvertWizardModel) Completed() bool {
	return m.completed
}

// Messages

type conversionStartedMsg struct {
	uuid string
	url  string
}

type pollTickMsg struct{}

type conversionCompleteMsg struct{}

type downloadProgressMsg struct {
	file     string
	progress float64
}

type downloadCompleteMsg struct{}

type wizardErrorMsg string

// Commands

func pollTickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return pollTickMsg{}
	})
}


func downloadFilesCmd(m *ConvertWizardModel) tea.Cmd {
	return func() tea.Msg {
		var resourcesPath string

		// Determine resources path based on server selection
		if m.externalMode == "current" {
			// Current directory
			currentDir, err := os.Getwd()
			if err != nil {
				return wizardErrorMsg(fmt.Sprintf("Failed to get current directory: %v", err))
			}
			resourcesPath = filepath.Join(currentDir, "resources")
		} else if m.externalMode == "custom" {
			// Custom path
			resourcesPath = m.customPath
		} else {
			// Registered server
			resourcesPath = filepath.Join(m.selectedServer.Path, "resources")
		}

		if err := os.MkdirAll(resourcesPath, 0755); err != nil {
			return wizardErrorMsg(fmt.Sprintf("Failed to create resources directory: %v", err))
		}

		var wg sync.WaitGroup
		errChan := make(chan error, len(m.downloads))

		for _, item := range m.conversions {
			if item.FileName == "" {
				continue
			}

			wg.Add(1)
			go func(convItem *ConversionItem) {
				defer wg.Done()

				// Create category subfolder (e.g., [vehicles]/)
				categoryFolder := fmt.Sprintf("[%s]", convItem.Category)
				categoryPath := filepath.Join(resourcesPath, categoryFolder)
				if err := os.MkdirAll(categoryPath, 0755); err != nil {
					errChan <- fmt.Errorf("failed to create category folder: %w", err)
					return
				}

				downloadURL := m.client.GetDownloadURL(convItem.FileName)
				destPath := filepath.Join(resourcesPath, filepath.Base(convItem.FileName))

				// Download using the downloader
				err := m.downloader.Download(downloadURL, destPath, func(progress download.Progress) {
					m.downloadProgress[convItem.FileName] = float64(progress.DownloadedBytes) / float64(progress.TotalBytes)
				})

				if err != nil {
					errChan <- fmt.Errorf("failed to download %s: %w", convItem.FileName, err)
					return
				}

				// Extract zip to category subfolder
				if err := extractZip(destPath, categoryPath); err != nil {
					errChan <- fmt.Errorf("failed to extract %s: %w", convItem.FileName, err)
					return
				}

				// Remove zip file after extraction
				os.Remove(destPath)
			}(item)
		}

		wg.Wait()
		close(errChan)

		// Check for errors
		if len(errChan) > 0 {
			return wizardErrorMsg(fmt.Sprintf("Download failed: %v", <-errChan))
		}

		return downloadCompleteMsg{}
	}
}

// extractZip extracts a zip file to the destination directory
func extractZip(zipPath, destPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destPath, f.Name)

		// Check for ZipSlip vulnerability
		cleanDest := filepath.Clean(destPath)
		cleanPath := filepath.Clean(fpath)
		if !strings.HasPrefix(cleanPath, cleanDest) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// extractCategory extracts the mod category from a gta5-mods.com URL
// e.g., "https://www.gta5-mods.com/vehicles/..." -> "vehicles"
func extractCategory(url string) string {
	// Split URL by "/" and find the category after gta5-mods.com
	parts := strings.Split(url, "/")
	for i, part := range parts {
		if part == "www.gta5-mods.com" || part == "gta5-mods.com" {
			if i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}
	return "misc" // Default category
}

// extractModName extracts a readable mod name from a gta5-mods.com URL
// e.g., "https://www.gta5-mods.com/vehicles/1995-mclaren-f1-lm-addon" -> "1995 McLaren F1 LM Addon"
func extractModName(url string) string {
	// Split URL by "/" and get the last part (slug)
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return url
	}

	slug := parts[len(parts)-1]

	// Remove any query parameters
	if idx := strings.Index(slug, "?"); idx != -1 {
		slug = slug[:idx]
	}

	// Replace hyphens with spaces and title case
	name := strings.ReplaceAll(slug, "-", " ")

	// Simple title case: capitalize first letter of each word
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	name = strings.Join(words, " ")

	// Limit length for display
	if len(name) > 50 {
		name = name[:47] + "..."
	}

	return name
}
