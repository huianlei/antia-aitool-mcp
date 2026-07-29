#!/bin/bash
# Development environment setup script
# Run this when first setting up the project

set -e

echo "==> Antia AI Tool MCP - Development Setup"
echo ""

# Check if gvm is installed
if ! command -v gvm &> /dev/null; then
    echo "⚠ gvm is not installed"
    echo "Please install gvm first: https://github.com/moovweb/gvm"
    exit 1
fi

# Check .go-version file
if [ ! -f .go-version ]; then
    echo "⚠ .go-version file not found"
    exit 1
fi

REQUIRED_GO_VERSION=$(cat .go-version)
echo "==> Required Go version: $REQUIRED_GO_VERSION"

# Check if the required Go version is installed
if ! gvm list | grep -q "$REQUIRED_GO_VERSION"; then
    echo "⚠ Go $REQUIRED_GO_VERSION is not installed in gvm"
    echo "Please install it first:"
    echo "  gvm install $REQUIRED_GO_VERSION"
    exit 1
fi

# Switch to the required Go version
echo "==> Switching to Go $REQUIRED_GO_VERSION..."
gvm use "$REQUIRED_GO_VERSION"

# Verify Go version
CURRENT_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [ "$CURRENT_VERSION" != "$REQUIRED_GO_VERSION" ]; then
    echo "⚠ Failed to switch to Go $REQUIRED_GO_VERSION"
    echo "Current version: $CURRENT_VERSION"
    exit 1
fi

echo "✓ Go version: $(go version)"
echo ""

# Install dependencies
echo "==> Installing Go dependencies..."
go mod download
go mod tidy

echo ""
echo "✓ Dependencies installed"
echo ""

# Create config file if not exists
if [ ! -f configs/config.yaml ]; then
    echo "==> Creating config.yaml from example..."
    cp configs/config.example.yaml configs/config.yaml
    echo "✓ Created configs/config.yaml"
    echo "  Please edit it with your Jenkins credentials"
    echo ""
fi

# Build project
echo "==> Building project..."
make build

echo ""
echo "✓ Setup complete!"
echo ""
echo "To use this project:"
echo "  1. Run 'gvm use $REQUIRED_GO_VERSION' to switch Go version"
echo "  2. Or add this to your shell profile:"
echo "     cd /path/to/antia-aitool-mcp && source .envrc"
echo ""
echo "Next steps:"
echo "  1. Edit configs/config.yaml with your settings"
echo "  2. Set environment variables: export JENKINS_PASSWORD='...'"
echo "  3. Run the server: make run"
echo ""
