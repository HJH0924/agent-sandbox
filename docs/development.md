# Agent Sandbox 开发指南

## 目录

- [快速开始](#快速开始)
- [开发环境设置](#开发环境设置)
- [项目结构](#项目结构)
- [开发流程](#开发流程)
- [API 开发](#api-开发)
- [测试](#测试)
- [部署](#部署)
- [故障排除](#故障排除)

## 快速开始

### 前置要求

- **Go 1.24+**: [下载地址](https://go.dev/dl/)
- **Make**: 用于构建和运行任务
- **Git**: 版本控制
- **Docker**（可选）: 用于容器化部署
- **开发工具**（可通过 `make install-deps` 自动安装）:
  - golangci-lint: 代码检查工具
  - gofumpt: Go 代码格式化工具
  - goimports: Import 语句整理工具
  - goimports-reviser: Import 语句排序工具
  - buf: Protocol Buffers 工具

### 安装步骤

```bash
# 1. 克隆仓库
git clone https://github.com/HJH0924/agent-sandbox.git
cd agent-sandbox

# 2. 安装开发依赖（推荐）
make install-deps

# 3. 安装 Go 模块依赖
go mod tidy

# 4. 生成 proto 代码
make generate

# 5. 构建项目
make build

# 6. 运行服务
./bin/api-server --config configs/config.yaml
```

服务启动后将监听 `http://localhost:8080`。

## 开发环境设置

### 安装开发依赖

运行 `make install-deps` 会自动完成以下操作：

1. **检查并安装开发工具**：
   - golangci-lint: 用于代码质量检查
   - gofumpt: 比标准 gofmt 更严格的格式化工具
   - goimports: 自动添加、删除和整理 import 语句
   - goimports-reviser: 按规则排序 import 语句
   - buf: 用于 Protocol Buffers 的编译和管理

2. **安装 Git pre-commit hook**：
   - 自动在每次提交前运行代码格式化（`make format`）
   - 自动运行代码检查（`make lint`）
   - 如果格式化或检查失败，将阻止提交
   - 如果格式化产生了变更，会提示先添加这些变更

### Pre-commit Hook

Pre-commit hook 位于 `scripts/pre-commit`，会在每次 `git commit` 时自动执行。

**手动安装 hook：**
```bash
make install-hooks
```

**跳过 hook（不推荐）：**
```bash
git commit --no-verify -m "your message"
```

**Hook 工作流程：**
1. 运行 `make format` 格式化代码
2. 运行 `make lint` 进行代码检查
3. 检查是否有格式化产生的未提交更改
4. 如果所有检查通过，允许提交；否则阻止提交并提示修复

### 开发工具使用

```bash
# 格式化代码
make format

# 运行代码检查
make lint

# 构建项目（包含 format 和 lint）
make build

# 运行测试
make test

# 生成 proto 代码
make generate
```

## 项目结构

```
agent-sandbox/
├── cmd/
│   └── server/              # 主程序入口
│       └── main.go
├── configs/
│   └── config.yaml          # 配置文件
├── docs/                    # 文档
│   └── development.md       # 本文档
├── internal/                # 私有应用代码
│   ├── config/             # 配置管理
│   │   └── config.go
│   ├── domain/             # 业务域
│   │   ├── core/          # 核心服务（沙箱初始化）
│   │   │   ├── service/   # 业务逻辑
│   │   │   └── core_handler.go
│   │   ├── file/          # 文件服务（读/写/编辑）
│   │   │   ├── service/
│   │   │   └── file_handler.go
│   │   └── shell/         # Shell 服务（命令执行）
│   │       ├── service/
│   │       └── shell_handler.go
│   └── middleware/         # 中间件（认证、日志）
│       └── auth.go
├── proto/                  # Protocol Buffer 定义
│   ├── core/v1/
│   ├── file/v1/
│   └── shell/v1/
├── scripts/                # 构建脚本和工具
│   └── pre-commit         # Git pre-commit hook
├── sdk/                    # 生成的 SDK 代码（不要手动编辑）
│   └── go/
├── e2b.Dockerfile         # E2B 和 Docker 统一的 Dockerfile
├── docker-compose.yml     # Docker Compose 配置
├── Makefile              # 构建任务
├── buf.gen.yaml          # Buf 代码生成配置
├── buf.yaml              # Buf 项目配置
└── go.mod
```

### 关键目录说明

- **cmd/**: 应用程序入口点
- **internal/**: 私有应用代码，不对外暴露
  - **domain/**: 领域驱动设计结构，每个域独立
  - **config/**: 配置加载和管理
  - **middleware/**: HTTP/gRPC 拦截器
- **proto/**: Protocol Buffer 定义（API 的真实来源）
- **sdk/**: 自动生成的代码（永远不要手动编辑）
- **scripts/**: 构建脚本和开发工具

## 开发流程

### 1. 添加新 API

#### 步骤 1: 定义 Proto 接口

在 `proto/<service>/v1/` 中编辑或创建 `.proto` 文件：

```protobuf
syntax = "proto3";

package myservice.v1;

service MyService {
  rpc DoSomething(DoSomethingRequest) returns (DoSomethingResponse) {}
}

message DoSomethingRequest {
  string param = 1;
}

message DoSomethingResponse {
  string result = 1;
}
```

#### 步骤 2: 生成代码

```bash
make generate
```

这将生成：
- `sdk/go/myservice/v1/myservice.pb.go` - 消息定义
- `sdk/go/myservice/v1/myservicev1connect/myservice.connect.go` - 服务接口

#### 步骤 3: 实现 Service 层

创建 `internal/domain/myservice/service/myservice.go`：

```go
package service

// Service 我的服务
type Service struct {
    // 依赖项
}

// NewService 创建服务实例
func NewService() *Service {
    return &Service{}
}

// DoSomething 执行某个操作
func (s *Service) DoSomething(param string) (string, error) {
    // 业务逻辑
    return "result", nil
}
```

#### 步骤 4: 实现 Handler 层

创建 `internal/domain/myservice/myservice_handler.go`：

```go
package myservice

import (
    "context"
    "log/slog"

    "connectrpc.com/connect"
    myservicev1 "github.com/HJH0924/agent-sandbox/sdk/go/myservice/v1"
    "github.com/HJH0924/agent-sandbox/internal/domain/myservice/service"
)

// Handler 处理器
type Handler struct {
    service *service.Service
    logger  *slog.Logger
}

// NewHandler 创建处理器
func NewHandler(svc *service.Service, logger *slog.Logger) *Handler {
    return &Handler{
        service: svc,
        logger:  logger,
    }
}

// DoSomething 处理请求
func (h *Handler) DoSomething(
    ctx context.Context,
    req *connect.Request[myservicev1.DoSomethingRequest],
) (*connect.Response[myservicev1.DoSomethingResponse], error) {
    h.logger.InfoContext(ctx, "processing request",
        slog.String("param", req.Msg.GetParam()))

    result, err := h.service.DoSomething(req.Msg.GetParam())
    if err != nil {
        h.logger.ErrorContext(ctx, "failed to process",
            slog.Any("error", err))
        return nil, connect.NewError(connect.CodeInternal, err)
    }

    return connect.NewResponse(&myservicev1.DoSomethingResponse{
        Result: result,
    }), nil
}
```

#### 步骤 5: 在 main.go 中注册

```go
// 创建服务和处理器
myService := myservice.NewService()
myHandler := myservice.NewHandler(myService, logger)

// 注册（可选是否需要认证）
path, handler := myservicev1connect.NewMyServiceHandler(
    myHandler,
    connect.WithInterceptors(authInterceptor), // 如需认证则添加
)
mux.Handle(path, handler)
```

### 2. 身份认证

- **InitSandbox**（核心服务）: 不需要认证
- **所有其他 API**: 需要 `X-Sandbox-Api-Key` 请求头

认证由 `middleware.AuthInterceptor` 处理：
- 跳过 `/InitSandbox` 端点的认证
- 验证所有其他请求的 API key
- 将 sandbox ID 存储在上下文中供处理器使用

### 3. 日志记录

使用 `slog` 进行结构化日志记录：

```go
// Info 级别
logger.InfoContext(ctx, "message",
    slog.String("key", "value"),
    slog.Int("count", 10))

// Error 级别
logger.ErrorContext(ctx, "error occurred",
    slog.Any("error", err),
    slog.String("detail", "additional info"))

// Debug 级别
logger.DebugContext(ctx, "debug info",
    slog.String("data", data))
```

## API 开发

### 核心服务

初始化沙箱并生成 API 密钥。

**端点**: `/core.v1.CoreService/InitSandbox`

**请求**: 空

**响应**:
```json
{
  "sandboxId": "uuid",
  "apiKey": "sk_...",
  "createdAt": "timestamp"
}
```

### 文件服务

在沙箱工作空间内的文件操作。

**操作**:
- **Read**: 获取文件内容
- **Write**: 创建或覆盖文件
- **Edit**: 更新现有文件

**认证**: 需要 `X-Sandbox-Api-Key`

### Shell 服务

执行 shell 命令。

**特性**:
- 在沙箱工作空间目录中运行
- 可配置超时
- 捕获 stdout 和 stderr

**认证**: 需要 `X-Sandbox-Api-Key`

## 测试

### 单元测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/domain/core/service/...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行测试并显示详细输出
go test -v ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 集成测试

使用 curl 或任何 HTTP 客户端：

```bash
# 初始化沙箱
RESPONSE=$(curl -s -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}')

# 提取 API key
API_KEY=$(echo $RESPONSE | jq -r '.apiKey')

# 使用 API key 进行文件操作
curl -X POST http://localhost:8080/file.v1.FileService/Write \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: $API_KEY" \
  -d '{"path": "test.txt", "content": "Hello"}'

# 执行命令
curl -X POST http://localhost:8080/shell.v1.ShellService/Execute \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: $API_KEY" \
  -d '{"command": "cat test.txt"}'
```

## 部署

### 本地测试环境（Docker）

**重要**：Docker 不是用于生产部署，而是用于本地测试和开发，模拟 E2B sandbox 环境。

本地 Docker 与 E2B 模板使用相同的 `e2b.Dockerfile`，确保环境一致性。

**使用 Docker 进行本地测试的优势**：
- ✅ **环境一致性**：本地测试环境与 E2B 云端环境完全一致
- ✅ **快速调试**：在本地复现和调试 E2B 环境问题
- ✅ **简化维护**：只需维护一个 Dockerfile
- ✅ **降低风险**：本地测试通过后可以确保 E2B 部署成功

#### 使用 Make（推荐）

```bash
# 构建并启动容器
make docker

# 输出示例：
# Building static binary for E2B...
# Binary built: bin/api-server
# Stopping existing container...
# Building and starting Docker container with docker-compose...
# ✅ Docker container started successfully!

# 进入容器 shell
make shell

# 查看日志
docker-compose logs -f

# 停止并删除容器
docker-compose down
```

#### 推荐的开发工作流

```bash
# 1. 本地开发和测试
make build
make run

# 2. 使用 Docker 模拟 E2B 环境进行测试
make docker

# 3. 测试 API
curl -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}'

# 4. 进入容器调试（如果需要）
make shell

# 5. 本地测试通过后，构建 E2B 模板用于生产环境
make e2b
```

### 生产部署（E2B）

**推荐生产环境使用 E2B 云端沙箱**。[E2B](https://e2b.dev) 是专为 AI Agent 设计的云端代码执行环境。

#### E2B 简介

E2B 是一个专为 AI Agent 设计的安全云端代码执行环境，提供以下特性：
- **即时沙箱创建**：几秒内创建全新的隔离环境
- **资源隔离**：每个沙箱独立的 CPU、内存和文件系统
- **API 友好**：简单的 SDK 集成，支持 TypeScript/JavaScript 和 Python
- **自动扩展**：按需创建和销毁沙箱，无需管理服务器

#### 前置要求

```bash
# 1. 安装 E2B CLI
npm install -g @e2b/cli

# 2. 获取 E2B API Key
# 访问 https://e2b.dev/dashboard 创建 API Key

# 3. 设置环境变量
export E2B_API_KEY=your_api_key
```

#### 构建 E2B 模板

项目采用本地编译方式，先在本地构建 Linux 静态二进制文件，然后使用 E2B CLI 构建模板。

**构建命令**：

```bash
# 使用 Make（推荐）
make e2b

# 输出示例：
# Building static binary for E2B...
# Binary built: bin/api-server
# Building E2B sandbox template...
# ✅ Template built successfully
# Template: agent-sandbox
```

**构建流程**：

1. **本地构建 Linux 静态二进制**（`make build-linux`）：
   - 交叉编译 Linux AMD64 静态二进制文件
   - 使用 CGO_ENABLED=0 确保静态链接
   - 使用优化的编译参数（-ldflags='-w -s'）减小体积
   - 输出到 `bin/api-server`

2. **使用 E2B CLI 构建模板**：
   - 使用 `e2b.Dockerfile` 定义容器环境
   - 将预构建的 Linux 二进制复制到容器中
   - 配置资源（CPU: 2核，内存: 2048MB）和启动命令

**优势**：
- ✅ 显著减少 E2B 构建时间（无需安装 Go 工具链）
- ✅ 减小最终镜像体积（不包含 Go 编译器和构建工具）
- ✅ 提高构建效率和可靠性
- ✅ 便于本地调试和测试

#### 使用 E2B CLI 管理沙箱

**创建并连接沙箱**：

```bash
# 创建并自动连接到沙箱（进入交互式终端）
e2b sandbox spawn agent-sandbox

# 或使用短命令
e2b sbx sp agent-sandbox

# 在沙箱终端中测试
ps aux | grep api-server          # 检查服务是否运行
cat /tmp/agent-sandbox/*.log      # 查看服务日志

# 测试 API
curl -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}'

# 退出沙箱终端
exit  # 或按 Ctrl+D
```

**管理沙箱**：

```bash
# 列出所有运行中的沙箱
e2b sandbox list
# 或短命令: e2b sbx ls

# 连接到已存在的沙箱
e2b sandbox connect <sandbox-id>
# 或短命令: e2b sbx cn <sandbox-id>

# 查看沙箱日志
e2b sandbox logs <sandbox-id>
# 或短命令: e2b sbx lg <sandbox-id>

# 终止沙箱
e2b sandbox kill <sandbox-id>
# 或短命令: e2b sbx kl <sandbox-id>

# 终止所有沙箱
e2b sandbox kill --all
```

#### 在代码中使用 E2B 模板

**TypeScript/JavaScript**：

```typescript
import { Sandbox } from '@e2b/code-interpreter'

async function main() {
  // 创建沙箱
  const sandbox = await Sandbox.create('agent-sandbox')
  console.log('Sandbox ID:', sandbox.sandboxId)

  try {
    // 初始化 Agent Sandbox
    const initResponse = await fetch('http://localhost:8080/core.v1.CoreService/InitSandbox', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({})
    })
    const { apiKey } = await initResponse.json()

    // 执行命令
    const execResponse = await fetch('http://localhost:8080/shell.v1.ShellService/Execute', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Sandbox-Api-Key': apiKey
      },
      body: JSON.stringify({ command: 'ls -la' })
    })
    const { output } = await execResponse.json()
    console.log('Output:', output)

  } finally {
    await sandbox.close()
  }
}
```

**Python**：

```python
from e2b import Sandbox
import requests

def main():
    # 创建沙箱
    sandbox = Sandbox(template='agent-sandbox')
    print(f'Sandbox ID: {sandbox.sandbox_id}')

    try:
        # 初始化 Agent Sandbox
        init_response = requests.post(
            'http://localhost:8080/core.v1.CoreService/InitSandbox',
            json={}
        )
        api_key = init_response.json()['apiKey']

        # 执行命令
        exec_response = requests.post(
            'http://localhost:8080/shell.v1.ShellService/Execute',
            headers={'X-Sandbox-Api-Key': api_key},
            json={'command': 'ls -la'}
        )
        print('Output:', exec_response.json()['output'])

    finally:
        sandbox.close()
```

#### 模板自定义

**1. 本地构建配置（Makefile）**：

修改编译参数或目标平台：

```makefile
# 构建 ARM64 版本
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build ...
```

**2. Dockerfile 配置（e2b.Dockerfile）**：

修改基础镜像、依赖或环境变量：

```dockerfile
FROM ubuntu:22.04

# 安装额外依赖
RUN apt-get update && \
    apt-get install -y ca-certificates curl vim && \
    rm -rf /var/lib/apt/lists/*

# 设置环境变量
ENV LOG_LEVEL=debug
ENV CUSTOM_VAR=value
```

**注意事项**：
- ⚠️ 不要在 Dockerfile 中添加注释（`#` 开头的行），E2B 不支持
- ⚠️ 不能使用多阶段构建（`FROM ... AS ...`）
- ⚠️ 必须使用 Debian 系基础镜像（Ubuntu、Debian）
- ⚠️ 确保先运行 `make build-linux` 构建二进制文件

**3. E2B CLI 参数（Makefile）**：

调整资源配置：

```makefile
e2b: build-linux
	e2b template build \
		--dockerfile e2b.Dockerfile \
		--name agent-sandbox \
		--cpu-count 4 \              # 调整 CPU
		--memory-mb 4096 \           # 调整内存
		--cmd "/app/bin/api-server -c /app/configs/config.yaml"
```

**应用更改**：

```bash
# 重新构建模板
make e2b

# 测试新模板
e2b sandbox spawn agent-sandbox
```

#### E2B vs 本地 Docker

| 特性 | E2B（生产环境） | 本地 Docker（测试环境） |
|------|-----|-------------|
| **用途** | 生产部署 | 本地测试和开发 |
| **部署速度** | 几秒钟创建沙箱 | 需要启动容器 |
| **隔离性** | 云端完全隔离 | 本地容器隔离 |
| **扩展性** | 自动扩展，无限沙箱 | 单个容器 |
| **成本** | 按使用付费 | 无额外成本 |
| **网络** | 公网可访问 | 本地访问 |
| **适用场景** | AI Agent 生产环境 | 本地开发和测试 |

#### 最佳实践

1. **版本管理**：
   ```bash
   # 创建带版本号的模板
   e2b template build --name agent-sandbox-v1.0.0 ...
   ```

2. **资源优化**：
   - 根据实际负载调整 CPU 和内存
   - 生产环境可按需增加资源配置

3. **安全考虑**：
   - 不要在构建脚本中硬编码敏感信息
   - 使用环境变量传递 API Key
   - 定期更新基础镜像和依赖

4. **构建优化**：
   - 本地编译二进制避免在 E2B 中安装 Go
   - 使用 `.dockerignore` 排除不必要的文件
   - 静态链接编译确保二进制独立运行
   - 最小化依赖减小镜像大小

## 故障排除

### 开发工具缺失

**问题**：
```
command not found: golangci-lint
command not found: gofumpt
```

**解决方法**：
```bash
# 自动安装所有开发依赖
make install-deps

# 验证安装
golangci-lint --version
gofumpt -version
```

### Proto 生成问题

**问题**：`make generate` 失败

**解决方法**：
```bash
# 检查 buf 安装
buf --version

# 清理并重新生成
rm -rf sdk/go
make generate
```

### 构建错误

**问题**：编译失败或依赖问题

**解决方法**：
```bash
# 清理缓存
go clean -cache

# 整理依赖
go mod tidy

# 重新构建
make build
```

### Pre-commit Hook 问题

**问题**：Hook 未生效或提交被阻止

**解决方法**：
```bash
# 重新安装 hook
make install-hooks

# 验证 hook 文件存在
ls -la .git/hooks/pre-commit

# 手动运行检查
make format
make lint

# 如需临时跳过 hook（不推荐）
git commit --no-verify -m "message"
```

### Docker 构建问题

**问题**：Docker 镜像构建失败

**解决方法**：
```bash
# 确保二进制已构建
make build-linux
ls -la bin/api-server

# 清理 Docker 缓存
docker system prune -a

# 重新构建
make docker
```

### E2B 构建失败

**问题 1**：找不到二进制文件

**解决方法**：
```bash
# 检查二进制是否存在
ls -la bin/api-server

# 如果不存在，先构建
make build-linux

# 验证二进制
file bin/api-server
# 应输出：bin/api-server: ELF 64-bit LSB executable, x86-64...
```

**问题 2**：E2B API Key 未设置

**解决方法**：
```bash
# 设置 API Key
export E2B_API_KEY=your_api_key

# 验证
echo $E2B_API_KEY
```

**问题 3**：不支持的基础镜像

E2B 只支持 Debian 系镜像（Ubuntu、Debian）。

**支持的镜像**：
```dockerfile
FROM ubuntu:22.04
FROM ubuntu:24.04
FROM debian:bookworm
```

**不支持的镜像**：
```dockerfile
FROM alpine:latest
FROM golang:1.24-alpine
FROM centos:7
```

**问题 4**：Dockerfile 注释错误

E2B 不支持 Dockerfile 中的注释。移除所有以 `#` 开头的注释行。

**问题 5**：服务未运行

进入沙箱后发现服务未运行：

```bash
# 检查进程
ps aux | grep api-server

# 手动启动
/app/bin/api-server -c /app/configs/config.yaml

# 检查日志
ls -la /tmp/agent-sandbox/

# 验证二进制可执行
ls -la /app/bin/api-server
/app/bin/api-server --help
```

### 测试失败

**问题**：测试运行失败

**解决方法**：
```bash
# 运行测试并显示详细输出
go test -v ./...

# 运行特定测试
go test -v -run TestName ./internal/domain/core/service/

# 清理测试缓存
go clean -testcache
go test ./...
```

## Makefile 命令参考

```bash
make help              # 显示所有可用命令
make install-deps      # 安装开发依赖和 pre-commit hook
make install-hooks     # 仅安装 Git pre-commit hooks
make format            # 格式化代码
make lint              # 运行代码检查
make build             # 构建本地二进制（包含 format 和 lint）
make build-linux       # 构建 Linux 静态二进制（用于 Docker 和 E2B）
make run               # 运行服务器
make test              # 运行测试并生成覆盖率报告
make clean             # 清理构建产物
make generate          # 生成 proto 代码
make docs              # 运行文档服务器
make docker            # 启动 Docker 模拟 E2B 环境（用于测试）
make shell             # 进入 Docker 容器 shell
make e2b               # 构建 E2B 模板（用于生产部署）
```

## 贡献

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
   - Pre-commit hook 会自动运行格式化和检查
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 资源

- [Connect RPC 文档](https://connectrpc.com/docs/)
- [Protocol Buffers 指南](https://protobuf.dev/)
- [Go slog 包](https://pkg.go.dev/log/slog)
- [Buf CLI 文档](https://buf.build/docs/)
- [E2B 文档](https://e2b.dev/docs)
- [golangci-lint 文档](https://golangci-lint.run/)
