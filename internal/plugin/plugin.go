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

// Manager manages loaded plugins and aggregates their registered capabilities.
type Manager struct {
	mu        sync.RWMutex
	plugins   []Plugin
	detectors map[string]func([]byte) float64
	decoders  map[string]func([]byte) ([]byte, error)
	encoders  map[string]func([]byte) ([]byte, error)
}

// NewManager creates a new empty plugin Manager.
func NewManager() *Manager {
	return &Manager{
		detectors: make(map[string]func([]byte) float64),
		decoders:  make(map[string]func([]byte) ([]byte, error)),
		encoders:  make(map[string]func([]byte) ([]byte, error)),
	}
}

// Load registers a plugin, calls its Register method to collect capabilities,
// and merges those capabilities into the manager's shared maps.
func (m *Manager) Load(p Plugin) {
	reg := newInternalRegistry()
	p.Register(reg)

	m.mu.Lock()
	defer m.mu.Unlock()
	m.plugins = append(m.plugins, p)
	for k, v := range reg.detectors {
		m.detectors[k] = v
	}
	for k, v := range reg.decoders {
		m.decoders[k] = v
	}
	for k, v := range reg.encoders {
		m.encoders[k] = v
	}
}

// List returns all loaded plugins.
func (m *Manager) List() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]Plugin, len(m.plugins))
	copy(result, m.plugins)
	return result
}

// Detectors returns a copy of all detector functions registered by plugins.
func (m *Manager) Detectors() map[string]func([]byte) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]func([]byte) float64, len(m.detectors))
	for k, v := range m.detectors {
		out[k] = v
	}
	return out
}

// Decoders returns a copy of all decoder functions registered by plugins.
func (m *Manager) Decoders() map[string]func([]byte) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]func([]byte) ([]byte, error), len(m.decoders))
	for k, v := range m.decoders {
		out[k] = v
	}
	return out
}

// Encoders returns a copy of all encoder functions registered by plugins.
func (m *Manager) Encoders() map[string]func([]byte) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]func([]byte) ([]byte, error), len(m.encoders))
	for k, v := range m.encoders {
		out[k] = v
	}
	return out
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
