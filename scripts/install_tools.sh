#!/usr/bin/env bash
# Install development tools for Go 1.25.12

set -e

echo "==> Installing Go development tools for Go 1.25.12"
echo ""

# Language Server - use older version compatible with Go 1.22-1.23
echo "Installing gopls (Go language server)..."
go install golang.org/x/tools/gopls@v0.21.1

# Debugger
echo "Installing delve (Go debugger)..."
go install github.com/go-delve/delve/cmd/dlv@v1.25.2

# Static analysis
echo "Installing staticcheck..."
go install honnef.co/go/tools/cmd/staticcheck@latest

# Linter (required by Makefile)
echo "Installing golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Test generation
echo "Installing gotests..."
go install github.com/cweill/gotests/gotests@latest

# Struct tag modifier
echo "Installing gomodifytags..."
go install github.com/fatih/gomodifytags@latest

# Interface implementation generator
echo "Installing impl..."
go install github.com/josharian/impl@latest

echo ""
echo "==> All tools installed successfully!"
echo ""
echo "Installed tools:"
echo "  - gopls (Language Server)"
echo "  - delve (Debugger)"
echo "  - staticcheck (Static Analysis)"
echo "  - golangci-lint (Linter)"
echo "  - gotests (Test Generator)"
echo "  - gomodifytags (Tag Modifier)"
echo "  - impl (Interface Generator)"
echo ""
echo "Run 'go version' to verify you're using Go 1.25.12"
echo "You can now use 'make lint' to run golangci-lint"
