package jenkins

import (
	"context"
	"fmt"
	"time"

	"github.com/huianlei/antia-aitool-mcp/internal/plugin"
	"github.com/huianlei/antia-aitool-mcp/pkg/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Plugin implements the Jenkins plugin
type Plugin struct {
	config plugin.Config
	client *Client
	logger *zap.Logger
}

func init() {
	// Register Jenkins plugin
	plugin.Register("jenkins", NewJenkinsPlugin)
}

// NewJenkinsPlugin creates a new Jenkins plugin instance
func NewJenkinsPlugin(config plugin.Config) (plugin.Plugin, error) {
	return &Plugin{
		config: config,
	}, nil
}

// Name returns the plugin name
func (p *Plugin) Name() string {
	return "jenkins"
}

// Version returns the plugin version
func (p *Plugin) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (p *Plugin) Description() string {
	return "Jenkins 2.204.1 integration plugin"
}

// Initialize initializes the plugin
func (p *Plugin) Initialize(config plugin.Config) error {
	p.config = config

	// Create logger
	logger, err := zap.NewProduction()
	if err != nil {
		return errors.Wrap(err, "failed to create logger")
	}
	p.logger = logger

	// Parse configuration
	jenkinsConfig := Config{
		URL:        config.GetString("url"),
		Username:   getNestedString(config.GetMap("auth"), "username"),
		Password:   getNestedString(config.GetMap("auth"), "password"),
		Timeout:    30 * time.Second,
		VerifySSL:  true,
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
	}

	// Get options if available
	if options := config.GetMap("options"); options != nil {
		if timeout, ok := options["timeout"].(string); ok {
			if d, err := time.ParseDuration(timeout); err == nil {
				jenkinsConfig.Timeout = d
			}
		}
		if verifySSL, ok := options["verify_ssl"].(bool); ok {
			jenkinsConfig.VerifySSL = verifySSL
		}
		if maxRetries, ok := options["max_retries"].(int); ok {
			jenkinsConfig.MaxRetries = maxRetries
		}
		if retryDelay, ok := options["retry_delay"].(string); ok {
			if d, err := time.ParseDuration(retryDelay); err == nil {
				jenkinsConfig.RetryDelay = d
			}
		}
	}

	// Validate configuration
	if jenkinsConfig.URL == "" {
		return fmt.Errorf("Jenkins URL is required")
	}
	if jenkinsConfig.Username == "" || jenkinsConfig.Password == "" {
		return fmt.Errorf("Jenkins username and password are required")
	}

	// Create Jenkins client
	p.client = NewClient(jenkinsConfig, p.logger)

	p.logger.Info("Jenkins plugin initialized",
		zap.String("url", jenkinsConfig.URL),
		zap.String("username", jenkinsConfig.Username),
	)

	return nil
}

// Start starts the plugin
func (p *Plugin) Start() error {
	// Test connection to Jenkins
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.client.Ping(ctx); err != nil {
		p.logger.Warn("failed to ping Jenkins", zap.Error(err))
		return errors.Wrap(err, "failed to connect to Jenkins")
	}

	p.logger.Info("Jenkins plugin started successfully")
	return nil
}

// Stop stops the plugin
func (p *Plugin) Stop() error {
	p.logger.Info("Jenkins plugin stopped")
	return nil
}

// HealthCheck checks plugin health
func (p *Plugin) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.client.Ping(ctx)
}

// GetTools returns the list of tools provided by this plugin
func (p *Plugin) GetTools() []models.Tool {
	return []models.Tool{
		{
			Name:        "jenkins_list_jobs",
			Description: "List all Jenkins jobs",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "jenkins_get_job",
			Description: "Get detailed information about a Jenkins job",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"job_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the Jenkins job",
					},
				},
				"required": []string{"job_name"},
			},
		},
		{
			Name:        "jenkins_trigger_build",
			Description: "Trigger a build for a Jenkins job",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"job_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the Jenkins job",
					},
					"parameters": map[string]interface{}{
						"type":        "object",
						"description": "Build parameters (optional)",
					},
				},
				"required": []string{"job_name"},
			},
		},
		{
			Name:        "jenkins_get_build",
			Description: "Get detailed information about a specific build",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"job_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the Jenkins job",
					},
					"build_number": map[string]interface{}{
						"type":        "integer",
						"description": "Build number",
					},
				},
				"required": []string{"job_name", "build_number"},
			},
		},
		{
			Name:        "jenkins_get_build_log",
			Description: "Get console log for a build",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"job_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the Jenkins job",
					},
					"build_number": map[string]interface{}{
						"type":        "integer",
						"description": "Build number",
					},
					"start": map[string]interface{}{
						"type":        "integer",
						"description": "Starting byte offset (optional)",
					},
				},
				"required": []string{"job_name", "build_number"},
			},
		},
		{
			Name:        "jenkins_list_builds",
			Description: "List build history for a Jenkins job",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"job_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the Jenkins job",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of builds to return (optional)",
					},
				},
				"required": []string{"job_name"},
			},
		},
	}
}

// ExecuteTool executes a tool
func (p *Plugin) ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	switch name {
	case "jenkins_list_jobs":
		return p.listJobs(ctx, params)
	case "jenkins_get_job":
		return p.getJob(ctx, params)
	case "jenkins_trigger_build":
		return p.triggerBuild(ctx, params)
	case "jenkins_get_build":
		return p.getBuild(ctx, params)
	case "jenkins_get_build_log":
		return p.getBuildLog(ctx, params)
	case "jenkins_list_builds":
		return p.listBuilds(ctx, params)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// Helper function to get nested string value
func getNestedString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
