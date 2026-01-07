package docker

import "context"

// DockerClient defines the interface for Docker operations
type DockerClient interface {
	ListContainers(ctx context.Context, all bool) ([]ContainerInfo, error)
	GetContainer(ctx context.Context, id string) (*ContainerInfo, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	RestartContainer(ctx context.Context, id string) error
	GetContainerLogs(ctx context.Context, id string, tail string) (string, error)
	GetContainerStats(ctx context.Context, id string) (*ContainerStats, error)
	ListVolumes(ctx context.Context) ([]VolumeInfo, error)
	CreateVolume(ctx context.Context, name string, driver string, labels map[string]string) (*VolumeInfo, error)
	DeleteVolume(ctx context.Context, name string, force bool) error
	ListNetworks(ctx context.Context) ([]NetworkInfo, error)
	CreateNetwork(ctx context.Context, name string, driver string) (*NetworkInfo, error)
	DeleteNetwork(ctx context.Context, id string) error
	ListImages(ctx context.Context) ([]ImageInfo, error)
	Close() error
}

// Ensure Client implements DockerClient
var _ DockerClient = (*Client)(nil)
