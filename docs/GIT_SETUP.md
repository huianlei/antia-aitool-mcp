# Git 配置建议

准备将项目推送到 GitHub 时的步骤：

## 1. 初始化 Git 仓库

```bash
cd /Users/admin/projects/antia/antia-aitool-mcp
git init
```

## 2. 配置 .gitignore

已包含在项目中，确保以下内容被忽略：
- 二进制文件 (`antia-aitool-mcp`)
- 配置文件中的敏感信息 (`configs/config.yaml`)
- 日志文件
- 临时文件

## 3. 首次提交

```bash
# 添加所有文件
git add .

# 创建首次提交
git commit -m "Initial commit: Antia AI Tool MCP Server v0.1.0

Features:
- MCP Protocol Layer (JSON-RPC 2.0 over stdio)
- Plugin system framework
- Jenkins plugin (6 MCP tools for Jenkins 2.204.1)
- Mock plugin for testing
- Configuration system with environment variable support
- Go version isolation with gvm"
```

## 4. 添加远程仓库

```bash
# 添加 GitHub 远程仓库
git remote add origin https://github.com/huianlei/antia-aitool-mcp.git

# 或使用 SSH
git remote add origin git@github.com:huianlei/antia-aitool-mcp.git
```

## 5. 推送到 GitHub

```bash
# 创建主分支并推送
git branch -M main
git push -u origin main
```

## 6. 推荐的分支策略

```bash
# 开发新功能
git checkout -b feature/redis-plugin
# ... 开发 ...
git commit -m "Add Redis plugin"
git push origin feature/redis-plugin
# 然后在 GitHub 上创建 Pull Request

# 修复 bug
git checkout -b fix/jenkins-auth-issue
# ... 修复 ...
git commit -m "Fix Jenkins authentication error handling"
git push origin fix/jenkins-auth-issue
```

## 7. 标签和版本

```bash
# 创建版本标签
git tag -a v0.1.0 -m "Release v0.1.0 - Initial release with Jenkins plugin"
git push origin v0.1.0

# 查看所有标签
git tag -l
```

## 8. 保护敏感信息

确保 **不要** 提交以下文件：
- ✅ `configs/config.yaml` (包含真实密码) - 已在 .gitignore
- ✅ `configs/config.example.yaml` (安全，使用占位符) - 可以提交
- ✅ 环境变量配置
- ✅ 日志文件

## 9. GitHub 仓库设置建议

### README Badges (可选)
```markdown
![Go Version](https://img.shields.io/badge/Go-1.25-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Build Status](https://github.com/huianlei/antia-aitool-mcp/workflows/CI/badge.svg)
```

### Topics (建议添加)
- `mcp`
- `model-context-protocol`
- `claude`
- `jenkins`
- `plugin-system`
- `golang`
- `ai-tools`

### GitHub Actions (CI/CD)

创建 `.github/workflows/ci.yml`:
```yaml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25'
    
    - name: Build
      run: make build
    
    - name: Test
      run: make test
    
    - name: Lint
      run: |
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        make lint
```

## 10. 文档检查清单

推送前确认：
- ✅ README.md 完整且更新
- ✅ CLAUDE.md 包含项目指南
- ✅ LICENSE 文件（如果需要）
- ✅ CONTRIBUTING.md（如果接受贡献）
- ✅ 所有 import 路径更新为 `github.com/huianlei/antia-aitool-mcp`

## 11. 推送前最后检查

```bash
# 检查所有文件状态
git status

# 检查没有硬编码的敏感信息
grep -r "password" --include="*.go" . | grep -v "Password"
grep -r "10.7.71.6" --include="*.go" .

# 确认 import 路径正确
grep -r "github.com/antia/antia-aitool-mcp" --include="*.go" .
# 应该没有结果

grep -r "github.com/huianlei/antia-aitool-mcp" --include="*.go" . | head -3
# 应该显示新的路径

# 最后测试构建
make clean
make build
./antia-aitool-mcp version
```

## 12. 首次推送命令序列

```bash
# 完整流程
cd /Users/admin/projects/antia/antia-aitool-mcp

# 初始化（如果还没有）
git init

# 添加文件
git add .

# 提交
git commit -m "Initial commit: Antia AI Tool MCP Server v0.1.0"

# 添加远程
git remote add origin https://github.com/huianlei/antia-aitool-mcp.git

# 推送
git branch -M main
git push -u origin main

# 添加标签
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

---

完成后，您的项目将在 https://github.com/huianlei/antia-aitool-mcp 上线！
