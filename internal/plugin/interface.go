package plugin

import (
	"context"

	"github.com/huianlei/antia-aitool-mcp/pkg/models"
)

// Plugin is the interface that all plugins must implement
type Plugin interface {
	// Metadata
	Name() string
	Version() string
	Description() string

	// Lifecycle
	Initialize(config Config) error
	Start() error
	Stop() error
	HealthCheck() error

	// Tool management
	GetTools() []models.Tool
	ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error)
}

// Config is the interface for plugin configuration
type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetMap(key string) map[string]interface{}
	GetStringSlice(key string) []string
}

// Factory is a function that creates a plugin instance
type Factory func(config Config) (Plugin, error)
