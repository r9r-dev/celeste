package stats

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemStats struct {
	CPUUsage      float64 `json:"cpuUsage"`
	CPUCores      int     `json:"cpuCores"`
	MemoryUsed    uint64  `json:"memoryUsed"`
	MemoryTotal   uint64  `json:"memoryTotal"`
	MemoryPercent float64 `json:"memoryPercent"`
	DiskUsed      uint64  `json:"diskUsed"`
	DiskTotal     uint64  `json:"diskTotal"`
	DiskPercent   float64 `json:"diskPercent"`
	Uptime        uint64  `json:"uptime"`
	Hostname      string  `json:"hostname"`
	OS            string  `json:"os"`
	Platform      string  `json:"platform"`
}

func GetSystemStats() (*SystemStats, error) {
	// CPU
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	cpuInfo, err := cpu.Info()
	cores := 0
	if err == nil && len(cpuInfo) > 0 {
		cores = int(cpuInfo[0].Cores)
	}

	// Memory
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Disk
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	// Host
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	cpuUsage := 0.0
	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	return &SystemStats{
		CPUUsage:      cpuUsage,
		CPUCores:      cores,
		MemoryUsed:    memInfo.Used,
		MemoryTotal:   memInfo.Total,
		MemoryPercent: memInfo.UsedPercent,
		DiskUsed:      diskInfo.Used,
		DiskTotal:     diskInfo.Total,
		DiskPercent:   diskInfo.UsedPercent,
		Uptime:        hostInfo.Uptime,
		Hostname:      hostInfo.Hostname,
		OS:            hostInfo.OS,
		Platform:      hostInfo.Platform,
	}, nil
}
