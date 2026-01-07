package stack

import (
	"context"
	"os"
	"path/filepath"

	"aperture-science-network/internal/docker"
)

// FilesystemProvider implements Provider using the local filesystem
type FilesystemProvider struct {
	stacksPath   string
	dockerClient docker.DockerClient
}

// NewFilesystemProvider creates a new filesystem-based stack provider
func NewFilesystemProvider(stacksPath string, dockerClient docker.DockerClient) *FilesystemProvider {
	return &FilesystemProvider{
		stacksPath:   stacksPath,
		dockerClient: dockerClient,
	}
}

func (p *FilesystemProvider) ListStacks() ([]StackInfo, error) {
	entries, err := os.ReadDir(p.stacksPath)
	if err != nil {
		return nil, err
	}

	// Get all containers to determine stack status
	containers, err := p.dockerClient.ListContainers(context.Background(), true)
	if err != nil {
		return nil, err
	}

	// Group containers by compose project
	containersByProject := make(map[string][]docker.ContainerInfo)
	for _, ctr := range containers {
		if project, ok := ctr.Labels["com.docker.compose.project"]; ok {
			containersByProject[project] = append(containersByProject[project], ctr)
		}
	}

	stacks := make([]StackInfo, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			composePath := filepath.Join(p.stacksPath, entry.Name(), "docker-compose.yml")
			if _, err := os.Stat(composePath); err == nil {
				stackName := entry.Name()
				projectContainers := containersByProject[stackName]

				// Count running services
				runningCount := 0
				for _, ctr := range projectContainers {
					if ctr.State == "running" {
						runningCount++
					}
				}

				// Determine stack status
				status := "stopped"
				totalServices := len(projectContainers)
				if totalServices > 0 {
					if runningCount == totalServices {
						status = "running"
					} else if runningCount > 0 {
						status = "partial"
					}
				}

				stacks = append(stacks, StackInfo{
					Name:            stackName,
					Path:            filepath.Join(p.stacksPath, stackName),
					Status:          status,
					Services:        totalServices,
					RunningServices: runningCount,
				})
			}
		}
	}

	return stacks, nil
}

func (p *FilesystemProvider) GetStack(name string) (*StackInfo, error) {
	stackPath := filepath.Join(p.stacksPath, name)

	if _, err := os.Stat(stackPath); os.IsNotExist(err) {
		return nil, err
	}

	return &StackInfo{
		Name: name,
		Path: stackPath,
	}, nil
}

func (p *FilesystemProvider) GetComposeFile(name string) (string, error) {
	composePath := filepath.Join(p.stacksPath, name, "docker-compose.yml")
	content, err := os.ReadFile(composePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (p *FilesystemProvider) UpdateComposeFile(name string, content string) error {
	composePath := filepath.Join(p.stacksPath, name, "docker-compose.yml")
	return os.WriteFile(composePath, []byte(content), 0644)
}

func (p *FilesystemProvider) StackExists(name string) bool {
	stackPath := filepath.Join(p.stacksPath, name)
	_, err := os.Stat(stackPath)
	return err == nil
}

func (p *FilesystemProvider) GetStackPath(name string) string {
	return filepath.Join(p.stacksPath, name)
}
