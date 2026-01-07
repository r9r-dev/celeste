package stack

// StackInfo represents a compose stack with its status
type StackInfo struct {
	Name            string `json:"name"`
	Path            string `json:"path"`
	Status          string `json:"status"`
	Services        int    `json:"services"`
	RunningServices int    `json:"runningServices"`
}

// Provider defines the interface for stack operations
type Provider interface {
	// ListStacks returns all available stacks
	ListStacks() ([]StackInfo, error)

	// GetStack returns a specific stack by name
	GetStack(name string) (*StackInfo, error)

	// GetComposeFile returns the content of a stack's docker-compose.yml
	GetComposeFile(name string) (string, error)

	// UpdateComposeFile updates the content of a stack's docker-compose.yml
	UpdateComposeFile(name string, content string) error

	// StackExists checks if a stack exists
	StackExists(name string) bool

	// GetStackPath returns the filesystem path for a stack
	GetStackPath(name string) string
}
