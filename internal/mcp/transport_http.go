package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/huianlei/antia-aitool-mcp/pkg/models"
	"go.uber.org/zap"
)

// HTTPTransport implements MCP over HTTP
type HTTPTransport struct {
	server *http.Server
	mcp    *Server
	logger *zap.Logger
}

// NewHTTPTransport creates a new HTTP transport
func NewHTTPTransport(host string, port int, mcp *Server, logger *zap.Logger) *HTTPTransport {
	return &HTTPTransport{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		mcp:    mcp,
		logger: logger,
	}
}

// Start starts the HTTP server
func (t *HTTPTransport) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", t.handleMCPRequest)
	mux.HandleFunc("/health", t.handleHealth)

	t.server.Handler = mux

	t.logger.Info("starting HTTP transport", zap.String("addr", t.server.Addr))

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := t.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		t.logger.Info("shutting down HTTP server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return t.server.Shutdown(shutdownCtx)
	}
}

// handleMCPRequest handles MCP JSON-RPC requests
func (t *HTTPTransport) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.logger.Error("failed to read request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON-RPC request
	var request models.MCPRequest
	if err := json.Unmarshal(body, &request); err != nil {
		t.logger.Error("failed to parse JSON-RPC request", zap.Error(err))
		http.Error(w, "Invalid JSON-RPC request", http.StatusBadRequest)
		return
	}

	t.logger.Debug("received MCP request",
		zap.String("method", request.Method),
		zap.Any("id", request.ID))

	// Handle the request
	response := t.mcp.HandleRequest(r.Context(), &request)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.logger.Error("failed to encode response", zap.Error(err))
	}
}

// handleHealth handles health check requests
func (t *HTTPTransport) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"server": "antia-aitool-mcp",
	})
}
