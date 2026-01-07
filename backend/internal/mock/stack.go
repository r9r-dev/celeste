package mock

import (
	"fmt"

	"aperture-science-network/internal/stack"
)

// StackProvider is a mock implementation of stack.Provider
type StackProvider struct {
	stacks map[string]*stack.StackInfo
}

// NewStackProvider creates a new mock stack provider
func NewStackProvider() *StackProvider {
	return &StackProvider{
		stacks: map[string]*stack.StackInfo{
			"celeste": {
				Name:            "celeste",
				Path:            "/stacks/celeste",
				Status:          "partial",
				Services:        3,
				RunningServices: 2,
			},
			"monitoring": {
				Name:            "monitoring",
				Path:            "/stacks/monitoring",
				Status:          "partial",
				Services:        2,
				RunningServices: 1,
			},
			"database": {
				Name:            "database",
				Path:            "/stacks/database",
				Status:          "stopped",
				Services:        2,
				RunningServices: 0,
			},
		},
	}
}

// Ensure StackProvider implements stack.Provider
var _ stack.Provider = (*StackProvider)(nil)

func (p *StackProvider) ListStacks() ([]stack.StackInfo, error) {
	stacks := make([]stack.StackInfo, 0, len(p.stacks))
	for _, s := range p.stacks {
		stacks = append(stacks, *s)
	}
	return stacks, nil
}

func (p *StackProvider) GetStack(name string) (*stack.StackInfo, error) {
	if s, ok := p.stacks[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("stack not found: %s", name)
}

func (p *StackProvider) GetComposeFile(name string) (string, error) {
	if _, ok := p.stacks[name]; !ok {
		return "", fmt.Errorf("stack not found: %s", name)
	}

	// Return mock compose content based on stack name
	switch name {
	case "celeste":
		return `version: "3.8"

services:
  frontend:
    image: celeste:latest
    ports:
      - "8080:80"
    depends_on:
      - backend

  backend:
    image: celeste:latest
    ports:
      - "8081:8080"
    depends_on:
      - redis

  redis:
    image: redis:alpine
    volumes:
      - redis_data:/data

volumes:
  redis_data:
`, nil

	case "monitoring":
		return `version: "3.8"

services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - prometheus_data:/prometheus

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  prometheus_data:
  grafana_data:
`, nil

	case "database":
		return `version: "3.8"

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_PASSWORD: secret
    volumes:
      - pg_data:/var/lib/postgresql/data

  adminer:
    image: adminer
    ports:
      - "8082:8080"

volumes:
  pg_data:
`, nil

	default:
		return `version: "3.8"

services:
  app:
    image: nginx:alpine
    ports:
      - "80:80"
`, nil
	}
}

func (p *StackProvider) UpdateComposeFile(name string, content string) error {
	if _, ok := p.stacks[name]; !ok {
		return fmt.Errorf("stack not found: %s", name)
	}
	// In mock mode, just pretend to save
	return nil
}

func (p *StackProvider) StackExists(name string) bool {
	_, ok := p.stacks[name]
	return ok
}

func (p *StackProvider) GetStackPath(name string) string {
	if s, ok := p.stacks[name]; ok {
		return s.Path
	}
	return "/stacks/" + name
}
