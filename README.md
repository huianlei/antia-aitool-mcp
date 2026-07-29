# Antia AI Tool MCP

A universal MCP (Model Context Protocol) server framework that enables Claude to interact with internal services through a plugin architecture.

## Features

- **Plugin Architecture**: Extensible design for multiple service integrations
- **MCP Protocol**: JSON-RPC 2.0 over stdio/HTTP
- **Phase 1**: Jenkins 2.204.1 integration with 6 core tools
- **Future**: Redis, Elasticsearch, Kubernetes, and more

## Quick Start

### Prerequisites

- Go 1.25+ (managed via gvm for version isolation)
- gvm (Go Version Manager) - https://github.com/moovweb/gvm
- Access to internal services (Jenkins, etc.)

### Installation

```bash
# Clone the repository
git clone https://github.com/huianlei/antia-aitool-mcp.git
cd antia-aitool-mcp

# Install required Go version (if not already installed)
gvm install go1.25.12

# Run setup script (handles Go version switching and dependencies)
./scripts/setup.sh

# Or manually:
gvm use go1.25.12
make install-deps
make build
```

### Go Version Isolation

This project uses **Go 1.25.12** isolated via gvm. Your global Go version (e.g., Go 1.16 for antia-server) remains unchanged.

**Project-level version file**: `.go-version` specifies `go1.25.12`

**To work on this project**:
```bash
cd /path/to/antia-aitool-mcp
gvm use go1.25.12              # Switch to project Go version
go version                      # Verify: should show go1.25.12
```

**To return to your global Go version**:
```bash
cd /path/to/antia-server
gvm use go1.16.15              # Switch back to your server's Go version
```

**Optional - Auto-switching**:
Add to your `~/.zshrc` or `~/.bashrc`:
```bash
# Auto-switch Go version based on .go-version file
cd() {
    builtin cd "$@"
    if [ -f ".go-version" ]; then
        gvm use $(cat .go-version) > /dev/null 2>&1
    fi
}
```

### Configuration

```bash
# Copy example config
cp configs/config.example.yaml configs/config.yaml

# Edit configuration
vim configs/config.yaml

# Set environment variables for sensitive data
export JENKINS_PASSWORD="your-jenkins-password"
```

### Running

```bash
# Run with config file
./antia-aitool-mcp --config configs/config.yaml

# Or use make
make run
```

## Jenkins Plugin (Phase 1)

### Available Tools

1. **jenkins_list_jobs** - List all Jenkins jobs
2. **jenkins_get_job** - Get job details and status
3. **jenkins_trigger_build** - Trigger a job build (with optional parameters)
4. **jenkins_get_build** - Get build details
5. **jenkins_get_build_log** - Get build console log
6. **jenkins_list_builds** - List build history for a job

### Jenkins Configuration

```yaml
plugins:
  jenkins:
    enabled: true
    url: "http://jenkins.internal.company.com"
    auth:
      username: "admin"
      password: "${JENKINS_PASSWORD}"
    options:
      timeout: 30s
      verify_ssl: false  # For self-signed certificates
      max_retries: 3
```

## Claude Desktop Integration

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "antia-jenkins": {
      "command": "/path/to/antia-aitool-mcp",
      "args": ["--config", "/path/to/config.yaml"],
      "env": {
        "JENKINS_PASSWORD": "your-password"
      }
    }
  }
}
```

## Development

### Project Structure

```
cmd/server/          # Entry point
internal/
  mcp/               # MCP protocol implementation
  plugin/            # Plugin management
  plugins/jenkins/   # Jenkins plugin
pkg/                 # Shared packages
configs/             # Configuration files
docs/                # Documentation
tests/               # Tests
```

### Commands

```bash
make help            # Show all available commands
make build           # Build binary
make test            # Run tests
make test-coverage   # Run tests with coverage
make lint            # Run linter
make fmt             # Format code
```

### Adding a New Plugin

See `docs/PLUGIN_DEV.md` for detailed instructions.

## Architecture

```
MCP Client (Claude) 
    ↓ stdio/HTTP
MCP Server (Protocol Layer)
    ↓
Plugin Manager (Router, Lifecycle)
    ↓
Plugins (Jenkins, Redis, ES, ...)
```

## Documentation

- [Development Plan](DEVELOPMENT_PLAN.md) - Complete development roadmap
- [Claude Code Guide](CLAUDE.md) - Instructions for Claude Code
- [API Documentation](docs/API.md) - MCP Tools API reference
- [Plugin Development](docs/PLUGIN_DEV.md) - How to create plugins
- [Deployment Guide](docs/DEPLOYMENT.md) - Production deployment

## Roadmap

- **Phase 1** ✓ Framework + Jenkins Plugin
- **Phase 2**: Redis Plugin
- **Phase 3**: Elasticsearch Plugin
- **Phase 4**: Additional plugins (K8s, Databases, Git)

## Go Version Management

✅ **Go Version Isolation**: This project uses Go 1.25.12 managed by gvm, isolated from your global Go installation.

- **This project**: Go 1.25.12 (via gvm)
- **antia-server**: Go 1.16.15 (unchanged)
- **Isolation method**: gvm + `.go-version` file

No conflict between projects!

## License

[Your License]

## Contributing

See `DEVELOPMENT_PLAN.md` for development guidelines.
