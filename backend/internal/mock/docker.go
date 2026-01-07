package mock

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"aperture-science-network/internal/docker"
)

// DockerClient is a mock implementation of docker.DockerClient
type DockerClient struct {
	mu        sync.Mutex
	startTime time.Time
}

// NewDockerClient creates a new mock Docker client
func NewDockerClient() *DockerClient {
	return &DockerClient{
		startTime: time.Now(),
	}
}

// Ensure MockDockerClient implements docker.DockerClient
var _ docker.DockerClient = (*DockerClient)(nil)

func (c *DockerClient) ListContainers(ctx context.Context, all bool) ([]docker.ContainerInfo, error) {
	containers := []docker.ContainerInfo{
		{
			ID:      "a1b2c3d4e5f6",
			Name:    "celeste-frontend",
			Image:   "celeste:latest",
			Status:  "Up 2 hours",
			State:   "running",
			Created: time.Now().Add(-2 * time.Hour).Unix(),
			Ports: []docker.PortBinding{
				{Private: 80, Public: 8080, Type: "tcp"},
			},
			Labels: map[string]string{
				"com.docker.compose.project": "celeste",
				"com.docker.compose.service": "frontend",
			},
			NetworkMode: "celeste_default",
		},
		{
			ID:      "b2c3d4e5f6a7",
			Name:    "celeste-backend",
			Image:   "celeste:latest",
			Status:  "Up 2 hours",
			State:   "running",
			Created: time.Now().Add(-2 * time.Hour).Unix(),
			Ports: []docker.PortBinding{
				{Private: 8080, Public: 8081, Type: "tcp"},
			},
			Labels: map[string]string{
				"com.docker.compose.project": "celeste",
				"com.docker.compose.service": "backend",
			},
			NetworkMode: "celeste_default",
		},
		{
			ID:      "c3d4e5f6a7b8",
			Name:    "prometheus",
			Image:   "prom/prometheus:latest",
			Status:  "Up 5 hours",
			State:   "running",
			Created: time.Now().Add(-5 * time.Hour).Unix(),
			Ports: []docker.PortBinding{
				{Private: 9090, Public: 9090, Type: "tcp"},
			},
			Labels: map[string]string{
				"com.docker.compose.project": "monitoring",
				"com.docker.compose.service": "prometheus",
			},
			NetworkMode: "monitoring_default",
		},
		{
			ID:      "d4e5f6a7b8c9",
			Name:    "grafana",
			Image:   "grafana/grafana:latest",
			Status:  "Exited (0) 1 hour ago",
			State:   "exited",
			Created: time.Now().Add(-6 * time.Hour).Unix(),
			Ports:   []docker.PortBinding{},
			Labels: map[string]string{
				"com.docker.compose.project": "monitoring",
				"com.docker.compose.service": "grafana",
			},
			NetworkMode: "monitoring_default",
		},
		{
			ID:      "e5f6a7b8c9d0",
			Name:    "redis",
			Image:   "redis:alpine",
			Status:  "Exited (0) 30 minutes ago",
			State:   "exited",
			Created: time.Now().Add(-3 * time.Hour).Unix(),
			Ports:   []docker.PortBinding{},
			Labels: map[string]string{
				"com.docker.compose.project": "celeste",
				"com.docker.compose.service": "redis",
			},
			NetworkMode: "celeste_default",
		},
	}

	if !all {
		// Filter to only running containers
		running := make([]docker.ContainerInfo, 0)
		for _, c := range containers {
			if c.State == "running" {
				running = append(running, c)
			}
		}
		return running, nil
	}

	return containers, nil
}

func (c *DockerClient) GetContainer(ctx context.Context, id string) (*docker.ContainerInfo, error) {
	containers, _ := c.ListContainers(ctx, true)
	for _, ctr := range containers {
		if ctr.ID == id || ctr.Name == id {
			return &ctr, nil
		}
	}
	return &containers[0], nil
}

func (c *DockerClient) StartContainer(ctx context.Context, id string) error {
	return nil
}

func (c *DockerClient) StopContainer(ctx context.Context, id string) error {
	return nil
}

func (c *DockerClient) RestartContainer(ctx context.Context, id string) error {
	return nil
}

func (c *DockerClient) GetContainerLogs(ctx context.Context, id string, tail string) (string, error) {
	return `2024-01-07T10:00:00.000Z [INFO] Container started
2024-01-07T10:00:01.000Z [INFO] Listening on port 8080
2024-01-07T10:00:02.000Z [DEBUG] Health check passed
2024-01-07T10:00:05.000Z [INFO] Connection established
2024-01-07T10:00:10.000Z [DEBUG] Processing request
`, nil
}

