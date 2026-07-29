# Antia AI Tool MCP - 开发计划

## 项目概述

为 antia-server 提供通用的 MCP (Model Context Protocol) 服务框架，支持插件化扩展多种内部服务。

**技术栈**: Go (Golang)  
**版本**: v0.1.0 (Phase 1 - Jenkins Plugin)  
**定位**: 通用 MCP Server 框架 + 插件生态系统

### 核心设计理念
- **插件化架构**: 每个服务（Jenkins、Redis、Elasticsearch）作为独立插件
- **统一接口**: 所有插件实现相同的Plugin接口，便于管理
- **按需加载**: 支持动态启用/禁用插件，减少资源占用
- **易扩展**: 新增服务只需开发新插件，无需修改核心框架

---

## 一、需求分析

### 1.1 功能需求

#### Phase 1 (一期) - 通用框架 + Jenkins Plugin
**核心框架**:
- 通用 MCP Server 实现（JSON-RPC 2.0）
- 插件管理系统（注册、加载、生命周期）
- 统一的配置系统（支持多插件配置）
- 日志和错误处理机制

**Jenkins Plugin** (v2.204.1 兼容):
- 支持用户名/密码认证（Basic Auth）
- 基于 REST API 实现（兼容低版本Jenkins）
- 提供 6 个核心 MCP Tools：
  - `jenkins_list_jobs` - 获取 Job 列表
  - `jenkins_get_job` - 查看 Job 详情和状态
  - `jenkins_trigger_build` - 触发 Job 构建（支持参数化）
  - `jenkins_get_build` - 获取构建详情
  - `jenkins_get_build_log` - 查看构建日志
  - `jenkins_list_builds` - 查看构建历史

#### Phase 2+ (后续扩展规划)
**预期支持的插件**:
- **Redis Plugin**: 键值操作、缓存管理、监控
- **Elasticsearch Plugin**: 索引查询、日志搜索、聚合分析
- **MySQL/PostgreSQL Plugin**: 数据库查询、Schema 查看
- **Kubernetes Plugin**: 集群管理、Pod 操作
- **GitLab/GitHub Plugin**: 代码仓库管理

**框架增强**:
- 插件热加载/热更新
- 插件间通信机制
- 统一的权限控制（RBAC）
- 操作审计日志
- 监控和告警集成

### 1.2 非功能需求
- **性能**: 响应时间 < 3s，支持并发请求
- **可靠性**: 自动重连机制，错误重试
- **安全性**: 支持 Jenkins 认证（用户名/密码、API Token）
- **可维护性**: 清晰的代码结构，完善的日志
- **可扩展性**: 模块化设计，易于添加新功能

---

## 二、技术方案

### 2.1 架构设计

#### 整体架构
```
┌──────────────────────────────────────┐
│         antia-server / Claude        │
│            (MCP Client)              │
└─────────────────┬────────────────────┘
                  │ MCP Protocol (stdio/HTTP)
                  │ JSON-RPC 2.0
┌─────────────────▼────────────────────┐
│      Antia AI Tool MCP Server        │
│  ┌────────────────────────────────┐  │
│  │     MCP Protocol Layer         │  │
│  │  • Initialize / Handshake      │  │
│  │  • tools/list, tools/call      │  │
│  │  • Request/Response Handler    │  │
│  └─────────────┬──────────────────┘  │
│                │                      │
│  ┌─────────────▼──────────────────┐  │
│  │    Plugin Management Layer     │  │
│  │  • Plugin Registry             │  │
│  │  • Lifecycle Manager           │  │
│  │  • Tool Router (统一分发)      │  │
│  │  • Config Injector             │  │
│  └─────────────┬──────────────────┘  │
│                │                      │
│  ┌─────────────▼──────────────────┐  │
│  │      Plugin Interface          │  │
│  │  ✓ GetTools()                  │  │
│  │  ✓ ExecuteTool(name, params)   │  │
│  │  ✓ Initialize(config)          │  │
│  └────────────────────────────────┘  │
└──────────────────┬───────────────────┘
                   │
        ┌──────────┼──────────┬─────────────┐
        │          │          │             │
   ┌────▼───┐ ┌───▼────┐ ┌──▼──────┐ ┌────▼────┐
   │Jenkins │ │ Redis  │ │Elastic  │ │ Future  │
   │Plugin  │ │Plugin  │ │Plugin   │ │ Plugins │
   │        │ │(Phase2)│ │(Phase2) │ │         │
   └────┬───┘ └───┬────┘ └──┬──────┘ └─────────┘
        │         │          │
   ┌────▼───┐ ┌──▼────┐ ┌───▼─────┐
   │Jenkins │ │ Redis │ │   ES    │
   │2.204.1 │ │Server │ │ Cluster │
   └────────┘ └───────┘ └─────────┘
```

