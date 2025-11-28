package wizard

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/VexoaXYZ/inkwash/internal/cache"
	"github.com/VexoaXYZ/inkwash/internal/download"
	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/server"
	"github.com/VexoaXYZ/inkwash/internal/ui"
	"github.com/VexoaXYZ/inkwash/internal/ui/components"
	"github.com/VexoaXYZ/inkwash/internal/validation"
	"github.com/VexoaXYZ/inkwash/pkg/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WizardStep represents the current step in the wizard
type WizardStep int

const (
	StepServerName WizardStep = iota
	StepBuild
	StepLicenseKey
	StepPort
	StepPath
	StepConfirm
	StepInstalling
	StepComplete
	StepError
)

// CreateWizardModel represents the state of the creation wizard
type CreateWizardModel struct {
	step          WizardStep
	installer     *server.Installer
	artifactClient *download.ArtifactClient
	keyVault      *cache.KeyVault
	registry      *registry.Registry

	// Input components
	nameInput     *components.TextInput
	portInput     *components.TextInput
	pathInput     *components.TextInput
	buildSelector *components.Selector
	keySelector   *components.Selector

	// Progress components
	progressBar   *components.ProgressBar
	spinner       *components.Spinner

	// State
	serverName    string
	buildNumber   int
	licenseKey    string
	port          int
	installPath   string
	builds        []types.Build
	keys          []cache.LicenseKey
	error         string
	installProgress server.InstallProgress
	quitting      bool
	completed     bool

	// Loading states
	loadingBuilds bool
	loadingKeys   bool
	width         int
	height        int
}

// NewCreateWizard creates a new creation wizard
func NewCreateWizard(installer *server.Installer, keyVault *cache.KeyVault, reg *registry.Registry) *CreateWizardModel {
	tier := ui.DetectAnimationTier()

	// Create input components
	nameInput := components.NewTextInput("Server Name", "My FiveM Server", 50)
	nameInput.SetValidator(func(s string) error {
		if s == "" {
			return fmt.Errorf("Server name cannot be empty")
		}
		if reg.Exists(s) {
			return fmt.Errorf("Server '%s' already exists", s)
		}
		return nil
	})

	portInput := components.NewTextInput("Port", "30120", 5)
	portInput.Value = "30120"
	portInput.SetValidator(func(s string) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Port must be a number")
		}
		if port < 1024 || port > 65535 {
			return fmt.Errorf("Port must be between 1024 and 65535")
		}
		return nil
	})

	defaultPath := filepath.Join(registry.GetDefaultConfigPath(), "servers")
	pathInput := components.NewTextInput("Installation Path", defaultPath, 255)
	pathInput.Value = defaultPath

	return &CreateWizardModel{
		step:           StepServerName,
		installer:      installer,
		artifactClient: download.NewArtifactClient(),
		keyVault:       keyVault,
		registry:       reg,
		nameInput:      nameInput,
		portInput:      portInput,
		pathInput:      pathInput,
		progressBar:    components.NewProgressBar(60),
		spinner:        components.NewSpinner(tier),
		port:           30120,
	}
}

// Init initializes the wizard
func (m *CreateWizardModel) Init() tea.Cmd {
	m.nameInput.Focus()
	return m.nameInput.BlinkCmd()
}

