package docker

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

type ContainerInfo struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Status       string            `json:"status"`
	State        string            `json:"state"`
	Created      int64             `json:"created"`
	Ports        []PortBinding     `json:"ports"`
	Labels       map[string]string `json:"labels"`
	NetworkMode  string            `json:"networkMode"`
}

type PortBinding struct {
	Private int    `json:"private"`
	Public  int    `json:"public"`
	Type    string `json:"type"`
}

type ContainerStats struct {
	CPUPercent    float64 `json:"cpuPercent"`
	MemoryUsage   uint64  `json:"memoryUsage"`
	MemoryLimit   uint64  `json:"memoryLimit"`
	MemoryPercent float64 `json:"memoryPercent"`
	NetworkRx     uint64  `json:"networkRx"`
	NetworkTx     uint64  `json:"networkTx"`
	BlockRead     uint64  `json:"blockRead"`
	BlockWrite    uint64  `json:"blockWrite"`
}

type VolumeInfo struct {
	Name       string   `json:"name"`
	Driver     string   `json:"driver"`
	Mountpoint string   `json:"mountpoint"`
	CreatedAt  string   `json:"createdAt"`
	Labels     map[string]string `json:"labels"`
	UsedBy     []string `json:"usedBy"`
}

type NetworkInfo struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Driver     string   `json:"driver"`
	Scope      string   `json:"scope"`
	Internal   bool     `json:"internal"`
	Containers []string `json:"containers"`
}

type ImageInfo struct {
	ID      string   `json:"id"`
	Tags    []string `json:"tags"`
	Size    int64    `json:"size"`
	Created int64    `json:"created"`
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) ListContainers(ctx context.Context, all bool) ([]ContainerInfo, error) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: all})
	if err != nil {
		return nil, err
	}

	result := make([]ContainerInfo, len(containers))
	for i, ctr := range containers {
		ports := make([]PortBinding, 0)
		for _, p := range ctr.Ports {
			ports = append(ports, PortBinding{
				Private: int(p.PrivatePort),
				Public:  int(p.PublicPort),
				Type:    p.Type,
			})
		}

		name := ""
		if len(ctr.Names) > 0 {
			name = strings.TrimPrefix(ctr.Names[0], "/")
		}

		result[i] = ContainerInfo{
			ID:          ctr.ID[:12],
			Name:        name,
			Image:       ctr.Image,
			Status:      ctr.Status,
			State:       ctr.State,
			Created:     ctr.Created,
			Ports:       ports,
			Labels:      ctr.Labels,
			NetworkMode: ctr.HostConfig.NetworkMode,
		}
	}

	return result, nil
}

func (c *Client) GetContainer(ctx context.Context, id string) (*ContainerInfo, error) {
	ctr, err := c.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}

	ports := make([]PortBinding, 0)
	for port, bindings := range ctr.NetworkSettings.Ports {
		for _, b := range bindings {
			var publicPort int
			if b.HostPort != "" {
				// Parse port string to int
				publicPort = 0 // simplified
			}
			ports = append(ports, PortBinding{
				Private: port.Int(),
				Public:  publicPort,
				Type:    port.Proto(),
			})
		}
	}

	var createdUnix int64
	if createdTime, err := time.Parse(time.RFC3339Nano, ctr.Created); err == nil {
		createdUnix = createdTime.Unix()
	}

	return &ContainerInfo{
		ID:          ctr.ID[:12],
		Name:        strings.TrimPrefix(ctr.Name, "/"),
		Image:       ctr.Config.Image,
		Status:      ctr.State.Status,
		State:       ctr.State.Status,
		Created:     createdUnix,
		Ports:       ports,
		Labels:      ctr.Config.Labels,
		NetworkMode: string(ctr.HostConfig.NetworkMode),
	}, nil
}

func (c *Client) StartContainer(ctx context.Context, id string) error {
	return c.cli.ContainerStart(ctx, id, container.StartOptions{})
}