#### 插件化设计原则
1. **插件独立性**: 每个插件独立目录，自包含所有逻辑
2. **统一接口**: 所有插件实现相同的 Plugin interface
3. **配置隔离**: 每个插件有独立的配置段
4. **错误隔离**: 单个插件错误不影响其他插件和核心框架

### 2.2 核心模块

#### 1. MCP Protocol Layer (MCP协议层)
**职责**: 实现 MCP 标准协议
- JSON-RPC 2.0 消息处理
- 标准 MCP 方法实现:
  - `initialize`: 初始化握手，返回服务器能力
  - `tools/list`: 聚合所有插件的工具列表
  - `tools/call`: 路由到对应插件执行
- 传输层抽象（stdio/HTTP）
- 请求验证和错误响应

#### 2. Plugin Management Layer (插件管理层)
**职责**: 插件生命周期和工具路由
- **Plugin Registry**: 插件注册表，维护插件实例
- **Lifecycle Manager**: 控制插件初始化、启动、停止
- **Tool Router**: 根据工具名称路由到对应插件
- **Config Injector**: 注入插件特定配置
- **Health Checker**: 插件健康检查

**插件发现机制**:
```go
// 自动发现和注册插件
func RegisterPlugin(name string, factory PluginFactory) {
    registry[name] = factory
}

// 插件工厂模式
type PluginFactory func(config interface{}) (Plugin, error)

// 内置插件自动注册
func init() {
    RegisterPlugin("jenkins", NewJenkinsPlugin)
    // 未来扩展
    // RegisterPlugin("redis", NewRedisPlugin)
    // RegisterPlugin("elasticsearch", NewElasticsearchPlugin)
}
```

#### 3. Plugin Interface (插件接口)
**统一插件接口定义**:
```go
// Plugin 是所有插件必须实现的接口
type Plugin interface {
    // 元信息
    Name() string
    Version() string
    Description() string
    
    // 生命周期
    Initialize(config PluginConfig) error
    Start() error
    Stop() error
    HealthCheck() error
    
    // 工具管理
    GetTools() []Tool
    ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error)
}

// Tool 定义 MCP 工具
type Tool struct {
    Name        string                 // 工具名称 (如: jenkins_list_jobs)
    Description string                 // 工具描述
    InputSchema map[string]interface{} // JSON Schema 参数定义
    OutputHint  string                 // 返回值说明
}

// PluginConfig 插件配置接口
type PluginConfig interface {
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    Unmarshal(v interface{}) error
}
```

#### 4. Jenkins Plugin (Jenkins插件)
**模块结构**:
```
plugins/jenkins/
├── plugin.go      # 插件主入口，实现Plugin接口
├── client.go      # Jenkins REST API客户端
├── tools.go       # MCP Tools定义和实现
├── config.go      # Jenkins配置结构
├── auth.go        # 认证管理（用户名/密码）
└── models.go      # 数据模型
```

**Jenkins Client 封装**:
- 基于 `net/http` 实现，兼容 Jenkins 2.204.1
- Basic Auth 认证（用户名/密码）
- REST API v1 endpoint 调用
- 连接池管理
- 自动重试机制
- 错误处理和日志记录

### 2.3 技术选型