// Update handles messages
func (m *CreateWizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.step == StepInstalling {
				return m, nil // Don't quit during installation
			}
			m.quitting = true
			return m, tea.Quit

		case "enter":
			return m.handleEnter()
		}

	case buildsLoadedMsg:
		m.builds = msg.builds
		m.loadingBuilds = false
		return m.setupBuildSelector(), nil

	case keysLoadedMsg:
		m.keys = msg.keys
		m.loadingKeys = false
		return m.setupKeySelector(), nil

	case installProgressMsg:
		m.installProgress = server.InstallProgress(msg)
		if m.installProgress.Progress >= 1.0 {
			m.step = StepComplete
			m.completed = true
		}
		m.progressBar.SetProgress(m.installProgress.Progress)
		return m, nil

	case installErrorMsg:
		m.error = string(msg)
		m.step = StepError
		return m, nil

	case components.SpinnerTickMsg:
		m.spinner.Tick()
		return m, m.spinner.TickCmd()

	case components.CursorBlinkMsg:
		// Pass to active input
		switch m.step {
		case StepServerName:
			cmd := m.nameInput.Update(msg)
			return m, cmd
		case StepPort:
			cmd := m.portInput.Update(msg)
			return m, cmd
		case StepPath:
			cmd := m.pathInput.Update(msg)
			return m, cmd
		}
	}

	// Update active component
	switch m.step {
	case StepServerName:
		cmd := m.nameInput.Update(msg)
		cmds = append(cmds, cmd)

	case StepBuild:
		if m.buildSelector != nil {
			cmd := m.buildSelector.Update(msg)
			cmds = append(cmds, cmd)
		}

	case StepLicenseKey:
		if m.keySelector != nil {
			cmd := m.keySelector.Update(msg)
			cmds = append(cmds, cmd)
		}

	case StepPort:
		cmd := m.portInput.Update(msg)
		cmds = append(cmds, cmd)

	case StepPath:
		cmd := m.pathInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleEnter processes Enter key for current step
func (m *CreateWizardModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepServerName:
		m.nameInput.Blur()
		if m.nameInput.Error != "" {
			return m, nil
		}
		m.serverName = m.nameInput.Value
		m.step = StepBuild
		m.loadingBuilds = true
		return m, tea.Batch(
			loadBuildsCmd(m.artifactClient),
			m.spinner.TickCmd(),
		)

	case StepBuild:
		if m.buildSelector != nil && m.buildSelector.Confirmed {
			if build, ok := m.buildSelector.SelectedValue().(types.Build); ok {
				m.buildNumber = build.Number
				m.step = StepLicenseKey
				m.loadingKeys = true
				return m, tea.Batch(
					loadKeysCmd(m.keyVault),
					m.spinner.TickCmd(),
				)
			}
		}
		// Let the normal update flow handle the Enter key
		return m, nil

	case StepLicenseKey:
		if m.keySelector != nil && m.keySelector.Confirmed {
			if key, ok := m.keySelector.SelectedValue().(string); ok {
				m.licenseKey = key
				m.step = StepPort
				m.portInput.Focus()
				return m, m.portInput.BlinkCmd()
			}
		}
		// Let the normal update flow handle the Enter key
		return m, nil

	case StepPort:
		m.portInput.Blur()
		if m.portInput.Error != "" {
			return m, nil
		}
		port, _ := strconv.Atoi(m.portInput.Value)
		m.port = port
		m.step = StepPath
		m.pathInput.Focus()
		return m, m.pathInput.BlinkCmd()

	case StepPath:
		m.pathInput.Blur()
		m.installPath = m.pathInput.Value
		m.step = StepConfirm

	case StepConfirm:
		m.step = StepInstalling
		return m, tea.Batch(
			installServerCmd(m),
			m.spinner.TickCmd(),
		)

	case StepComplete, StepError:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// setupBuildSelector creates the build selector with loaded builds
func (m *CreateWizardModel) setupBuildSelector() *CreateWizardModel {
	items := make([]components.SelectorItem, len(m.builds))
	for i, build := range m.builds {
		label := fmt.Sprintf("Build %d", build.Number)
		desc := ""
		if build.Recommended {
			label += " (Recommended)"
			desc = "Stable build recommended for production"
		} else if build.Optional {
			label += " (Optional)"
			desc = "Latest features, may be unstable"
		}

		items[i] = components.SelectorItem{
			Label:       label,
			Description: desc,
			Value:       build,
		}
	}

	m.buildSelector = components.NewSelector("Select FXServer Build", items)
	m.buildSelector.MaxHeight = 10
	m.buildSelector.Focus()
	return m
}

// setupKeySelector creates the key selector with loaded keys
func (m *CreateWizardModel) setupKeySelector() *CreateWizardModel {
	items := make([]components.SelectorItem, len(m.keys)+1)

	// Add existing keys
	for i, key := range m.keys {
		items[i] = components.SelectorItem{
			Label:       key.Label,
			Description: validation.MaskKey(key.Key),
			Value:       key.Key,
		}
	}

	// Add manual entry option
	items[len(m.keys)] = components.SelectorItem{
		Label:       "Enter manually",
		Description: "Type your license key",
		Value:       "manual",
	}

	m.keySelector = components.NewSelector("Select License Key", items)
	m.keySelector.MaxHeight = 10
	m.keySelector.Focus()
	return m
}

// View renders the wizard
func (m *CreateWizardModel) View() string {
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

	b.WriteString(titleStyle.Render("Create New FiveM Server"))
	b.WriteString("\n\n")

	// Step indicator
	stepStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray)

	stepNum := int(m.step) + 1
	totalSteps := 6 // Not counting Installing, Complete, Error
	if m.step >= StepInstalling {
		stepNum = totalSteps
	}

	b.WriteString(stepStyle.Render(fmt.Sprintf("Step %d of %d", stepNum, totalSteps)))
	b.WriteString("\n\n")

	// Render current step
	switch m.step {
	case StepServerName:
		b.WriteString(m.nameInput.View())

	case StepBuild:
		if m.loadingBuilds {
			b.WriteString(m.spinner.View())
			b.WriteString(" Loading available builds...")
		} else if m.buildSelector != nil {
			b.WriteString(m.buildSelector.View())
		}

	case StepLicenseKey:
		if m.loadingKeys {
			b.WriteString(m.spinner.View())
			b.WriteString(" Loading license keys...")
		} else if m.keySelector != nil {
			b.WriteString(m.keySelector.View())
		}

	case StepPort:
		b.WriteString(m.portInput.View())

	case StepPath:
		b.WriteString(m.pathInput.View())

	case StepConfirm:
		b.WriteString(m.renderConfirmation())

	case StepInstalling:
		b.WriteString(m.renderProgress())

	case StepComplete:
		b.WriteString(m.renderComplete())

	case StepError:
		b.WriteString(m.renderError())
	}

	// Help text
	if m.step != StepInstalling && m.step != StepComplete && m.step != StepError {
		b.WriteString("\n\n")
		helpStyle := lipgloss.NewStyle().
			Foreground(ui.ColorMediumGray).
			Italic(true)
		b.WriteString(helpStyle.Render("Esc: Cancel  •  Enter: Continue"))
	}

	return b.String()
}

// renderConfirmation renders the confirmation step
func (m *CreateWizardModel) renderConfirmation() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPureWhite).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray)

	valueStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPrimary)

	b.WriteString(headerStyle.Render("Confirm Configuration"))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Server Name:    "))
	b.WriteString(valueStyle.Render(m.serverName))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Build Number:   "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", m.buildNumber)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("License Key:    "))
	b.WriteString(valueStyle.Render(validation.MaskKey(m.licenseKey)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Port:           "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", m.port)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Install Path:   "))
	b.WriteString(valueStyle.Render(m.installPath))
	b.WriteString("\n\n")

	b.WriteString(headerStyle.Render("Press Enter to start installation"))

	return b.String()
}

