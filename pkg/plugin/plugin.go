package plugin

import "sync"

// Registry is the interface used by plugins to register their capabilities.
type Registry interface {
	RegisterDetector(name string, detect func([]byte) float64)
	RegisterDecoder(name string, decode func([]byte) ([]byte, error))
	RegisterEncoder(name string, encode func([]byte) ([]byte, error))
}

// Plugin is the interface that all ForgeCLI plugins must implement.
type Plugin interface {
	Name() string
	Version() string
	Register(reg Registry)
}

// Manager manages loaded plugins.
type Manager struct {
	mu      sync.RWMutex
	plugins []Plugin
}

// NewManager creates a new empty plugin Manager.
func NewManager() *Manager {
	return &Manager{}
}

// Load registers a plugin with the manager.
func (m *Manager) Load(p Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.plugins = append(m.plugins, p)
}

// List returns all loaded plugins.
func (m *Manager) List() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]Plugin, len(m.plugins))
	copy(result, m.plugins)
	return result
}

// internalRegistry is the default implementation of Registry used internally.
type internalRegistry struct {
	detectors map[string]func([]byte) float64
	decoders  map[string]func([]byte) ([]byte, error)
	encoders  map[string]func([]byte) ([]byte, error)
}

func newInternalRegistry() *internalRegistry {
	return &internalRegistry{
		detectors: make(map[string]func([]byte) float64),
		decoders:  make(map[string]func([]byte) ([]byte, error)),
		encoders:  make(map[string]func([]byte) ([]byte, error)),
	}
}

func (r *internalRegistry) RegisterDetector(name string, detect func([]byte) float64) {
	r.detectors[name] = detect
}

func (r *internalRegistry) RegisterDecoder(name string, decode func([]byte) ([]byte, error)) {
	r.decoders[name] = decode
}

func (r *internalRegistry) RegisterEncoder(name string, encode func([]byte) ([]byte, error)) {
	r.encoders[name] = encode
}