| 组件 | 技术选择 | 理由 |
|------|---------|------|
| **MCP 实现** | 基于 `github.com/mark3labs/mcp-go` | 官方推荐的 Go MCP SDK |
| **HTTP 客户端** | 标准库 `net/http` + 自定义封装 | Jenkins 2.204.1 REST API 兼容，无需第三方库 |
| **配置管理** | `github.com/spf13/viper` | 支持多格式配置、环境变量、动态加载 |
| **日志** | `go.uber.org/zap` | 高性能结构化日志，适合生产环境 |
| **CLI** | `github.com/spf13/cobra` | 强大的命令行框架，易于扩展子命令 |
| **测试** | `github.com/stretchr/testify` | 丰富的断言和Mock支持 |
| **错误处理** | `github.com/pkg/errors` | 错误包装和堆栈追踪 |
| **JSON Schema** | `github.com/xeipuuv/gojsonschema` | MCP Tool 参数验证 |

**为什么不使用 gojenkins**:
- `gojenkins` 可能不完全兼容 Jenkins 2.204.1 的老版本 API
- 直接使用 REST API 更灵活，可精确控制兼容性
- 减少依赖，降低版本冲突风险
- 便于根据内网环境定制（如自签名证书处理）

---

## 三、Claude Code 完整开发流程

### Phase 1: 项目初始化 (预计 1-2 小时)

#### Task 1.1: 项目结构搭建
```
antia-aitool-mcp/
├── cmd/
│   └── server/
│       └── main.go                  # 服务入口
├── internal/
│   ├── mcp/
│   │   ├── server.go                # MCP 服务器核心
│   │   ├── handler.go               # 请求处理器
│   │   ├── transport_stdio.go       # stdio 传输实现
│   │   ├── transport_http.go        # HTTP 传输实现（可选）
│   │   └── protocol.go              # MCP 协议定义
│   ├── plugin/
│   │   ├── interface.go             # Plugin 接口定义
│   │   ├── manager.go               # 插件管理器
│   │   ├── registry.go              # 插件注册表
│   │   ├── config.go                # 插件配置抽象
│   │   └── router.go                # 工具路由器
│   └── plugins/
│       ├── jenkins/
│       │   ├── plugin.go            # Jenkins 插件实现
│       │   ├── client.go            # REST API 客户端
│       │   ├── tools.go             # 6个 MCP Tools
│       │   ├── auth.go              # Basic Auth 认证
│       │   ├── config.go            # Jenkins 配置
│       │   └── models.go            # 数据模型
│       └── README.md                # 插件开发指南
├── pkg/
│   ├── models/                      # 公共数据模型
│   │   ├── mcp.go                   # MCP 协议模型
│   │   └── errors.go                # 错误定义
│   └── utils/                       # 工具函数
│       ├── http.go                  # HTTP 工具
│       └── json.go                  # JSON 处理
├── configs/
│   ├── config.yaml                  # 默认配置
│   └── config.example.yaml          # 配置示例（含注释）
├── scripts/
│   ├── build.sh                     # 构建脚本
│   ├── test.sh                      # 测试脚本
│   └── install.sh                   # 安装脚本
├── docs/
│   ├── API.md                       # MCP Tools API 文档
│   ├── PLUGIN_DEV.md                # 插件开发指南
│   ├── DEPLOYMENT.md                # 部署文档
│   └── ARCHITECTURE.md              # 架构设计文档
├── tests/
│   ├── integration/                 # 集成测试
│   │   └── jenkins_test.go
│   ├── fixtures/                    # 测试数据
│   │   └── jenkins_responses.json
│   └── mocks/                       # Mock 对象
├── .github/
│   └── workflows/
│       └── ci.yml                   # CI/CD 配置
├── go.mod
├── go.sum
├── Makefile                         # 构建任务
├── .gitignore
├── README.md                        # 项目介绍
├── DEVELOPMENT_PLAN.md              # 本文档
├── CLAUDE.md                        # Claude Code 项目文档
└── LICENSE
```

**开发步骤**:
1. 初始化 Go 模块: `go mod init github.com/antia/antia-aitool-mcp`
2. 创建目录结构
3. 添加核心依赖:
   ```bash
   go get github.com/mark3labs/mcp-go@latest
   go get github.com/spf13/cobra@latest
   go get github.com/spf13/viper@latest
   go get go.uber.org/zap@latest
   go get github.com/pkg/errors@latest
   go get github.com/stretchr/testify@latest
   ```
