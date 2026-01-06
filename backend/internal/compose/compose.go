package compose

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Manager handles docker-compose operations via CLI
type Manager struct {
	dockerCmd string
}

// ServiceStatus represents the status of a compose service
type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Health string `json:"health,omitempty"`
}

// NewManager creates a new compose manager, detecting the appropriate CLI command
func NewManager() *Manager {
	// Check if 'container' CLI is available (macOS)
	if _, err := exec.LookPath("container"); err == nil {
		return &Manager{dockerCmd: "container"}
	}
	return &Manager{dockerCmd: "docker"}
}

// Up starts all services in a compose stack
func (m *Manager) Up(ctx context.Context, stackPath string) error {
	composeFile := filepath.Join(stackPath, "docker-compose.yml")
	cmd := exec.CommandContext(ctx, m.dockerCmd, "compose", "-f", composeFile, "up", "-d")
	cmd.Dir = stackPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compose up failed: %s: %w", stderr.String(), err)
	}
	return nil
}

// Down stops and removes all services in a compose stack
func (m *Manager) Down(ctx context.Context, stackPath string) error {
	composeFile := filepath.Join(stackPath, "docker-compose.yml")
	cmd := exec.CommandContext(ctx, m.dockerCmd, "compose", "-f", composeFile, "down")
	cmd.Dir = stackPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compose down failed: %s: %w", stderr.String(), err)
	}
	return nil
}

// Restart restarts all services in a compose stack
func (m *Manager) Restart(ctx context.Context, stackPath string) error {
	composeFile := filepath.Join(stackPath, "docker-compose.yml")
	cmd := exec.CommandContext(ctx, m.dockerCmd, "compose", "-f", composeFile, "restart")
	cmd.Dir = stackPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compose restart failed: %s: %w", stderr.String(), err)
	}
	return nil
}

// Pull pulls the latest images for all services in a compose stack
func (m *Manager) Pull(ctx context.Context, stackPath string) error {
	composeFile := filepath.Join(stackPath, "docker-compose.yml")
	cmd := exec.CommandContext(ctx, m.dockerCmd, "compose", "-f", composeFile, "pull")
	cmd.Dir = stackPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compose pull failed: %s: %w", stderr.String(), err)
	}
	return nil
}

// Logs retrieves logs for a specific service or all services
func (m *Manager) Logs(ctx context.Context, stackPath string, service string, tail int) (string, error) {
	composeFile := filepath.Join(stackPath, "docker-compose.yml")
	args := []string{"compose", "-f", composeFile, "logs", "--no-color", fmt.Sprintf("--tail=%d", tail)}
	if service != "" {
		args = append(args, service)
	}

	cmd := exec.CommandContext(ctx, m.dockerCmd, args...)
	cmd.Dir = stackPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("compose logs failed: %s: %w", stderr.String(), err)
	}
	return stdout.String(), nil
}

// PS returns the status of services in a compose stack
func (m *Manager) PS(ctx context.Context, stackPath string) ([]ServiceStatus, error) {
	composeFile := filepath.Join(stackPath, "docker-compose.yml")
	cmd := exec.CommandContext(ctx, m.dockerCmd, "compose", "-f", composeFile, "ps", "--format", "{{.Service}}|{{.State}}|{{.Health}}")
	cmd.Dir = stackPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// If no containers are running, ps might fail - return empty list
		if strings.Contains(stderr.String(), "no such service") || strings.Contains(stderr.String(), "no configuration") {
			return []ServiceStatus{}, nil
		}
		return nil, fmt.Errorf("compose ps failed: %s: %w", stderr.String(), err)
	}

	var services []ServiceStatus
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 2 {
			status := ServiceStatus{
				Name:   parts[0],
				Status: parts[1],
			}
			if len(parts) >= 3 {
				status.Health = parts[2]
			}
			services = append(services, status)
		}
	}
	return services, nil
}
