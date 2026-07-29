# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Antia AI Tool MCP** is a universal MCP (Model Context Protocol) server framework written in Go that enables Claude to interact with internal services through a plugin architecture. Phase 1 focuses on Jenkins 2.204.1 integration, with future phases expanding to Redis, Elasticsearch, and other services.

**Key Characteristics**:
- Plugin-based architecture where each service is an independent plugin
- MCP protocol implementation over stdio (primary) and HTTP (optional)
- Designed for internal network environments (self-signed certs, Basic Auth)
- Jenkins 2.204.1 compatibility using REST API v1 (no third-party Jenkins libraries)

## Architecture

### Three-Layer Design

```
MCP Protocol Layer (internal/mcp/)
    ↓ JSON-RPC 2.0 over stdio/HTTP
Plugin Management Layer (internal/plugin/)
    ↓ Tool routing, lifecycle management
Plugin Implementations (internal/plugins/)
    ↓ Jenkins, Redis, Elasticsearch, etc.
```

**Critical Design Decisions**:
1. **Plugin isolation**: Each plugin is self-contained in `internal/plugins/{name}/` with its own client, tools, config, and models
2. **Tool naming convention**: `{plugin}_{action}` (e.g., `jenkins_list_jobs`) enables automatic routing to the correct plugin
3. **No third-party Jenkins libraries**: Direct REST API implementation for Jenkins 2.204.1 compatibility and control
4. **Config-driven plugin loading**: Plugins enabled/disabled via YAML config (`plugins.{name}.enabled`)

### Plugin Interface

All plugins implement this interface (defined in `internal/plugin/interface.go`):

```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    
    Initialize(config PluginConfig) error
    Start() error
    Stop() error
    HealthCheck() error
    
    GetTools() []Tool
    ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error)
}
```

### Tool Router Logic

The Plugin Manager routes MCP tool calls by parsing the tool name prefix:
- `jenkins_list_jobs` → routes to Jenkins plugin
- `redis_get` → routes to Redis plugin
- Tool name format is `{plugin_name}_{tool_action}`

## Project Structure

```
cmd/server/          # Entry point (main.go)
internal/
  mcp/               # MCP protocol implementation (JSON-RPC 2.0, stdio/HTTP transport)
  plugin/            # Plugin management (registry, lifecycle, router, config injection)
  plugins/
    jenkins/         # Jenkins plugin implementation
      plugin.go      # Plugin interface implementation
      client.go      # REST API client (Basic Auth, Jenkins 2.204.1)
      tools.go       # 6 MCP tools definitions
      auth.go        # Authentication handling
      config.go      # Jenkins-specific config
      models.go      # Data models
pkg/
  models/            # Shared data models (MCP protocol, errors)
  utils/             # Common utilities (HTTP, JSON helpers)
configs/             # YAML configuration files
tests/
  integration/       # Integration tests (testcontainers for Jenkins)
  fixtures/          # Test data (mock API responses)
```

## Development Commands

### Setup
```bash
# Initialize project
go mod init github.com/antia/antia-aitool-mcp
go mod tidy

# Install dependencies
go get github.com/mark3labs/mcp-go@latest
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get go.uber.org/zap@latest
go get github.com/pkg/errors@latest
go get github.com/stretchr/testify@latest
```

### Build
```bash
# Build binary
go build -o antia-aitool-mcp ./cmd/server

# Build with version info
go build -ldflags "-X main.version=$(git describe --tags)" -o antia-aitool-mcp ./cmd/server
```

### Run
```bash
# Run with default config
./antia-aitool-mcp --config configs/config.yaml

# Run with environment variables
export JENKINS_PASSWORD="your-password"
./antia-aitool-mcp --config configs/config.yaml

# Run in stdio mode (for Claude Desktop)
./antia-aitool-mcp --config configs/config.yaml

# Run in HTTP mode (for testing)
./antia-aitool-mcp --config configs/config.yaml --http
```

### Test
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/plugins/jenkins/...

# Run integration tests (requires Jenkins test instance)
go test -tags=integration ./tests/integration/...

# Run single test
go test -run TestJenkinsClient_GetJobs ./internal/plugins/jenkins/
```

### Lint
```bash
# Run golangci-lint
golangci-lint run

# Format code
go fmt ./...
goimports -w .
```

## Configuration System

Configuration uses Viper with YAML files. Priority order: CLI args > env vars > config file > defaults.

**Environment variable prefix**: `ANTIA_`
- Example: `ANTIA_SERVER_LOG_LEVEL=debug`
- Sensitive values use `${VAR}` syntax in YAML: `password: "${JENKINS_PASSWORD}"`

**Plugin configuration structure**:
```yaml
plugins:
  jenkins:
    enabled: true
    url: "http://jenkins.internal.company.com"
    auth:
      username: "admin"
      password: "${JENKINS_PASSWORD}"  # From env var
    options:
      timeout: 30s
      verify_ssl: false  # Internal network with self-signed certs
      max_retries: 3
