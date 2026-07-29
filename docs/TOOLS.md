# Development Tools

This document lists the development tools used in this project and their versions.

## Required Tools

### Go 1.25.12
This project requires Go 1.25.12 or higher. Use gvm for version management:
```bash
gvm install go1.25.12
gvm use go1.25.12
```

### Development Tools

Install all development tools at once:
```bash
./scripts/install_tools.sh
```

Or install individually:

#### gopls v0.17.1
Go language server for IDE integration.
```bash
go install golang.org/x/tools/gopls@v0.17.1
```

#### delve v1.24.0
Go debugger.
```bash
go install github.com/go-delve/delve/cmd/dlv@v1.24.0
```

#### staticcheck v0.5.1
Static analysis tool.
```bash
go install honnef.co/go/tools/cmd/staticcheck@v0.5.1
```

#### golangci-lint v1.62.2
Comprehensive linter (required by `make lint`).
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
```

#### gotests
Test generation tool.
```bash
go install github.com/cweill/gotests/gotests@latest
```

#### gomodifytags
Struct tag modifier.
```bash
go install github.com/fatih/gomodifytags@latest
```

#### impl
Interface implementation generator.
```bash
go install github.com/josharian/impl@latest
```

## Tool Versions Compatibility

| Tool | Version | Go 1.25.12 Compatible | Notes |
|------|---------|----------------------|-------|
| gopls | v0.17.1 | ✅ Yes | Updated for Go 1.25 support |
| delve | v1.24.0 | ✅ Yes | Latest stable debugger |
| staticcheck | v0.5.1 | ✅ Yes | Latest static analyzer |
| golangci-lint | v1.62.2 | ✅ Yes | Required by Makefile |
| gotests | latest | ✅ Yes | Test generator |
| gomodifytags | latest | ✅ Yes | Tag modifier |
| impl | latest | ✅ Yes | Interface generator |

## Updating Tools

To update all tools to the latest compatible versions:
```bash
./scripts/install_tools.sh
```

## IDE Configuration

### VS Code
Install the Go extension and it will automatically use gopls.

### GoLand / IntelliJ IDEA
Built-in Go support, no additional tools needed.

### Vim/Neovim
Use vim-go or coc-go plugin with gopls.

## Verification

Check installed tools:
```bash
gopls version
dlv version
staticcheck -version
golangci-lint version
```

Expected output:
```
gopls v0.17.1
Delve Debugger v1.24.0
staticcheck v0.5.1
golangci-lint version 1.62.2
```
