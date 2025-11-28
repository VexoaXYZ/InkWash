package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/VexoaXYZ/inkwash/pkg/types"
)

// ProcessManager handles server process lifecycle
type ProcessManager struct{}

// NewProcessManager creates a new process manager
func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

// Start starts a server process
func (pm *ProcessManager) Start(server *types.Server) error {
	if server.IsRunning() {
		return fmt.Errorf("server '%s' is already running (PID: %d)", server.Name, server.PID)
	}

	// Get script path
	scriptPath := pm.getScriptPath(server)

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("launch script not found: %s", scriptPath)
	}

	// Create command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", scriptPath)
	} else {
		cmd = exec.Command("bash", scriptPath)
	}

	cmd.Dir = server.Path

	// Create logs directory
	logsDir := filepath.Join(server.Path, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Redirect output to log file
	logPath := filepath.Join(logsDir, "server.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Start process in background
	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Close log file after command exits (in background)
	go func() {
		cmd.Wait()
		logFile.Close()
	}()

	server.PID = cmd.Process.Pid
	server.LastStarted = time.Now()

	return nil
}

// Stop stops a server process
func (pm *ProcessManager) Stop(server *types.Server) error {
	if !server.IsRunning() {
		return fmt.Errorf("server '%s' is not running", server.Name)
	}

	proc, err := process.NewProcess(int32(server.PID))
	if err != nil {
		// Process doesn't exist, update PID
		server.PID = 0
		return nil
	}

	// Graceful shutdown
	if runtime.GOOS == "windows" {
		// On Windows, use taskkill for graceful termination
		cmd := exec.Command("taskkill", "/PID", strconv.Itoa(server.PID), "/T")
		if err := cmd.Run(); err != nil {
			// If graceful fails, force kill
			cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(server.PID), "/T")
			cmd.Run()
		}
	} else {
		// On Linux, send SIGTERM
		if err := proc.SendSignal(syscall.SIGTERM); err != nil {
			// If SIGTERM fails, send SIGKILL
			proc.Kill()
		}
	}

	// Wait for shutdown (timeout 30s)
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			// Force kill if still running
			proc.Kill()
			server.PID = 0
			return nil

		case <-ticker.C:
			exists, _ := process.PidExists(int32(server.PID))
			if !exists {
				server.PID = 0
				return nil
			}
		}
	}
}

// IsRunning checks if a server process is actually running
func (pm *ProcessManager) IsRunning(server *types.Server) bool {
	if server.PID == 0 {
		return false
	}

	exists, err := process.PidExists(int32(server.PID))
	if err != nil || !exists {
		return false
	}

	return true
}

// GetStatus returns detailed process status
func (pm *ProcessManager) GetStatus(server *types.Server) string {
	if !pm.IsRunning(server) {
		return "Stopped"
	}

	proc, err := process.NewProcess(int32(server.PID))
	if err != nil {
		return "Unknown"
	}

	status, err := proc.Status()
	if err != nil {
		return "Running"
	}

	// status is an array, get first element
	if len(status) > 0 {
		return status[0]
	}

	return "Running"
}

// Restart restarts a server
func (pm *ProcessManager) Restart(server *types.Server) error {
	if pm.IsRunning(server) {
		if err := pm.Stop(server); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}

		// Wait a moment for cleanup
		time.Sleep(2 * time.Second)
	}

	return pm.Start(server)
}

// getScriptPath returns the launch script path for a server
func (pm *ProcessManager) getScriptPath(server *types.Server) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(server.Path, "run.cmd")
	}
	return filepath.Join(server.Path, "run.sh")
}

// GetMemoryUsage returns memory usage in bytes
func (pm *ProcessManager) GetMemoryUsage(server *types.Server) (uint64, error) {
	if !pm.IsRunning(server) {
		return 0, fmt.Errorf("server is not running")
	}

	proc, err := process.NewProcess(int32(server.PID))
	if err != nil {
		return 0, err
	}

	memInfo, err := proc.MemoryInfo()
	if err != nil {
		return 0, err
	}

	return memInfo.RSS, nil
}

// GetCPUPercent returns CPU usage percentage
func (pm *ProcessManager) GetCPUPercent(server *types.Server) (float64, error) {
	if !pm.IsRunning(server) {
		return 0, fmt.Errorf("server is not running")
	}

	proc, err := process.NewProcess(int32(server.PID))
	if err != nil {
		return 0, err
	}

	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		return 0, err
	}

	return cpuPercent, nil
}
