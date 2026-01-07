package mock

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"aperture-science-network/internal/stats"
)

// Ensure StatsProvider implements stats.Provider
var _ stats.Provider = (*StatsProvider)(nil)

// StatsProvider is a mock implementation for system stats
type StatsProvider struct {
	mu        sync.Mutex
	startTime time.Time
	baseCPU   float64
}

// NewStatsProvider creates a new mock stats provider
func NewStatsProvider() *StatsProvider {
	return &StatsProvider{
		startTime: time.Now(),
		baseCPU:   30.0,
	}
}

// GetSystemStats returns mock system statistics with slight variations
func (p *StatsProvider) GetSystemStats() (*stats.SystemStats, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Simulate CPU usage oscillating between 15-45%
	elapsed := time.Since(p.startTime).Seconds()
	cpuBase := 30.0 + 15.0*math.Sin(elapsed/10.0)
	cpuVariation := (rand.Float64() - 0.5) * 10
	cpuUsage := cpuBase + cpuVariation
	if cpuUsage < 5 {
		cpuUsage = 5
	}
	if cpuUsage > 95 {
		cpuUsage = 95
	}

	// Memory: ~60% used with slight variations
	memTotal := uint64(16 * 1024 * 1024 * 1024) // 16 GB
	memBaseUsed := uint64(10 * 1024 * 1024 * 1024)
	memVariation := uint64((rand.Float64() - 0.5) * float64(512*1024*1024))
	memUsed := memBaseUsed + memVariation
	memPercent := float64(memUsed) / float64(memTotal) * 100

	// Disk: ~40% used (stable)
	diskTotal := uint64(500 * 1024 * 1024 * 1024) // 500 GB
	diskUsed := uint64(200 * 1024 * 1024 * 1024)  // 200 GB
	diskPercent := float64(diskUsed) / float64(diskTotal) * 100

	// Uptime: incrementing from start
	uptime := uint64(time.Since(p.startTime).Seconds()) + 86400 // +1 day base

	return &stats.SystemStats{
		CPUUsage:      cpuUsage,
		CPUCores:      4,
		MemoryUsed:    memUsed,
		MemoryTotal:   memTotal,
		MemoryPercent: memPercent,
		DiskUsed:      diskUsed,
		DiskTotal:     diskTotal,
		DiskPercent:   diskPercent,
		Uptime:        uptime,
		Hostname:      "debug-server",
		OS:            "linux",
		Platform:      "debian",
	}, nil
}
