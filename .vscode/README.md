# VS Code Configuration

This directory contains VS Code workspace settings for the Antia AI Tool MCP project.

## Files Overview

### settings.json
Workspace-specific settings for Go development:
- **Language Server**: gopls configuration
- **Formatting**: goimports on save
- **Linting**: golangci-lint integration
- **Testing**: Test flags and coverage
- **Editor**: Tab settings, rulers, etc.

### launch.json
Debug configurations:
1. **Launch Server** - Run MCP server with config
2. **Launch Server (Mock Plugin Only)** - Quick testing with mock plugin
3. **Test Current File** - Debug single test file
4. **Test Current Package** - Debug package tests
5. **Test All** - Debug all tests
6. **Test with Coverage** - Run tests with coverage report
7. **Attach to Process** - Attach debugger to running process

### tasks.json
Build and development tasks (accessible via `Cmd+Shift+B` or `Ctrl+Shift+P` → Tasks):
- **Build** - Compile the project
- **Clean** - Remove build artifacts
- **Test** - Run all tests
- **Test with Coverage** - Generate coverage report
- **Lint** - Run golangci-lint
- **Format Code** - Format all Go files
- **Run Server** - Start the MCP server
- **Install Dependencies** - Run `go mod download`
- **Install Dev Tools** - Install development tools

### extensions.json
Recommended VS Code extensions:
- `golang.go` - Official Go extension (required)
- `eamodio.gitlens` - Git history and blame
- `ms-vscode.makefile-tools` - Makefile support
- `redhat.vscode-yaml` - YAML validation
- `streetsidesoftware.code-spell-checker` - Spell checker
- `davidanson.vscode-markdownlint` - Markdown linter
- `yzhang.markdown-all-in-one` - Markdown enhancements

## Quick Start

### 1. Install Recommended Extensions
When you open this workspace, VS Code will prompt you to install recommended extensions. Click "Install All".

Or manually:
```
Cmd/Ctrl+Shift+P → Extensions: Show Recommended Extensions
```

### 2. Install Go Tools
```bash
./scripts/install_tools.sh
```

### 3. Verify Setup
Open a `.go` file and check:
- ✅ Syntax highlighting works
- ✅ Code completion works (Ctrl+Space)
- ✅ "Go to Definition" works (F12)
- ✅ Save auto-formats the code

### 4. Run the Server
Press `F5` or:
```
Cmd/Ctrl+Shift+P → Debug: Start Debugging → Launch Server
```

## Keyboard Shortcuts

### Building and Testing
- `Cmd/Ctrl+Shift+B` - Build project
- `Cmd/Ctrl+Shift+T` - Run tests

### Debugging
- `F5` - Start debugging
- `Shift+F5` - Stop debugging
- `F9` - Toggle breakpoint
- `F10` - Step over
- `F11` - Step into
- `Shift+F11` - Step out

### Code Navigation
- `F12` - Go to definition
- `Alt+F12` - Peek definition
- `Shift+F12` - Find all references
- `Cmd/Ctrl+T` - Go to symbol

### Editing
- `Shift+Alt+F` - Format document
- `Cmd/Ctrl+.` - Quick fix
- `F2` - Rename symbol

## Debugging Tips

### 1. Debug Server with Breakpoints
1. Open `cmd/server/main.go`
2. Click left of line number to set breakpoint
3. Press `F5` to start debugging
4. Server will pause at breakpoint

### 2. Debug Specific Plugin
Set breakpoint in plugin code:
- `internal/plugins/jenkins/plugin.go`
- `internal/plugins/jenkins/client.go`

### 3. Debug Tests
1. Open test file (e.g., `client_test.go`)
2. Set breakpoint in test function
3. Use "Test Current File" debug config

## Troubleshooting

### gopls Not Working
```bash
# Reinstall gopls
go install golang.org/x/tools/gopls@v0.17.1

# Restart VS Code
Cmd/Ctrl+Shift+P → Developer: Reload Window
```

### Linter Errors
```bash
# Install golangci-lint
./scripts/install_tools.sh

# Verify installation
golangci-lint version
```

### Go Version Issues
Ensure you're using Go 1.25.12:
```bash
gvm use go1.25.12
go version
```

Then restart VS Code.

## Configuration Customization

### Personal Settings
Don't edit `.vscode/settings.json` directly. Instead, add personal preferences to:
- **Workspace settings**: `Cmd/Ctrl+,` → Workspace tab
- **User settings**: `Cmd/Ctrl+,` → User tab

### Environment Variables
Edit `.vscode/launch.json` to add environment variables for debugging:
```json
"env": {
  "JENKINS_PASSWORD": "your-password",
  "DEBUG": "true"
}
```

## Git Ignore
The `.vscode/` directory is tracked in git because it contains useful shared configurations. Only `.vscode/settings.json` contains project-wide settings; personal settings should go in User Settings.

If you want to ignore your personal modifications:
```bash
git update-index --skip-worktree .vscode/settings.json
```