```

## Jenkins Plugin Implementation Notes

### Jenkins 2.204.1 Compatibility
- Uses REST API v1 endpoints (stable, backward-compatible)
- No Blue Ocean or Pipeline-specific APIs (may not exist in old versions)
- Basic Auth with username/password (no API token in this version)

### Key REST API Endpoints
```
GET  /api/json                          # Server info
GET  /api/json?tree=jobs[name,color]    # Job list
GET  /job/{name}/api/json               # Job details
POST /job/{name}/build                  # Trigger build (no params)
POST /job/{name}/buildWithParameters    # Trigger build (with params)
GET  /job/{name}/{number}/api/json      # Build details
GET  /job/{name}/{number}/consoleText   # Build log
```

### Jenkins Client Implementation
Located in `internal/plugins/jenkins/client.go`:
- Uses standard `net/http` (no third-party Jenkins library)
- Implements exponential backoff retry for failed requests
- Handles 401/403 auth errors distinctly from 404/500 errors
- Supports `verify_ssl: false` for self-signed certificates
- Connection pooling via http.Client configuration

## Adding a New Plugin

1. Create plugin directory: `internal/plugins/{name}/`
2. Implement files:
   - `plugin.go` - implements Plugin interface
   - `client.go` - service-specific API client
   - `tools.go` - MCP tool definitions
   - `config.go` - plugin configuration struct
3. Register plugin in `init()`:
   ```go
   func init() {
       plugin.RegisterPlugin("myplugin", NewMyPlugin)
   }
   ```
4. Add config section to `configs/config.yaml`:
   ```yaml
   plugins:
     myplugin:
       enabled: false
       # plugin-specific config
   ```
5. Update `docs/PLUGIN_DEV.md` with plugin-specific details

## MCP Tools Naming Convention

Tool names follow the pattern: `{plugin_name}_{action}`

**Jenkins Plugin Tools** (Phase 1):
- `jenkins_list_jobs` - List all jobs
- `jenkins_get_job` - Get job details
- `jenkins_trigger_build` - Trigger a build
- `jenkins_get_build` - Get build details
- `jenkins_get_build_log` - Get build console log
- `jenkins_list_builds` - List build history

Each tool defines its input schema using JSON Schema for parameter validation.

## Testing Strategy

**Unit Tests**: Test individual components in isolation
- Mock external dependencies (Jenkins API, HTTP clients)
- Focus on business logic, error handling, edge cases
- Located alongside source files (`*_test.go`)

**Integration Tests**: Test against real services
- Use testcontainers to spin up Jenkins instance
- Test full request/response cycle
- Located in `tests/integration/`
- Run with: `go test -tags=integration`

**End-to-End Tests**: Test MCP protocol integration
- Test with MCP Inspector or Claude Desktop
- Verify tool discovery, parameter validation, execution
- Manual testing checklist in `docs/DEPLOYMENT.md`

## Error Handling

Use `github.com/pkg/errors` for error wrapping with context:

```go
if err != nil {
    return errors.Wrap(err, "failed to fetch Jenkins jobs")
}
```

Plugin errors should be descriptive and include:
- What operation failed
- Why it failed (HTTP status, network error, etc.)
- Relevant context (job name, build number, etc.)

## Logging

Uses `go.uber.org/zap` for structured logging:
- Log levels: debug, info, warn, error
- Never log sensitive data (passwords, tokens)
- Include request context (plugin name, tool name, request ID)
- Use structured fields: `logger.Info("job triggered", zap.String("job", jobName))`

## Security Considerations

1. **Credential handling**: Never log or expose passwords/tokens in responses
2. **Input validation**: Validate all tool parameters against JSON Schema
3. **SSL/TLS**: Support `verify_ssl: false` for internal networks, but warn in logs
4. **Error messages**: Avoid leaking internal paths or sensitive configuration
5. **Environment variables**: Store credentials in env vars, not config files

## Claude Desktop Integration

Add to Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "antia-jenkins": {
      "command": "/path/to/antia-aitool-mcp",
      "args": ["--config", "/path/to/config.yaml"],
      "env": {
        "JENKINS_PASSWORD": "your-password-here"
      }
    }
  }
}
```

## Development Workflow

**Recommended order**:
1. Implement core MCP Server (Protocol Layer)
2. Implement Plugin Manager (Management Layer)
3. Create a simple mock plugin to validate framework
4. Implement Jenkins plugin with one tool (`jenkins_list_jobs`)
5. Test end-to-end with MCP Inspector
6. Add remaining Jenkins tools
7. Add comprehensive tests and documentation

**When modifying plugin interface**: Update all existing plugins to maintain compatibility.

**When adding new tools**: Update tool documentation in `docs/API.md` with examples.

## Future Expansion

Phase 2+ will add:
- Redis plugin (key-value operations, monitoring)
- Elasticsearch plugin (search, aggregations)
- Database plugins (MySQL, PostgreSQL)
- Kubernetes plugin (cluster management)

The plugin architecture is designed to support these additions without modifying core MCP or plugin management layers.
