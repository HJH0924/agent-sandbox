# Agent Sandbox Service

一个为 AI Agent 设计的安全沙箱服务，提供隔离的文件操作和命令执行环境。基于 Go 和 gRPC/Connect-RPC 构建。

## 特性

- **安全沙箱环境**：每个沙箱实例拥有独立的工作空间
- **API Key 认证**：使用 X-Sandbox-Api-Key 进行安全访问控制
- **文件操作**：支持文件读取、写入和编辑
- **Shell 命令执行**：支持超时控制的命令执行
- **结构化日志**：使用 slog 进行 JSON 格式日志记录
- **云端部署**：支持 E2B 和 Docker 部署

## 快速开始

### 前置要求

- Go 1.24+
- Make
- Docker（可选，用于容器部署）

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/HJH0924/agent-sandbox.git
cd agent-sandbox

# 安装开发依赖（推荐）
make install-deps

# 安装 Go 模块依赖
go mod tidy

# 构建并运行
make build
./bin/api-server --config configs/config.yaml
```

服务将在 `http://localhost:8080` 启动。

## 基本使用

### 1. 初始化沙箱

```bash
curl -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}'
```

响应示例：
```json
{
  "sandboxId": "550e8400-e29b-41d4-a716-446655440000",
  "apiKey": "sk_0123456789abcdef...",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

### 2. 执行命令

```bash
curl -X POST http://localhost:8080/shell.v1.ShellService/Execute \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"command": "ls -la"}'
```

### 3. 文件操作

**写入文件：**
```bash
curl -X POST http://localhost:8080/file.v1.FileService/Write \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"path": "example.txt", "content": "Hello, World!"}'
```

**读取文件：**
```bash
curl -X POST http://localhost:8080/file.v1.FileService/Read \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"path": "example.txt"}'
```

## 本地测试

### 使用 Docker 模拟 E2B 环境

本地可以使用 Docker 来模拟 E2B sandbox 环境，用于开发和测试：

```bash
# 启动 Docker 容器（使用与 E2B 相同的 Dockerfile）
make docker

# 查看日志
docker-compose logs -f

# 进入容器测试
make shell
```

## 生产部署

### E2B Sandbox 部署

**生产环境推荐使用 E2B 云端沙箱**，专为 AI Agent 设计。

```bash
# 安装 E2B CLI
npm install -g @e2b/cli

# 设置 API Key（从 https://e2b.dev/dashboard 获取）
export E2B_API_KEY=your_api_key

# 构建 E2B 模板
make e2b

# 创建沙箱
e2b sandbox spawn agent-sandbox
```

详细的 E2B 使用说明请查看[开发指南](docs/development.md#e2b-模板部署)。

## 配置

编辑 `configs/config.yaml`：

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

sandbox:
  workspace_dir: "/tmp/agent-sandbox"
  max_file_size: 104857600  # 100MB
  shell_timeout: 300        # 5分钟

log:
  level: "info"  # debug, info, warn, error
  format: "json" # json, text
```

## 开发

查看[开发指南](docs/development.md)了解详细的开发说明。

**常用命令：**

```bash
make install-deps      # 安装开发依赖（包含 pre-commit hook）
make format            # 格式化代码
make lint              # 运行代码检查
make build             # 构建项目
make test              # 运行测试
make generate          # 生成 proto 代码
make docker            # 启动 Docker 模拟 E2B 环境（用于测试）
make e2b               # 构建 E2B 模板（用于生产部署）
make help              # 查看所有命令
```

## 项目结构

```
agent-sandbox/
├── cmd/server/           # 主程序入口
├── configs/              # 配置文件
├── docs/                 # 文档
├── internal/             # 私有应用代码
│   ├── config/          # 配置管理
│   ├── domain/          # 业务域（core/file/shell）
│   └── middleware/      # 中间件（认证、日志）
├── proto/               # Protocol Buffer 定义
├── scripts/             # 构建脚本和工具
└── sdk/                 # 生成的 SDK 代码
```

## 文档

- [开发指南](docs/development.md) - 详细的开发说明、API 开发、测试和部署
- [API 文档](docs/api.md) - API 接口详细说明（即将推出）

## 许可证

MIT License
