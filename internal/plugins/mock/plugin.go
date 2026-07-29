package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/huianlei/antia-aitool-mcp/internal/plugin"
	"github.com/huianlei/antia-aitool-mcp/pkg/models"
)

// MockPlugin is a simple test plugin
type MockPlugin struct {
	config plugin.Config
}

func init() {
	// Register mock plugin
	plugin.Register("mock", NewMockPlugin)
}

// NewMockPlugin creates a new mock plugin instance
func NewMockPlugin(config plugin.Config) (plugin.Plugin, error) {
	return &MockPlugin{
		config: config,
	}, nil
}

// Name returns the plugin name
func (p *MockPlugin) Name() string {
	return "mock"
}

// Version returns the plugin version
func (p *MockPlugin) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (p *MockPlugin) Description() string {
	return "Mock plugin for testing MCP framework"
}

// Initialize initializes the plugin
func (p *MockPlugin) Initialize(config plugin.Config) error {
	p.config = config
	return nil
}

// Start starts the plugin
func (p *MockPlugin) Start() error {
	return nil
}

// Stop stops the plugin
func (p *MockPlugin) Stop() error {
	return nil
}

// HealthCheck checks plugin health
func (p *MockPlugin) HealthCheck() error {
	return nil
}

// GetTools returns the list of tools provided by this plugin
func (p *MockPlugin) GetTools() []models.Tool {
	return []models.Tool{
		{
			Name:        "mock_echo",
			Description: "Echo back the input message",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"message": map[string]interface{}{
						"type":        "string",
						"description": "Message to echo",
					},
				},
				"required": []string{"message"},
			},
		},
		{
			Name:        "mock_time",
			Description: "Get current server time",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "mock_add",
			Description: "Add two numbers",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "number",
						"description": "First number",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "Second number",
					},
				},
				"required": []string{"a", "b"},
			},
		},
	}
}

// ExecuteTool executes a tool
func (p *MockPlugin) ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	switch name {
	case "mock_echo":
		return p.executeEcho(params)
	case "mock_time":
		return p.executeTime(params)
	case "mock_add":
		return p.executeAdd(params)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func (p *MockPlugin) executeEcho(params map[string]interface{}) (interface{}, error) {
	message, ok := params["message"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'message' parameter")
	}

	return map[string]interface{}{
		"echo":      message,
		"timestamp": time.Now().Format(time.RFC3339),
	}, nil
}

func (p *MockPlugin) executeTime(params map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"time":     time.Now().Format(time.RFC3339),
		"unix":     time.Now().Unix(),
		"timezone": time.Now().Location().String(),
	}, nil
}

func (p *MockPlugin) executeAdd(params map[string]interface{}) (interface{}, error) {
	a, aOk := params["a"].(float64)
	b, bOk := params["b"].(float64)

	if !aOk || !bOk {
		return nil, fmt.Errorf("parameters 'a' and 'b' must be numbers")
	}

	return map[string]interface{}{
		"a":      a,
		"b":      b,
		"result": a + b,
	}, nil
}
