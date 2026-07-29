# Go 版本隔离指南

## 问题背景

您有两个项目需要不同的 Go 版本：
- **antia-server**（游戏服务器）: Go 1.16.15
- **antia-aitool-mcp**（本项目）: Go 1.25.12

使用 gvm (Go Version Manager) 可以实现版本隔离，互不影响。

## 解决方案

### 1. 项目级别版本控制

本项目使用 `.go-version` 文件指定 Go 版本：

```bash
# 文件: .go-version
go1.25.12
```

### 2. 使用 gvm 切换版本

**工作在本项目时**:
```bash
cd /path/to/antia-aitool-mcp
gvm use go1.25.12
go version          # 验证: go version go1.25.12
make build
```

**切换回游戏服务器**:
```bash
cd /path/to/antia-server
gvm use go1.16.15
go version          # 验证: go version go1.16.15
```

### 3. 自动切换（推荐）

在您的 `~/.zshrc` 或 `~/.bashrc` 中添加：

```bash
# 自动根据 .go-version 切换 Go 版本
cd() {
    builtin cd "$@"
    if [ -f ".go-version" ]; then
        REQUIRED_VERSION=$(cat .go-version)
        if command -v gvm &> /dev/null; then
            gvm use "$REQUIRED_VERSION" > /dev/null 2>&1
        fi
    fi
}
```

这样每次 `cd` 到项目目录时会自动切换 Go 版本。

### 4. 使用 direnv（更强大的方案）

安装 direnv:
```bash
brew install direnv
```

添加到 shell 配置:
```bash
# ~/.zshrc
eval "$(direnv hook zsh)"
```

本项目已包含 `.envrc` 文件，会自动切换 Go 版本。

启用 direnv:
```bash
cd antia-aitool-mcp
direnv allow .
```

## 验证版本隔离

### 测试 1: 本项目使用正确版本
```bash
cd /path/to/antia-aitool-mcp
gvm use go1.25.12
go version
# 输出: go version go1.25.12 darwin/amd64

go env GOROOT
# 输出: /Users/admin/.gvm/gos/go1.25.12

./antia-aitool-mcp version
# 输出: Version: dev, Build Time: ...
```

### 测试 2: 游戏服务器版本未受影响
```bash
cd /path/to/antia-server
gvm use go1.16.15
go version
# 输出: go version go1.16.15 darwin/amd64

go env GOROOT
# 输出: /Users/admin/.gvm/gos/go1.16.15
```

## 常见问题

### Q1: 为什么 go.mod 显示 go 1.25？
**A**: 项目使用 Go 1.25.12。`go.mod` 中的 `go 1.25` 表示最低要求版本，`toolchain go1.25.12` 指定具体工具链版本。

### Q2: 这会影响我的 antia-server 吗？
**A**: 不会。只要在 antia-server 目录下使用 `gvm use go1.16.15`，就完全隔离。每个项目使用各自的 Go 版本。

### Q3: 如何确认当前使用的 Go 版本？
**A**: 运行以下命令：
```bash
go version              # 查看版本
go env GOROOT          # 查看 Go 安装路径
which go               # 查看 Go 可执行文件路径
```

### Q4: 如果不小心在 antia-server 目录运行了错误的 Go 版本怎么办？
**A**: 
1. 立即切换回正确版本: `gvm use go1.16.15`
2. 重新构建: `go build`
3. gvm 的版本切换是会话级别的，不会永久改变全局设置

## 最佳实践

1. **每次进入项目目录都手动切换版本**（如果没有自动切换）:
   ```bash
   cd antia-aitool-mcp && gvm use go1.25.12
   cd antia-server && gvm use go1.16.15
   ```

2. **在 shell 提示符显示当前 Go 版本**（可选）:
   ```bash
   # ~/.zshrc
   PS1='[$(go version | cut -d" " -f3)] %~ $ '
   ```

3. **使用项目专用的 Makefile**:
   本项目的 Makefile 会使用当前激活的 Go 版本，不会影响其他项目。

4. **CI/CD 中指定版本**:
   在 GitHub Actions / GitLab CI 中使用：
   ```yaml
   - uses: actions/setup-go@v4
     with:
       go-version: '1.25'
   ```

## 项目隔离保证

✅ **go.mod**: 每个项目有独立的 go.mod 文件  
✅ **GOROOT**: gvm 为每个版本创建独立的 GOROOT  
✅ **依赖**: `go mod vendor` 或 `go.mod` 确保依赖版本独立  
✅ **二进制**: 编译的二进制文件包含所需的运行时，不依赖 Go 版本  

**结论**: 两个项目完全隔离，互不影响！