4. 创建 CLAUDE.md 文档（项目上下文和开发指南）
5. 设置 .gitignore（排除配置敏感信息）

#### Task 1.2: 配置系统设计
- 定义配置文件格式（YAML）
- 实现配置加载和验证
- 支持环境变量覆盖
- 敏感信息处理（密码从环境变量读取）
- 插件配置隔离（每个插件独立配置段）

**配置示例**:
```yaml
# 服务器配置
server:
  mode: stdio              # stdio | http
  http:
    enabled: false         # HTTP 模式（可选，便于调试）
    host: "0.0.0.0"
    port: 8080
  log:
    level: info            # debug | info | warn | error
    format: json           # json | console
    file: ""               # 日志文件路径（空则输出到stderr）

# 插件配置
plugins:
  # Jenkins 插件配置
  jenkins:
    enabled: true
    url: "http://jenkins.internal.company.com"
    auth:
      username: "admin"
      # 密码从环境变量读取，避免明文存储
      password: "${JENKINS_PASSWORD}"
    options:
      timeout: 30s         # API 请求超时
      verify_ssl: false    # 内网自签名证书，设为false
      max_retries: 3       # 失败重试次数
      retry_delay: 2s      # 重试间隔
    
  # Redis 插件配置（Phase 2，暂时禁用）
  redis:
    enabled: false
    host: "localhost"
    port: 6379
    password: "${REDIS_PASSWORD}"
    db: 0
  
  # Elasticsearch 插件配置（Phase 2，暂时禁用）
  elasticsearch:
    enabled: false
    urls:
      - "http://es-node1:9200"
      - "http://es-node2:9200"
    username: "elastic"
    password: "${ES_PASSWORD}"
```

**配置加载优先级**:
1. 默认值（代码中定义）
2. 配置文件（`--config` 指定或默认路径）
3. 环境变量（`ANTIA_` 前缀，如 `ANTIA_SERVER_LOG_LEVEL=debug`）
4. 命令行参数（最高优先级）

---

### Phase 2: MCP Server 核心实现 (预计 3-4 小时)

#### Task 2.1: MCP 协议实现
- 实现 MCP JSON-RPC 2.0 协议
- 支持标准 MCP 方法:
  - `initialize`: 初始化连接
  - `tools/list`: 列出可用工具
  - `tools/call`: 调用工具
  - `resources/list`: 列出可用资源（可选）
- 错误处理和响应格式化

#### Task 2.2: 传输层实现
- **Stdio 传输**: 标准输入输出通信（Claude Desktop 默认）
- **HTTP 传输**: RESTful API 接口（可选，便于测试）
- 消息序列化/反序列化
- 连接管理和心跳检测

#### Task 2.3: 插件管理系统
- 定义插件接口 `Plugin interface`
- 实现插件注册机制（工厂模式）
- 插件生命周期管理（Initialize, Start, Stop）
- 插件配置注入
- 工具路由器（根据工具名称分发到对应插件）
- 插件健康检查

**插件接口设计**:
```go
// Plugin 是所有插件必须实现的接口
type Plugin interface {
    // 元信息
    Name() string                    // 插件名称（如: jenkins）
    Version() string                 // 插件版本（如: v1.0.0）
    Description() string             // 插件描述
    
    // 生命周期
    Initialize(config PluginConfig) error  // 初始化配置
    Start() error                          // 启动插件（建立连接等）
    Stop() error                           // 停止插件（清理资源）
    HealthCheck() error                    // 健康检查
    
    // 工具管理
    GetTools() []Tool                      // 返回插件提供的所有工具
    ExecuteTool(ctx context.Context, 
                name string, 
                params map[string]interface{}) (interface{}, error)
}

// Tool 定义 MCP 工具
type Tool struct {
    Name        string                 // 工具名称（如: jenkins_list_jobs）
    Description string                 // 工具描述
    InputSchema map[string]interface{} // JSON Schema 参数定义
}

// PluginManager 管理所有插件
type PluginManager struct {
    plugins  map[string]Plugin        // 已加载的插件
    registry map[string]PluginFactory // 插件工厂注册表
    logger   *zap.Logger
}

// 工厂函数类型
type PluginFactory func(config PluginConfig) (Plugin, error)

// 注册插件（在各插件的 init() 中调用）
func RegisterPlugin(name string, factory PluginFactory) {
    defaultRegistry[name] = factory
}
```

