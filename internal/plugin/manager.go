package plugin

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/huianlei/antia-aitool-mcp/pkg/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Manager manages the lifecycle of all plugins
type Manager struct {
	plugins map[string]Plugin
	logger  *zap.Logger
	mu      sync.RWMutex
}

// NewManager creates a new plugin manager
func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
		logger:  logger,
	}
}

// LoadPlugin loads a plugin by name with the given configuration
func (m *Manager) LoadPlugin(name string, config Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if plugin is already loaded
	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin %s is already loaded", name)
	}

	// Get plugin factory
	factory, exists := GetFactory(name)
	if !exists {
		return errors.Wrapf(models.ErrPluginNotFound, "plugin: %s", name)
	}

	// Create plugin instance
	plugin, err := factory(config)
	if err != nil {
		return errors.Wrapf(models.ErrPluginInitFailed, "plugin: %s, error: %v", name, err)
	}

	// Initialize plugin
	if err := plugin.Initialize(config); err != nil {
		return errors.Wrapf(err, "failed to initialize plugin: %s", name)
	}

	// Start plugin
	if err := plugin.Start(); err != nil {
		return errors.Wrapf(err, "failed to start plugin: %s", name)
	}

	m.plugins[name] = plugin
	m.logger.Info("plugin loaded", zap.String("plugin", name), zap.String("version", plugin.Version()))

	return nil
}

// UnloadPlugin unloads a plugin by name
func (m *Manager) UnloadPlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return errors.Wrapf(models.ErrPluginNotFound, "plugin: %s", name)
	}

	if err := plugin.Stop(); err != nil {
		m.logger.Warn("error stopping plugin", zap.String("plugin", name), zap.Error(err))
	}

	delete(m.plugins, name)
	m.logger.Info("plugin unloaded", zap.String("plugin", name))

	return nil
}

// GetPlugin retrieves a loaded plugin by name
func (m *Manager) GetPlugin(name string) (Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return nil, errors.Wrapf(models.ErrPluginNotFound, "plugin: %s", name)
	}

	return plugin, nil
}

// GetAllTools returns all tools from all loaded plugins
func (m *Manager) GetAllTools() []models.Tool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var allTools []models.Tool
	for _, plugin := range m.plugins {
		tools := plugin.GetTools()
		allTools = append(allTools, tools...)
	}

	return allTools
}

// ExecuteTool routes a tool call to the appropriate plugin
func (m *Manager) ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (interface{}, error) {
	// Parse tool name: {plugin_name}_{tool_action}
	parts := strings.SplitN(toolName, "_", 2)
	if len(parts) < 2 {
		return nil, errors.Wrapf(models.ErrInvalidToolName, "tool: %s", toolName)
	}

	pluginName := parts[0]

	// Get plugin
	plugin, err := m.GetPlugin(pluginName)
	if err != nil {
		return nil, err
	}

	// Execute tool
	m.logger.Debug("executing tool",
		zap.String("plugin", pluginName),
		zap.String("tool", toolName),
	)

	result, err := plugin.ExecuteTool(ctx, toolName, params)
	if err != nil {
		return nil, errors.Wrapf(err, "tool execution failed: %s", toolName)
	}

	return result, nil
}

// HealthCheck checks the health of all plugins
func (m *Manager) HealthCheck() map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make(map[string]error)
	for name, plugin := range m.plugins {
		results[name] = plugin.HealthCheck()
	}

	return results
}

// Shutdown stops all plugins
func (m *Manager) Shutdown() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for name, plugin := range m.plugins {
		if err := plugin.Stop(); err != nil {
			errs = append(errs, errors.Wrapf(err, "plugin: %s", name))
		}
	}

	m.plugins = make(map[string]Plugin)

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}

	return nil
}

// ListLoadedPlugins returns the names of all loaded plugins
func (m *Manager) ListLoadedPlugins() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.plugins))
	for name := range m.plugins {
		names = append(names, name)
	}

	return names
}
