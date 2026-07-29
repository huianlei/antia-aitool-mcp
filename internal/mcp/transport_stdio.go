package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/huianlei/antia-aitool-mcp/pkg/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// StdioTransport handles MCP communication over stdin/stdout
type StdioTransport struct {
	server *Server
	logger *zap.Logger
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewStdioTransport creates a new stdio transport
func NewStdioTransport(server *Server, logger *zap.Logger) *StdioTransport {
	return &StdioTransport{
		server: server,
		logger: logger,
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}

// Start starts the stdio transport loop
func (t *StdioTransport) Start(ctx context.Context) error {
	t.logger.Info("starting stdio transport")

	// Channel to receive read results
	type readResult struct {
		req *models.MCPRequest
		err error
	}
	readChan := make(chan readResult, 1)

	for {
		// Start async read
		go func() {
			req, err := t.readRequest()
			readChan <- readResult{req: req, err: err}
		}()

		select {
		case <-ctx.Done():
			t.logger.Info("stdio transport shutting down")
			return ctx.Err()

		case result := <-readChan:
			if result.err != nil {
				if result.err == io.EOF {
					t.logger.Info("stdin closed, shutting down")
					return nil
				}
				t.logger.Error("failed to read request", zap.Error(result.err))
				continue
			}

			// Process request
			resp := t.server.HandleRequest(ctx, result.req)

			// Write response to stdout
			if err := t.writeResponse(resp); err != nil {
				t.logger.Error("failed to write response", zap.Error(err))
				continue
			}
		}
	}
}

// readRequest reads a JSON-RPC request from stdin
func (t *StdioTransport) readRequest() (*models.MCPRequest, error) {
	// Read line from stdin
	line, err := t.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// Parse JSON-RPC request
	var req models.MCPRequest
	if err := json.Unmarshal(line, &req); err != nil {
		return nil, errors.Wrap(err, "failed to parse JSON-RPC request")
	}

	t.logger.Debug("received request",
		zap.String("method", req.Method),
		zap.Any("id", req.ID),
	)

	return &req, nil
}

// writeResponse writes a JSON-RPC response to stdout
func (t *StdioTransport) writeResponse(resp *models.MCPResponse) error {
	// Marshal response to JSON
	data, err := json.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON-RPC response")
	}

	// Write to stdout
	if _, err := t.writer.Write(data); err != nil {
		return errors.Wrap(err, "failed to write to stdout")
	}

	// Write newline
	if err := t.writer.WriteByte('\n'); err != nil {
		return errors.Wrap(err, "failed to write newline")
	}

	// Flush buffer
	if err := t.writer.Flush(); err != nil {
		return errors.Wrap(err, "failed to flush stdout")
	}

	t.logger.Debug("sent response", zap.Any("id", resp.ID))

	return nil
}
