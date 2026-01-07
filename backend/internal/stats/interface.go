package stats

// Provider defines the interface for system stats collection
type Provider interface {
	GetSystemStats() (*SystemStats, error)
}

// DefaultProvider is the default implementation using gopsutil
type DefaultProvider struct{}

// NewDefaultProvider creates a new default stats provider
func NewDefaultProvider() *DefaultProvider {
	return &DefaultProvider{}
}

// GetSystemStats returns the current system statistics
func (p *DefaultProvider) GetSystemStats() (*SystemStats, error) {
	return GetSystemStats()
}

// Ensure DefaultProvider implements Provider
var _ Provider = (*DefaultProvider)(nil)