**工具路由逻辑**:
```go
// 根据工具名称路由到对应插件
func (m *PluginManager) ExecuteTool(ctx context.Context, 
                                     toolName string, 
                                     params map[string]interface{}) (interface{}, error) {
    // 工具名称格式: {plugin_name}_{tool_name}
    // 例如: jenkins_list_jobs -> jenkins 插件的 list_jobs 工具
    parts := strings.SplitN(toolName, "_", 2)
    if len(parts) < 2 {
        return nil, errors.New("invalid tool name format")
    }
    
    pluginName := parts[0]
    plugin, exists := m.plugins[pluginName]
    if !exists {
        return nil, fmt.Errorf("plugin %s not found", pluginName)
    }
    
    return plugin.ExecuteTool(ctx, toolName, params)
}
```

---

### Phase 3: Jenkins Plugin 开发 (预计 4-5 小时)

#### Task 3.1: Jenkins REST API 客户端封装
**目标**: 实现兼容 Jenkins 2.204.1 的 REST API 客户端

**核心功能**:
1. **认证管理**
   - Basic Auth（用户名/密码）
   - 自动添加 Authorization header
   - 会话保持和重用

2. **REST API 封装**
   - 基于 Jenkins REST API v1
   - 关键 endpoints:
     ```
     GET  /api/json                          # 服务器信息
     GET  /api/json?tree=jobs[name,color]    # Job 列表
     GET  /job/{name}/api/json               # Job 详情
     POST /job/{name}/build                  # 触发构建（无参数）
     POST /job/{name}/buildWithParameters    # 触发构建（带参数）
     GET  /job/{name}/{number}/api/json      # 构建详情
     GET  /job/{name}/{number}/consoleText   # 构建日志
     GET  /job/{name}/api/json?tree=builds[number,result,timestamp,duration]  # 构建历史
     ```

3. **错误处理**
   - HTTP 状态码处理（401/403认证错误、404未找到、500服务器错误）
   - 网络超时和重试
   - 响应解析错误处理
   - 友好的错误信息转换

4. **连接管理**
   - HTTP 客户端配置（超时、连接池）
   - SSL/TLS 配置（支持跳过证书验证）
   - 自动重试机制（指数退避）

**Jenkins 2.204.1 兼容性注意事项**:
- 使用稳定的 REST API v1，避免使用新版本特性
- 某些新 API 可能不可用（如 Blue Ocean API），使用传统 API
- 测试时注意 API 响应格式的差异

**客户端结构**:
```go
type JenkinsClient struct {
    baseURL    string
    httpClient *http.Client
    username   string
    password   string
    logger     *zap.Logger
}

func (c *JenkinsClient) GetJobs(folder string) ([]Job, error)
func (c *JenkinsClient) GetJob(name string) (*JobDetail, error)
func (c *JenkinsClient) TriggerBuild(name string, params map[string]string) error
func (c *JenkinsClient) GetBuild(job string, number int) (*BuildDetail, error)
func (c *JenkinsClient) GetBuildLog(job string, number int, start int) (string, error)
func (c *JenkinsClient) GetBuilds(job string, limit int) ([]Build, error)
```

#### Task 3.2: MCP Tools 定义
实现以下 MCP Tools:

1. **jenkins_list_jobs**
   - 描述: 获取所有 Jenkins Job 列表
   - 参数: `folder` (可选) - Job 文件夹路径
   - 返回: Job 名称、类型、状态列表

2. **jenkins_get_job**
   - 描述: 获取指定 Job 的详细信息
   - 参数: `job_name` (必需) - Job 名称
   - 返回: Job 配置、最近构建信息

3. **jenkins_trigger_build**
   - 描述: 触发 Job 构建
   - 参数: `job_name`, `parameters` (可选)
   - 返回: 构建队列信息

