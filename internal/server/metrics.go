package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/VexoaXYZ/inkwash/pkg/types"
)

// MetricsCollector collects server metrics in background
type MetricsCollector struct {
	servers  map[string]*types.ServerMetrics
	interval time.Duration
	stopChan chan struct{}
	mu       sync.RWMutex
	pm       *ProcessManager
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(interval time.Duration) *MetricsCollector {
	if interval == 0 {
		interval = 2 * time.Second
	}

	return &MetricsCollector{
		servers:  make(map[string]*types.ServerMetrics),
		interval: interval,
		stopChan: make(chan struct{}),
		pm:       NewProcessManager(),
	}
}

// Start starts the metrics collection loop
func (mc *MetricsCollector) Start() {
	go mc.collectLoop()
}

// Stop stops the metrics collection
func (mc *MetricsCollector) Stop() {
	close(mc.stopChan)
}

// Track adds a server to track
func (mc *MetricsCollector) Track(server *types.Server) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if server.IsRunning() {
		mc.servers[server.Name] = types.NewServerMetrics(server.PID)
	}
}

// Untrack removes a server from tracking
func (mc *MetricsCollector) Untrack(serverName string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.servers, serverName)
}

// Get returns metrics for a server
func (mc *MetricsCollector) Get(serverName string) *types.ServerMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	return mc.servers[serverName]
}

// GetAll returns all tracked metrics
func (mc *MetricsCollector) GetAll() map[string]*types.ServerMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// Return a copy
	metrics := make(map[string]*types.ServerMetrics)
	for k, v := range mc.servers {
		metrics[k] = v
	}

	return metrics
}

// collectLoop runs the collection loop
func (mc *MetricsCollector) collectLoop() {
	ticker := time.NewTicker(mc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-mc.stopChan:
			return
		case <-ticker.C:
			mc.collect()
		}
	}
}

// collect collects metrics for all tracked servers
func (mc *MetricsCollector) collect() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for name, metrics := range mc.servers {
		if err := mc.collectOne(metrics); err != nil {
			// If collection fails, the process may have stopped
			// Remove from tracking
			delete(mc.servers, name)
		}
	}
}

// collectOne collects metrics for a single server
func (mc *MetricsCollector) collectOne(metrics *types.ServerMetrics) error {
	proc, err := process.NewProcess(int32(metrics.PID))
	if err != nil {
		return fmt.Errorf("process not found: %w", err)
	}

	// Collect CPU percentage
	cpu, err := proc.CPUPercent()
	if err == nil {
		metrics.AddCPUSample(cpu)
	}

	// Collect memory usage
	memInfo, err := proc.MemoryInfo()
	if err == nil {
		ramGB := float64(memInfo.RSS) / 1024 / 1024 / 1024
		metrics.AddRAMSample(ramGB)
	}

	// Collect network I/O
	ioCounters, err := proc.IOCounters()
	if err == nil {
		// Calculate delta from last measurement
		if metrics.LastUpdate.IsZero() {
			metrics.NetworkTX = 0
			metrics.NetworkRX = 0
		} else {
			elapsed := time.Since(metrics.LastUpdate).Seconds()
			if elapsed > 0 {
				txDelta := ioCounters.WriteBytes
				rxDelta := ioCounters.ReadBytes

				metrics.NetworkTX = uint64(float64(txDelta) / elapsed)
				metrics.NetworkRX = uint64(float64(rxDelta) / elapsed)
			}
		}
	}

	// TODO: Get player count from server logs or query endpoint
	// For now, set to 0
	metrics.PlayerCount = 0

	metrics.LastUpdate = time.Now()

	return nil
}

// UpdatePlayerCount manually updates player count for a server
func (mc *MetricsCollector) UpdatePlayerCount(serverName string, count int) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if metrics, ok := mc.servers[serverName]; ok {
		metrics.PlayerCount = count
	}
}
