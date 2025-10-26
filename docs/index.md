---
layout: home

hero:
  name: "Agent Sandbox"
  text: "为 AI Agent 提供安全的沙箱服务"
  tagline: 在隔离环境中安全地执行 shell 命令和文件操作
  actions:
    - theme: brand
      text: 快速开始
      link: /development
    - theme: alt
      text: GitHub
      link: https://github.com/HJH0924/agent-sandbox

features:
  - title: 🔒 安全的沙箱环境
    details: 为每个沙箱实例提供隔离的工作空间，使用 API Key 进行身份验证
  - title: 📁 文件操作
    details: 在沙箱工作空间内安全地读取、写入和编辑文件
  - title: 🖥️ Shell 命令执行
    details: 执行 shell 命令，支持超时控制和输出捕获
  - title: 📊 结构化日志
    details: 使用 slog 记录 JSON 格式日志，便于观察和调试
  - title: 🔑 API Key 认证
    details: 使用 X-SANDBOX-API-KEY 请求头进行安全的访问控制
  - title: 🐳 Docker 支持
    details: 使用 Docker 和 Docker Compose 轻松部署
---

## 快速开始

### 安装和运行

```bash
# 克隆仓库
git clone https://github.com/HJH0924/agent-sandbox.git
cd agent-sandbox

# 构建并运行
make build
./bin/agent-sandbox
```

### 或使用 Docker

```bash
docker-compose up -d
```

## API 示例

### 初始化沙箱

```bash
curl -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}'
```

### 文件操作

```bash
# 写入文件
curl -X POST http://localhost:8080/file.v1.FileService/Write \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: YOUR_API_KEY" \
  -d '{"path": "hello.txt", "content": "Hello, World!"}'
```

### 执行命令

```bash
# 执行 shell 命令
curl -X POST http://localhost:8080/shell.v1.ShellService/Execute \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: YOUR_API_KEY" \
  -d '{"command": "ls -la"}'
```

## 功能特性

### 核心服务
初始化沙箱实例并生成安全的 API 密钥用于身份验证。

### 文件服务
在沙箱工作空间内执行文件操作，包括读取、写入和编辑。

### Shell 服务
执行 shell 命令，支持可配置的超时和输出捕获。

## 了解更多

- [开发指南](/development)
- [核心服务](/core/index)
- [文件服务](/file/index)
- [Shell 服务](/shell/index)