4. **jenkins_get_build**
   - 描述: 获取指定构建的详细信息
   - 参数: `job_name`, `build_number`
   - 返回: 构建状态、持续时间、结果

5. **jenkins_get_build_log**
   - 描述: 获取构建日志
   - 参数: `job_name`, `build_number`, `start` (可选)
   - 返回: 日志文本

6. **jenkins_list_builds**
   - 描述: 获取 Job 的构建历史
   - 参数: `job_name`, `limit` (可选)
   - 返回: 构建列表

#### Task 3.3: 数据转换和格式化
- Jenkins API 响应转换为 MCP 标准格式
- 错误信息规范化
- 日志输出优化（大日志分页）

---

### Phase 4: 测试和验证 (预计 3-4 小时)

#### Task 4.1: 单元测试
- MCP Server 核心逻辑测试
- Jenkins 客户端 Mock 测试
- 配置加载测试
- 覆盖率目标: > 70%

#### Task 4.2: 集成测试
- 使用 testcontainers 启动 Jenkins 测试实例
- 端到端工具调用测试
- 错误场景测试（网络超时、认证失败）

#### Task 4.3: 手动测试
- 使用 MCP Inspector 工具测试
- 集成到 Claude Desktop 测试
- 性能和压力测试

**测试清单**:
```bash
# 启动服务
./antia-aitool-mcp --config configs/config.yaml

# 测试工具列表
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | ./antia-aitool-mcp

# 测试 Jenkins 连接
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"jenkins_list_jobs"},"id":2}' | ./antia-aitool-mcp
```

---

### Phase 5: 文档和部署 (预计 2-3 小时)

#### Task 5.1: 文档编写
- README.md: 项目介绍、快速开始
- API.md: MCP Tools 详细文档
- DEPLOYMENT.md: 部署指南
- PLUGIN_DEV.md: 插件开发指南（为后续扩展准备）

#### Task 5.2: 部署准备
- 编写 Dockerfile
- 创建 systemd service 文件
- 编写部署脚本
- CI/CD 配置（GitHub Actions / GitLab CI）

#### Task 5.3: Claude Desktop 集成配置
创建 Claude Desktop 配置文件示例:
```json
{
  "mcpServers": {
    "antia-jenkins": {
      "command": "/path/to/antia-aitool-mcp",
      "args": ["--config", "/path/to/config.yaml"],
      "env": {
        "JENKINS_API_TOKEN": "your-token-here"
      }
    }
  }
}
```

---

## 四、开发建议和最佳实践

### 4.1 使用 Claude Code 的建议

1. **初始化阶段**
   - 使用 `/init` 创建 CLAUDE.md 项目文档
   - 描述项目架构、依赖、构建流程
   - 记录关键决策和设计原则

2. **开发阶段**
   - 使用 Task 系统跟踪进度: `/task create`, `/task list`
   - 遇到复杂问题使用 `/plan` 模式规划方案
   - 使用 `/review` 进行代码审查
   - 使用 `/simplify` 优化代码质量

3. **测试阶段**
   - 使用 `/run` 快速启动测试
   - 使用 Agent 工具并行执行测试任务
   - 使用 `/security-review` 检查安全问题

4. **文档阶段**
   - 使用 Claude 生成 API 文档
   - 使用 `codebase-memory` 技能理解代码结构
   - 生成架构图和流程图

### 4.2 Go 开发最佳实践

1. **代码组织**
   - 遵循 Go 标准项目布局
   - internal/ 目录保护内部 API
   - pkg/ 目录提供可重用的公共 API

2. **错误处理**
   - 使用 `errors.Wrap` 添加上下文
   - 定义清晰的错误类型
   - 避免 panic，使用 error 返回

3. **并发安全**
   - Jenkins 客户端使用连接池
   - 使用 sync.Mutex 保护共享状态
   - 使用 context 控制超时和取消

4. **性能优化**
   - 使用 pprof 进行性能分析
   - 合理使用缓存（Job 信息缓存）
   - 异步处理长时间操作

5. **安全考虑**
   - 敏感信息不记录到日志
   - 使用环境变量存储凭证
   - 验证所有外部输入
   - 支持 SSL/TLS 验证

