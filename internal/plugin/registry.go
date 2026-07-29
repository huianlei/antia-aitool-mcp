package plugin

import (
	"sync"
)

var (
	// Global plugin registry
	defaultRegistry = &Registry{
		factories: make(map[string]Factory),
	}
)

// Registry manages plugin factories
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// Register registers a plugin factory with the given name
func Register(name string, factory Factory) {
	defaultRegistry.Register(name, factory)
}

// Register registers a plugin factory
func (r *Registry) Register(name string, factory Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = factory
}

// Get retrieves a plugin factory by name
func (r *Registry) Get(name string) (Factory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	factory, exists := r.factories[name]
	return factory, exists
}

// List returns all registered plugin names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}

// GetFactory retrieves a plugin factory by name from the default registry
func GetFactory(name string) (Factory, bool) {
	return defaultRegistry.Get(name)
}

// ListPlugins returns all registered plugin names from the default registry
func ListPlugins() []string {
	return defaultRegistry.List()
}
