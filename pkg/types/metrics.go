package types

import "time"

// ServerMetrics represents real-time metrics for a server
type ServerMetrics struct {
	PID         int
	RAM         []float64 // Last 20 samples (for sparkline) in GB
	CPU         []float64 // Last 20 samples (percentage)
	NetworkTX   uint64    // Bytes transmitted per second
	NetworkRX   uint64    // Bytes received per second
	PlayerCount int
	LastUpdate  time.Time
}

// NewServerMetrics creates a new ServerMetrics instance
func NewServerMetrics(pid int) *ServerMetrics {
	return &ServerMetrics{
		PID:        pid,
		RAM:        make([]float64, 20),
		CPU:        make([]float64, 20),
		LastUpdate: time.Now(),
	}
}

// AddRAMSample adds a RAM usage sample (sliding window)
func (m *ServerMetrics) AddRAMSample(ramGB float64) {
	m.RAM = append(m.RAM[1:], ramGB)
}

// AddCPUSample adds a CPU usage sample (sliding window)
func (m *ServerMetrics) AddCPUSample(cpuPercent float64) {
	m.CPU = append(m.CPU[1:], cpuPercent)
}

// CurrentRAM returns the most recent RAM usage in GB
func (m *ServerMetrics) CurrentRAM() float64 {
	if len(m.RAM) == 0 {
		return 0
	}
	return m.RAM[len(m.RAM)-1]
}

// CurrentCPU returns the most recent CPU usage percentage
func (m *ServerMetrics) CurrentCPU() float64 {
	if len(m.CPU) == 0 {
		return 0
	}
	return m.CPU[len(m.CPU)-1]
}
