package mcp

import (
	"context"
	"fmt"

	"github.com/huianlei/antia-aitool-mcp/internal/plugin"
	"github.com/huianlei/antia-aitool-mcp/pkg/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// MCP Protocol Version
	ProtocolVersion = "2024-11-05"

	// Server Info
	ServerName    = "antia-aitool-mcp"
	ServerVersion = "0.1.0"
)

// Server represents the MCP server
type Server struct {
	pluginManager *plugin.Manager
	logger        *zap.Logger
	serverInfo    models.ServerInfo
}

// NewServer creates a new MCP server instance
func NewServer(pluginManager *plugin.Manager, logger *zap.Logger) *Server {
	return &Server{
		pluginManager: pluginManager,
		logger:        logger,
		serverInfo: models.ServerInfo{
			Name:    ServerName,
			Version: ServerVersion,
		},
	}
}

// HandleRequest processes an MCP request and returns a response
func (s *Server) HandleRequest(ctx context.Context, req *models.MCPRequest) *models.MCPResponse {
	s.logger.Debug("handling MCP request",
		zap.String("method", req.Method),
		zap.Any("id", req.ID),
	)

	// Route to appropriate handler
	switch req.Method {
	case "initialize":
		return s.handleInitialize(ctx, req)
	case "tools/list":
		return s.handleToolsList(ctx, req)
	case "tools/call":
		return s.handleToolsCall(ctx, req)
	default:
		return s.errorResponse(req.ID, models.ErrorCodeMethodNotFound,
			fmt.Sprintf("method not found: %s", req.Method))
	}
}

// handleInitialize handles the initialize method
func (s *Server) handleInitialize(ctx context.Context, req *models.MCPRequest) *models.MCPResponse {
	s.logger.Info("initializing MCP server")

	result := models.InitializeResult{
		ProtocolVersion: ProtocolVersion,
		Capabilities:    s.serverInfo,
		ServerInfo:      s.serverInfo,
	}

	return &models.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleToolsList handles the tools/list method
func (s *Server) handleToolsList(ctx context.Context, req *models.MCPRequest) *models.MCPResponse {
	s.logger.Debug("listing tools")

	// Get all tools from plugin manager
	tools := s.pluginManager.GetAllTools()

	result := models.ToolsListResult{
		Tools: tools,
	}

	s.logger.Info("tools listed", zap.Int("count", len(tools)))

	return &models.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleToolsCall handles the tools/call method
func (s *Server) handleToolsCall(ctx context.Context, req *models.MCPRequest) *models.MCPResponse {
	// Extract tool call parameters
	params, ok := req.Params["arguments"].(map[string]interface{})
	if !ok {
		params = req.Params
	}

	toolName, ok := req.Params["name"].(string)
	if !ok {
		return s.errorResponse(req.ID, models.ErrorCodeInvalidParams,
			"missing or invalid 'name' parameter")
	}

	s.logger.Info("executing tool",
		zap.String("tool", toolName),
	)

	// Execute tool via plugin manager
	result, err := s.pluginManager.ExecuteTool(ctx, toolName, params)
	if err != nil {
		s.logger.Error("tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
		)

		// Map error to appropriate error code
		code := models.ErrorCodeInternalError
		if errors.Is(err, models.ErrPluginNotFound) || errors.Is(err, models.ErrToolNotFound) {
			code = models.ErrorCodeMethodNotFound
		} else if errors.Is(err, models.ErrInvalidParameters) {
			code = models.ErrorCodeInvalidParams
		}

		return s.errorResponse(req.ID, code, err.Error())
	}

	// Wrap result in content block
	toolResult := models.ToolCallResult{
		Content: []models.ContentBlock{
			{
				Type: "text",
				Text: fmt.Sprintf("%v", result),
			},
		},
	}

	s.logger.Info("tool executed successfully",
		zap.String("tool", toolName),
	)

	return &models.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  toolResult,
	}
}

// errorResponse creates an error response
func (s *Server) errorResponse(id interface{}, code int, message string) *models.MCPResponse {
	return &models.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &models.MCPError{
			Code:    code,
			Message: message,
		},
	}
}

// Shutdown gracefully shuts down the MCP server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down MCP server")
	return s.pluginManager.Shutdown()
}