func (c *DockerClient) GetContainerStats(ctx context.Context, id string) (*docker.ContainerStats, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Generate slightly varying stats
	baseStats := map[string]struct {
		cpu    float64
		memory uint64
		limit  uint64
	}{
		"a1b2c3d4e5f6": {5.0, 128 * 1024 * 1024, 512 * 1024 * 1024},   // celeste-frontend
		"b2c3d4e5f6a7": {10.0, 256 * 1024 * 1024, 1024 * 1024 * 1024}, // celeste-backend
		"c3d4e5f6a7b8": {7.0, 512 * 1024 * 1024, 2048 * 1024 * 1024},  // prometheus
	}

	base, ok := baseStats[id]
	if !ok {
		base = baseStats["a1b2c3d4e5f6"]
	}

	// Add some variation
	cpuVariation := (rand.Float64() - 0.5) * 6 // +/- 3%
	memVariation := uint64((rand.Float64() - 0.5) * float64(base.memory) * 0.2)

	cpu := base.cpu + cpuVariation
	if cpu < 0 {
		cpu = 0
	}

	memory := base.memory + memVariation
	memPercent := float64(memory) / float64(base.limit) * 100

	return &docker.ContainerStats{
		CPUPercent:    cpu,
		MemoryUsage:   memory,
		MemoryLimit:   base.limit,
		MemoryPercent: memPercent,
		NetworkRx:     uint64(rand.Intn(1000000)),
		NetworkTx:     uint64(rand.Intn(500000)),
		BlockRead:     uint64(rand.Intn(10000000)),
		BlockWrite:    uint64(rand.Intn(5000000)),
	}, nil
}

func (c *DockerClient) ListVolumes(ctx context.Context) ([]docker.VolumeInfo, error) {
	return []docker.VolumeInfo{
		{
			Name:       "celeste_data",
			Driver:     "local",
			Mountpoint: "/var/lib/docker/volumes/celeste_data/_data",
			CreatedAt:  time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			Labels:     map[string]string{"project": "celeste"},
			UsedBy:     []string{"celeste-backend"},
		},
		{
			Name:       "prometheus_data",
			Driver:     "local",
			Mountpoint: "/var/lib/docker/volumes/prometheus_data/_data",
			CreatedAt:  time.Now().Add(-48 * time.Hour).Format(time.RFC3339),
			Labels:     map[string]string{"project": "monitoring"},
			UsedBy:     []string{"prometheus"},
		},
		{
			Name:       "grafana_data",
			Driver:     "local",
			Mountpoint: "/var/lib/docker/volumes/grafana_data/_data",
			CreatedAt:  time.Now().Add(-48 * time.Hour).Format(time.RFC3339),
			Labels:     map[string]string{"project": "monitoring"},
			UsedBy:     []string{"grafana"},
		},
	}, nil
}

func (c *DockerClient) CreateVolume(ctx context.Context, name string, driver string, labels map[string]string) (*docker.VolumeInfo, error) {
	return &docker.VolumeInfo{
		Name:       name,
		Driver:     driver,
		Mountpoint: "/var/lib/docker/volumes/" + name + "/_data",
		CreatedAt:  time.Now().Format(time.RFC3339),
		Labels:     labels,
	}, nil
}

func (c *DockerClient) DeleteVolume(ctx context.Context, name string, force bool) error {
	return nil
}

func (c *DockerClient) ListNetworks(ctx context.Context) ([]docker.NetworkInfo, error) {
	return []docker.NetworkInfo{
		{
			ID:         "net1a2b3c4d5e",
			Name:       "bridge",
			Driver:     "bridge",
			Scope:      "local",
			Internal:   false,
			Containers: []string{},
		},
		{
			ID:         "net2b3c4d5e6f",
			Name:       "celeste_default",
			Driver:     "bridge",
			Scope:      "local",
			Internal:   false,
			Containers: []string{"celeste-frontend", "celeste-backend", "redis"},
		},
		{
			ID:         "net3c4d5e6f7a",
			Name:       "monitoring_default",
			Driver:     "bridge",
			Scope:      "local",
			Internal:   false,
			Containers: []string{"prometheus", "grafana"},
		},
	}, nil
}

func (c *DockerClient) CreateNetwork(ctx context.Context, name string, driver string) (*docker.NetworkInfo, error) {
	return &docker.NetworkInfo{
		ID:     "netnew12345",
		Name:   name,
		Driver: driver,
		Scope:  "local",
	}, nil
}

func (c *DockerClient) DeleteNetwork(ctx context.Context, id string) error {
	return nil
}

func (c *DockerClient) ListImages(ctx context.Context) ([]docker.ImageInfo, error) {
	return []docker.ImageInfo{
		{
			ID:      "sha256:abc123",
			Tags:    []string{"celeste:latest", "celeste:2.1.1"},
			Size:    150 * 1024 * 1024,
			Created: time.Now().Add(-24 * time.Hour).Unix(),
		},
		{
			ID:      "sha256:def456",
			Tags:    []string{"prom/prometheus:latest"},
			Size:    250 * 1024 * 1024,
			Created: time.Now().Add(-72 * time.Hour).Unix(),
		},
		{
			ID:      "sha256:ghi789",
			Tags:    []string{"grafana/grafana:latest"},
			Size:    400 * 1024 * 1024,
			Created: time.Now().Add(-72 * time.Hour).Unix(),
		},
		{
			ID:      "sha256:jkl012",
			Tags:    []string{"redis:alpine"},
			Size:    30 * 1024 * 1024,
			Created: time.Now().Add(-168 * time.Hour).Unix(),
		},
	}, nil
}

func (c *DockerClient) Close() error {
	return nil
}