### 4.3 MCP 协议注意事项

1. **Tool Schema 设计**
   - 使用 JSON Schema 定义参数
   - 提供清晰的描述和示例
   - 合理的默认值

2. **错误响应**
   - 遵循 JSON-RPC 2.0 错误格式
   - 提供有用的错误信息
   - 区分客户端错误和服务器错误

3. **性能考虑**
   - 大数据分页返回
   - 实现超时机制
   - 支持取消操作

---

## 五、里程碑和时间规划

| 阶段 | 任务 | 预计时间 | 交付物 |
|------|------|---------|--------|
| Phase 1 | 项目初始化 | 1-2 小时 | 项目结构、配置系统 |
| Phase 2 | MCP Server 实现 | 3-4 小时 | MCP 协议、插件系统 |
| Phase 3 | Jenkins Plugin | 4-5 小时 | Jenkins 集成、6 个 Tools |
| Phase 4 | 测试验证 | 3-4 小时 | 测试套件、集成测试 |
| Phase 5 | 文档部署 | 2-3 小时 | 完整文档、部署方案 |
| **总计** | | **13-18 小时** | 生产就绪的 MCP 服务 |

---

## 六、后续扩展方向

### 6.1 功能扩展

#### Redis Plugin (示例)
**预期 MCP Tools**:
- `redis_get` - 获取键值
- `redis_set` - 设置键值
- `redis_keys` - 查询键列表
- `redis_info` - 服务器信息
- `redis_monitor` - 实时监控命令

**插件结构**:
```
plugins/redis/
├── plugin.go      # 实现 Plugin 接口
├── client.go      # Redis 客户端封装
├── tools.go       # MCP Tools 定义
└── config.go      # Redis 配置
```

#### Elasticsearch Plugin (示例)
**预期 MCP Tools**:
- `es_search` - 搜索文档
- `es_indices` - 列出索引
- `es_get_document` - 获取文档
- `es_aggregation` - 聚合查询
- `es_cluster_health` - 集群健康状态

#### 插件开发流程
1. 在 `internal/plugins/` 下创建新目录
2. 实现 `Plugin` 接口
3. 定义 MCP Tools
4. 在 `init()` 函数中注册插件
5. 更新配置文件模板
6. 编写单元测试和集成测试
7. 更新文档

**插件开发模板**: 详见 `docs/PLUGIN_DEV.md`

### 6.2 性能和可靠性
- 添加缓存层（Redis）
- 实现分布式部署
- 添加监控和告警（Prometheus）
- 实现请求限流和熔断

### 6.3 安全和权限
- OAuth2 认证集成
- RBAC 权限控制
- 操作审计日志
- 敏感数据加密存储

---

## 七、风险和挑战

### 7.1 技术风险
- **Jenkins API 兼容性**: Jenkins 2.204.1 是 2019 年的版本，需确保 REST API 兼容
  - *缓解*: 使用稳定的 REST API v1，充分测试所有 endpoints
  - *验证*: 在真实 Jenkins 2.204.1 环境测试所有 API 调用
  
- **用户名/密码认证**: Basic Auth 需处理特殊字符和编码
  - *缓解*: 使用标准库 base64 编码，测试特殊字符场景
  
- **网络稳定性**: 内网环境可能存在网络波动
  - *缓解*: 实现自动重试（指数退避）、超时控制、健康检查
  
- **插件扩展性**: 架构是否真正支持插件化
  - *缓解*: Phase 1 设计时预留扩展点，使用接口隔离
  - *验证*: Phase 2 尝试快速添加一个简单插件验证架构

### 7.2 集成风险
- **MCP 协议变更**: MCP 协议仍在演进中
  - *缓解*: 使用稳定的 SDK 版本，关注官方更新
- **Claude Desktop 兼容性**: 不同版本可能有差异
  - *缓解*: 测试多个版本，提供配置示例

---

## 八、成功标准

### 8.1 功能完整性
- ✅ 实现 6 个核心 Jenkins MCP Tools
- ✅ 支持 stdio 传输方式
- ✅ 配置系统完整可用
- ✅ 错误处理完善

