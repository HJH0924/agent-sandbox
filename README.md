# Agent Sandbox Service

[English](#english) | [中文](#中文)

## English

### Overview

Agent Sandbox is a secure sandbox service designed for AI Agents to safely execute shell commands and perform file operations. Built with Go and gRPC/Connect-RPC, it provides isolated environments with fine-grained access control.

### Features

- **Secure Sandbox Environment**: Isolated workspace for each sandbox instance
- **API Key Authentication**: Secure access control using X-SANDBOX-API-KEY header
- **File Operations**: Read, write, and edit files within the sandbox
- **Shell Command Execution**: Execute shell commands with timeout control
- **JSON Logging**: Structured logging using slog
- **Health Checks**: Built-in health check endpoint
- **Docker Support**: Easy deployment with Docker and Docker Compose

### Architecture

- **Core Service**: Initialize sandbox and generate API keys
- **File Service**: File read/write/edit operations
- **Shell Service**: Execute shell commands with timeout
- **Middleware**: API key authentication interceptor

### Quick Start

#### Prerequisites

- Go 1.23+
- Protocol Buffers compiler (for development)
- Docker and Docker Compose (for deployment)

#### Local Development

```bash
# Clone the repository
git clone https://github.com/HJH0924/agent-sandbox.git
cd agent-sandbox

# Install dependencies
go mod download

# Build the project
make build

# Run the server
./bin/agent-sandbox --config configs/config.yaml
```

#### Docker Deployment

```bash
# Build and start with Docker Compose
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f agent-sandbox
```

### API Usage

#### 1. Initialize Sandbox

```bash
curl -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}'
```

Response:
```json
{
  "sandboxId": "550e8400-e29b-41d4-a716-446655440000",
  "apiKey": "sk_0123456789abcdef...",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

#### 2. File Operations

**Read File:**
```bash
curl -X POST http://localhost:8080/file.v1.FileService/Read \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"path": "example.txt"}'
```

**Write File:**
```bash
curl -X POST http://localhost:8080/file.v1.FileService/Write \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"path": "example.txt", "content": "Hello, World!"}'
```

#### 3. Execute Shell Command

```bash
curl -X POST http://localhost:8080/shell.v1.ShellService/Execute \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"command": "ls -la"}'
```

### Configuration

Edit `configs/config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

sandbox:
  workspace_dir: "/tmp/agent-sandbox"
  max_file_size: 104857600  # 100MB
  shell_timeout: 300        # 5 minutes

log:
  level: "info"  # debug, info, warn, error
  format: "json" # json, text
```

### Development

See [Development Guide](docs/development.md) for detailed development instructions.

```bash
# Generate proto code
make generate

# Run tests
make test

# Run linter
make lint

# Format code
make format
```

### Project Structure

```
agent-sandbox/
├── cmd/
│   └── server/           # Main application entry
├── configs/              # Configuration files
├── docs/                 # Documentation
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Business domains
│   │   ├── core/       # Core service
│   │   ├── file/       # File service
│   │   └── shell/      # Shell service
│   └── middleware/      # Middleware (auth, logging)
├── proto/               # Protocol Buffer definitions
└── sdk/                 # Generated SDK code
```

### License

MIT License

---

## 中文

### 概述

Agent Sandbox 是一个为 AI Agent 设计的安全沙箱服务，用于安全地执行 shell 命令和文件操作。基于 Go 和 gRPC/Connect-RPC 构建，提供隔离环境和细粒度的访问控制。

### 特性

- **安全的沙箱环境**：每个沙箱实例拥有独立的工作空间
- **API Key 认证**：使用 X-SANDBOX-API-KEY 请求头进行安全访问控制
- **文件操作**：在沙箱内读取、写入和编辑文件
- **Shell 命令执行**：支持超时控制的 shell 命令执行
- **JSON 日志**：使用 slog 进行结构化日志记录
- **健康检查**：内置健康检查端点
- **Docker 支持**：使用 Docker 和 Docker Compose 轻松部署

### 架构

- **Core Service**：初始化沙箱并生成 API 密钥
- **File Service**：文件读写编辑操作
- **Shell Service**：执行带超时控制的 shell 命令
- **Middleware**：API 密钥认证拦截器

### 快速开始

#### 前置要求

- Go 1.23+
- Protocol Buffers 编译器（开发用）
- Docker 和 Docker Compose（部署用）

#### 本地开发

```bash
# 克隆仓库
git clone https://github.com/HJH0924/agent-sandbox.git
cd agent-sandbox

# 安装依赖
go mod download

# 构建项目
make build

# 运行服务
./bin/agent-sandbox --config configs/config.yaml
```

#### Docker 部署

```bash
# 使用 Docker Compose 构建并启动
docker-compose up -d

# 检查服务状态
docker-compose ps

# 查看日志
docker-compose logs -f agent-sandbox
```

### API 使用

#### 1. 初始化沙箱

```bash
curl -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}'
```

响应：
```json
{
  "sandboxId": "550e8400-e29b-41d4-a716-446655440000",
  "apiKey": "sk_0123456789abcdef...",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

#### 2. 文件操作

**读取文件：**
```bash
curl -X POST http://localhost:8080/file.v1.FileService/Read \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"path": "example.txt"}'
```

**写入文件：**
```bash
curl -X POST http://localhost:8080/file.v1.FileService/Write \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"path": "example.txt", "content": "Hello, World!"}'
```

#### 3. 执行 Shell 命令

```bash
curl -X POST http://localhost:8080/shell.v1.ShellService/Execute \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"command": "ls -la"}'
```

### 配置

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

### 开发

查看[开发指南](docs/development.md)了解详细的开发说明。

```bash
# 生成 proto 代码
make generate

# 运行测试
make test

# 运行代码检查
make lint

# 格式化代码
make format
```

### 项目结构

```
agent-sandbox/
├── cmd/
│   └── server/           # 主程序入口
├── configs/              # 配置文件
├── docs/                 # 文档
├── internal/
│   ├── config/          # 配置管理
│   ├── domain/          # 业务域
│   │   ├── core/       # 核心服务
│   │   ├── file/       # 文件服务
│   │   └── shell/      # Shell服务
│   └── middleware/      # 中间件（认证、日志）
├── proto/               # Protocol Buffer 定义
└── sdk/                 # 生成的 SDK 代码
```

### 许可证

MIT License