func (c *Client) StopContainer(ctx context.Context, id string) error {
	timeout := 10
	return c.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout})
}

func (c *Client) RestartContainer(ctx context.Context, id string) error {
	timeout := 10
	return c.cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout})
}

func (c *Client) GetContainerLogs(ctx context.Context, id string, tail string) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
		Timestamps: true,
	}

	reader, err := c.cli.ContainerLogs(ctx, id, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(logs), nil
}

func (c *Client) GetContainerStats(ctx context.Context, id string) (*ContainerStats, error) {
	stats, err := c.cli.ContainerStats(ctx, id, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	var v types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
		return nil, err
	}

	// Calculate CPU percentage
	cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(v.CPUStats.SystemUsage - v.PreCPUStats.SystemUsage)
	cpuPercent := 0.0
	if systemDelta > 0 && cpuDelta > 0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	// Calculate memory percentage
	memPercent := 0.0
	if v.MemoryStats.Limit > 0 {
		memPercent = float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100.0
	}

	// Calculate network I/O
	var rxBytes, txBytes uint64
	for _, net := range v.Networks {
		rxBytes += net.RxBytes
		txBytes += net.TxBytes
	}

	return &ContainerStats{
		CPUPercent:    cpuPercent,
		MemoryUsage:   v.MemoryStats.Usage,
		MemoryLimit:   v.MemoryStats.Limit,
		MemoryPercent: memPercent,
		NetworkRx:     rxBytes,
		NetworkTx:     txBytes,
	}, nil
}

func (c *Client) ListVolumes(ctx context.Context) ([]VolumeInfo, error) {
	volumes, err := c.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]VolumeInfo, len(volumes.Volumes))
	for i, vol := range volumes.Volumes {
		result[i] = VolumeInfo{
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
			CreatedAt:  vol.CreatedAt,
			Labels:     vol.Labels,
		}
	}

	return result, nil
}

func (c *Client) CreateVolume(ctx context.Context, name string, driver string, labels map[string]string) (*VolumeInfo, error) {
	vol, err := c.cli.VolumeCreate(ctx, volume.CreateOptions{
		Name:   name,
		Driver: driver,
		Labels: labels,
	})
	if err != nil {
		return nil, err
	}

	return &VolumeInfo{
		Name:       vol.Name,
		Driver:     vol.Driver,
		Mountpoint: vol.Mountpoint,
		Labels:     vol.Labels,
	}, nil
}

func (c *Client) DeleteVolume(ctx context.Context, name string, force bool) error {
	return c.cli.VolumeRemove(ctx, name, force)
}

func (c *Client) ListNetworks(ctx context.Context) ([]NetworkInfo, error) {
	networks, err := c.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]NetworkInfo, len(networks))
	for i, net := range networks {
		containers := make([]string, 0, len(net.Containers))
		for _, ctr := range net.Containers {
			containers = append(containers, ctr.Name)
		}

		result[i] = NetworkInfo{
			ID:         net.ID[:12],
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			Internal:   net.Internal,
			Containers: containers,
		}
	}

	return result, nil
}

func (c *Client) CreateNetwork(ctx context.Context, name string, driver string) (*NetworkInfo, error) {
	resp, err := c.cli.NetworkCreate(ctx, name, network.CreateOptions{
		Driver: driver,
	})
	if err != nil {
		return nil, err
	}

	return &NetworkInfo{
		ID:     resp.ID[:12],
		Name:   name,
		Driver: driver,
	}, nil
}

func (c *Client) DeleteNetwork(ctx context.Context, id string) error {
	return c.cli.NetworkRemove(ctx, id)
}

func (c *Client) ListImages(ctx context.Context) ([]ImageInfo, error) {
	images, err := c.cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]ImageInfo, len(images))
	for i, img := range images {
		result[i] = ImageInfo{
			ID:      img.ID[7:19], // Remove "sha256:" prefix and truncate
			Tags:    img.RepoTags,
			Size:    img.Size,
			Created: img.Created,
		}
	}

	return result, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}