### 8.2 质量标准
- ✅ 单元测试覆盖率 > 70%
- ✅ 集成测试通过
- ✅ 代码通过 golangci-lint 检查
- ✅ 无已知安全漏洞

### 8.3 可用性标准
- ✅ 能够成功集成到 Claude Desktop
- ✅ 所有 Tools 在实际 Jenkins 环境测试通过
- ✅ 完整的用户文档和部署指南
- ✅ 响应时间符合性能要求

---

## 九、开始开发

### 准备工作

**环境要求**:
- Go 1.21+ 
- Git
- 访问内网 Jenkins 2.204.1 服务器
- Claude Desktop（用于测试 MCP 集成）

**建议的开发流程**:

```bash
# 1. 初始化 Go 项目
cd /Users/admin/projects/antia/antia-aitool-mcp
go mod init github.com/antia/antia-aitool-mcp

# 2. 创建项目结构
mkdir -p cmd/server internal/{mcp,plugin,plugins/jenkins} pkg/{models,utils} configs docs tests/{integration,fixtures}

# 3. 使用 Claude Code 初始化项目文档
# 在 Claude Code 中执行: /init
```

### 推荐开发顺序

**第一步: 核心框架** (先跑通最小可用系统)
1. 实现最简单的 MCP Server（处理 initialize 和 tools/list）
2. 实现最简单的插件接口
3. 手写一个 mock 插件验证框架可用

**第二步: Jenkins 插件** (实现业务功能)
1. 先实现一个最简单的 Jenkins Tool（如 `jenkins_list_jobs`）
2. 测试端到端流程（从 MCP Client 到 Jenkins API）
3. 逐步添加其他 5 个 Tools

**第三步: 完善和优化**
1. 添加错误处理、日志、重试
2. 编写单元测试和集成测试
3. 完善配置系统
4. 编写文档

### 使用 Claude Code 的建议命令

```bash
# 开始 Phase 1 开发
"Let's start Phase 1: initialize the project structure and implement basic config system"

# 实现 MCP Server
"Implement Phase 2: MCP Protocol Layer with stdio transport"

# 实现 Jenkins Plugin
"Implement Phase 3: Jenkins Plugin with REST API client and 6 MCP Tools"

# 代码审查和优化
/simplify    # 简化代码
/review      # 代码审查
/security-review  # 安全审查

# 运行测试
/run         # 启动服务测试
```

### 开发里程碑验证

**Milestone 1**: 框架跑通
- [ ] MCP Server 能响应 `initialize` 请求
- [ ] 能正确列出工具列表（`tools/list`）
- [ ] 插件系统能加载 mock 插件

**Milestone 2**: Jenkins 集成
- [ ] 能连接 Jenkins 2.204.1（Basic Auth）
- [ ] `jenkins_list_jobs` 能返回 Job 列表
- [ ] 其他 5 个 Tools 全部实现

**Milestone 3**: 生产就绪
- [ ] 所有测试通过（单元 + 集成）
- [ ] 能集成到 Claude Desktop
- [ ] 文档完整（README、API、部署）
- [ ] 性能和稳定性满足要求

---

### 关键决策记录

**为什么不使用第三方 Jenkins 库?**
- Jenkins 2.204.1 是 2019 年版本，第三方库可能不兼容
- 直接使用 REST API 更灵活，便于调试和定制
- 减少依赖，降低维护成本

**为什么选择插件化架构?**
- 后续需要扩展 Redis、Elasticsearch 等多种服务
- 插件化便于独立开发、测试、部署
- 核心框架稳定后，可以并行开发多个插件

**为什么优先支持 stdio 而非 HTTP?**
- Claude Desktop 默认使用 stdio 通信
- stdio 更简单、更安全（无需暴露端口）
- HTTP 作为可选项，便于调试和测试

---

**准备好开始了吗？** 建议首先执行：

```
/init
```

然后告诉我："Let's start with Phase 1, Task 1.1: Initialize the Go project and create the directory structure"

---

**文档版本**: v1.0  
**创建日期**: 2026-07-24  
**作者**: Claude Code Assistant  
**项目代号**: antia-aitool-mcp