// renderProgress renders the installation progress
func (m *CreateWizardModel) renderProgress() string {
	var b strings.Builder

	b.WriteString(m.spinner.View())
	b.WriteString(" ")
	b.WriteString(m.installProgress.Step)
	b.WriteString("\n\n")

	b.WriteString(m.progressBar.Render())
	b.WriteString("\n\n")

	progressStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray)

	b.WriteString(progressStyle.Render(fmt.Sprintf("Step %d of %d",
		m.installProgress.CompletedSteps, m.installProgress.TotalSteps)))

	if m.installProgress.CurrentFile != "" {
		b.WriteString("\n")
		b.WriteString(progressStyle.Render(m.installProgress.CurrentFile))
	}

	return b.String()
}

// renderComplete renders the completion screen
func (m *CreateWizardModel) renderComplete() string {
	var b strings.Builder

	successStyle := lipgloss.NewStyle().
		Foreground(ui.ColorSuccess).
		Bold(true)

	b.WriteString(successStyle.Render(ui.SymbolCheck + " Server created successfully!"))
	b.WriteString("\n\n")

	codeStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPrimary)

	b.WriteString("Start your server:\n")
	b.WriteString(codeStyle.Render(fmt.Sprintf("  inkwash start %s", m.serverName)))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray).
		Italic(true)

	b.WriteString(helpStyle.Render("Press Enter to exit"))

	return b.String()
}

// renderError renders the error screen
func (m *CreateWizardModel) renderError() string {
	var b strings.Builder

	errorStyle := lipgloss.NewStyle().
		Foreground(ui.ColorError).
		Bold(true)

	b.WriteString(errorStyle.Render(ui.SymbolCross + " Installation failed"))
	b.WriteString("\n\n")

	b.WriteString(m.error)
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMediumGray).
		Italic(true)

	b.WriteString(helpStyle.Render("Press Enter to exit"))

	return b.String()
}

// Completed returns whether the wizard completed successfully
func (m *CreateWizardModel) Completed() bool {
	return m.completed
}

// ServerName returns the created server name
func (m *CreateWizardModel) ServerName() string {
	return m.serverName
}

// Messages

type buildsLoadedMsg struct {
	builds []types.Build
}

type keysLoadedMsg struct {
	keys []cache.LicenseKey
}

type installProgressMsg server.InstallProgress

type installErrorMsg string

// Commands

func loadBuildsCmd(client *download.ArtifactClient) tea.Cmd {
	return func() tea.Msg {
		builds, err := client.FetchBuilds()
		if err != nil {
			return installErrorMsg(fmt.Sprintf("Failed to fetch builds: %v", err))
		}
		return buildsLoadedMsg{builds: builds}
	}
}

func loadKeysCmd(vault *cache.KeyVault) tea.Cmd {
	return func() tea.Msg {
		return keysLoadedMsg{keys: vault.List()}
	}
}

func installServerCmd(m *CreateWizardModel) tea.Cmd {
	return func() tea.Msg {
		err := m.installer.Install(
			m.serverName,
			m.installPath,
			m.buildNumber,
			m.licenseKey,
			m.port,
			func(progress server.InstallProgress) {
				// This will be called from a different goroutine
				// We need to send it back as a message
			},
		)

		if err != nil {
			return installErrorMsg(fmt.Sprintf("Installation failed: %v", err))
		}

		return installProgressMsg{
			Step:           "Complete",
			Progress:       1.0,
			TotalSteps:     8,
			CompletedSteps: 8,
		}
	}
}
