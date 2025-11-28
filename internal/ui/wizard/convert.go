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

	// Progress components
	progressBar *components.ProgressBar
	spinner     *components.Spinner

	// State
	selectedServer *types.Server
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

	// UI state
	width  int
	height int
}

// NewConvertWizard creates a new conversion wizard
func NewConvertWizard(reg *registry.Registry) *ConvertWizardModel {
	tier := ui.DetectAnimationTier()

	// Create multi-line input for URLs
	urlInput := components.NewTextInput("GTA5 Mod URLs (one per line)", "https://www.gta5-mods.com/...", 2000)
	urlInput.SetValidator(func(s string) error {
		if s == "" {
			return fmt.Errorf("Please enter at least one URL")
		}
		// Validate each line is a gta5-mods.com URL
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.Contains(line, "gta5-mods.com") {
				return fmt.Errorf("All URLs must be from gta5-mods.com")
			}
		}
		return nil
	})

	return &ConvertWizardModel{
		step:             ConvertStepSelectServer,
		client:           convert.NewClient(),
		downloader:       download.NewDownloader(3),
		registry:         reg,
		urlInput:         urlInput,
		progressBar:      components.NewProgressBar(60),
		spinner:          components.NewSpinner(tier),
		conversions:      make(map[string]*ConversionItem),
		downloadProgress: make(map[string]float64),
	}
}

// Init initializes the wizard
func (m *ConvertWizardModel) Init() tea.Cmd {
	return m.setupServerSelector()
}

// setupServerSelector creates the server selector
func (m *ConvertWizardModel) setupServerSelector() tea.Cmd {
	servers := m.registry.List()
	if len(servers) == 0 {
		m.error = "No servers found. Please create a server first."
		m.step = ConvertStepError
		return nil
	}

	items := make([]components.SelectorItem, len(servers))
	for i, srv := range servers {
		items[i] = components.SelectorItem{
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
			return m.handleEnter()
		}

	case conversionStartedMsg:
		item := m.conversions[msg.uuid]
		if item != nil {
			item.UUID = msg.uuid
		}
		return m, nil

	case conversionProgressMsg:
		if _, ok := m.conversions[string(msg)]; ok {
			// Update status
			m.updateConversionProgress()
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
				if srv, ok := m.serverSelector.SelectedValue().(types.Server); ok {
					m.selectedServer = &srv
					m.step = ConvertStepEnterURLs
					m.urlInput.Focus()
					return m, m.urlInput.BlinkCmd()
				}
			}
		}
		return m, nil

	case ConvertStepEnterURLs:
		m.urlInput.Blur()
		if m.urlInput.Error != "" {
			return m, nil
		}

		// Parse URLs
		lines := strings.Split(m.urlInput.Value, "\n")
		var urls []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				urls = append(urls, line)
			}
		}
		m.urls = urls

		// Initialize conversion items
		for _, url := range urls {
			m.conversions[url] = &ConversionItem{
				URL:      url,
				Category: extractCategory(url),
			}
		}

		m.step = ConvertStepConverting
		return m, tea.Batch(
			startConversionsCmd(m),
			m.spinner.TickCmd(),
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

	case ConvertStepEnterURLs:
		b.WriteString(m.urlInput.View())

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

	b.WriteString(headerStyle.Render("Converting Mods"))
	b.WriteString("\n\n")

	// Spinner and status
	spinnerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPrimary)

	b.WriteString(spinnerStyle.Render(m.spinner.View()))
	b.WriteString(" Converting...")
	b.WriteString("\n\n")

	// Progress bar
	b.WriteString(m.progressBar.Render())
	b.WriteString("\n\n")

	// Individual conversion statuses
	for _, item := range m.conversions {
		statusStyle := lipgloss.NewStyle().
			Foreground(ui.ColorMediumGray)

		status := "Waiting..."
		if item.Status != nil {
			status = fmt.Sprintf("%d%% - %s", item.Status.Progress, item.Status.Message)
		}
		if item.Error != nil {
			statusStyle = statusStyle.Foreground(ui.ColorError)
			status = fmt.Sprintf("Error: %s", item.Error)
		}

		b.WriteString(statusStyle.Render(fmt.Sprintf("  %s", status)))
		b.WriteString("\n")
	}

	return b.String()
}

// renderDownloading renders the download progress
func (m *ConvertWizardModel) renderDownloading() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Bold(true)

	b.WriteString(headerStyle.Render("Downloading Resources"))
	b.WriteString("\n\n")

	// Spinner and status
	spinnerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPrimary)

	b.WriteString(spinnerStyle.Render(m.spinner.View()))
	b.WriteString(" Downloading...")
	b.WriteString("\n\n")

	// Progress bar
	b.WriteString(m.progressBar.Render())
	b.WriteString("\n\n")

	// Individual download statuses
	for file, progress := range m.downloadProgress {
		statusStyle := lipgloss.NewStyle().
			Foreground(ui.ColorMediumGray)

		b.WriteString(statusStyle.Render(fmt.Sprintf("  %s: %.0f%%", filepath.Base(file), progress*100)))
		b.WriteString("\n")
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

	b.WriteString(labelStyle.Render("Server: "))
	b.WriteString(nameStyle.Render(m.selectedServer.Name))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Resources Path: "))
	b.WriteString(nameStyle.Render(filepath.Join(m.selectedServer.Path, "resources")))
	b.WriteString("\n\n")

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
}

type conversionProgressMsg string // UUID

type conversionCompleteMsg struct{}

type downloadProgressMsg struct {
	file     string
	progress float64
}

type downloadCompleteMsg struct{}

type wizardErrorMsg string

// Commands

func startConversionsCmd(m *ConvertWizardModel) tea.Cmd {
	return func() tea.Msg {
		var wg sync.WaitGroup
		errChan := make(chan error, len(m.urls))

		// Start all conversions
		for _, url := range m.urls {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()

				uuid, err := m.client.StartConversion(u)
				if err != nil {
					errChan <- fmt.Errorf("failed to start conversion for %s: %w", u, err)
					return
				}

				if item, ok := m.conversions[u]; ok {
					item.UUID = uuid
					m.conversionList = append(m.conversionList, uuid)
				}

				// Poll until complete
				for {
					status, err := m.client.QueryProgress(uuid)
					if err != nil {
						errChan <- fmt.Errorf("failed to query progress: %w", err)
						return
					}

					if item, ok := m.conversions[u]; ok {
						item.Status = status
						if status.Progress >= 100 {
							item.FileName = status.File
							m.downloads = append(m.downloads, status.File)
							break
						}
					}

					time.Sleep(2 * time.Second)
				}
			}(url)
		}

		wg.Wait()
		close(errChan)

		// Check for errors
		if len(errChan) > 0 {
			return wizardErrorMsg(fmt.Sprintf("Conversion failed: %v", <-errChan))
		}

		return conversionCompleteMsg{}
	}
}

func downloadFilesCmd(m *ConvertWizardModel) tea.Cmd {
	return func() tea.Msg {
		resourcesPath := filepath.Join(m.selectedServer.Path, "resources")
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
