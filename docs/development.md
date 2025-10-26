# Agent Sandbox 开发指南

## 目录

- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [开发流程](#开发流程)
- [API 开发](#api-开发)
- [测试](#测试)
- [部署](#部署)

## 快速开始

### 前置要求

- **Go 1.23+**: [下载地址](https://go.dev/dl/)
- **Protocol Buffers 编译器**: 安装 buf CLI
  ```bash
  # macOS
  brew install bufbuild/buf/buf

  # 其他系统: https://buf.build/docs/installation
  ```
- **Make**: 用于构建和运行任务

### 安装

```bash
# 克隆仓库
git clone https://github.com/HJH0924/agent-sandbox.git
cd agent-sandbox

# 安装依赖
go mod download

# 生成 proto 代码
make generate

# 构建项目
make build
```

### 运行服务器

```bash
# 使用默认配置运行
./bin/agent-sandbox

# 使用自定义配置运行
./bin/agent-sandbox --config path/to/config.yaml
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
├── internal/                # 私有应用代码
│   ├── config/             # 配置管理
│   ├── domain/             # 业务域
│   │   ├── core/          # 核心服务（沙箱初始化）
│   │   │   ├── service/   # 业务逻辑
│   │   │   └── core_handler.go
│   │   ├── file/          # 文件服务（读/写/编辑）
│   │   │   ├── service/
│   │   │   └── file_handler.go
│   │   └── shell/         # Shell 服务（命令执行）
│   │       ├── serivce/
│   │       └── shell_handler.go
│   └── middleware/         # 中间件（认证、日志）
├── proto/                  # Protocol Buffer 定义
│   ├── core/v1/
│   ├── file/v1/
│   └── shell/v1/
├── sdk/                    # 生成的 SDK 代码（不要手动编辑）
│   └── go/
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── buf.gen.yaml           # Buf 代码生成配置
├── buf.yaml               # Buf 项目配置
└── go.mod
```

### 关键目录说明

- **cmd/**: 应用程序入口
- **internal/**: 私有应用代码
  - **domain/**: 领域驱动设计结构，每个域都是独立的
  - **config/**: 配置加载和管理
  - **middleware/**: HTTP/gRPC 拦截器
- **proto/**: Protocol Buffer 定义（API 的真实来源）
- **sdk/**: 自动生成的代码（永远不要手动编辑）

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

初始化沙箱并生成 API 密钥：

```bash
curl -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}'
```

### 文件服务

在沙箱工作空间内的文件操作：

- **Read**: 获取文件内容
- **Write**: 创建或覆盖文件
- **Edit**: 更新现有文件

### Shell 服务

执行 shell 命令：

- 在沙箱工作空间目录中运行
- 可配置超时
- 捕获 stdout 和 stderr

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
```

## 部署

### Docker

```bash
# 构建镜像
docker build -t agent-sandbox:latest .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/configs/config.yaml:/app/configs/config.yaml \
  agent-sandbox:latest
```

### Docker Compose

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 查看日志
docker-compose logs -f
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
  shell_timeout: 300        # 5 分钟

log:
  level: "info"
  format: "json"
```

## 最佳实践

1. **错误处理**: 始终使用上下文记录错误
2. **日志记录**: 使用 slog 进行结构化日志记录
3. **测试**: 为 service 层编写单元测试
4. **文档**: 添加新功能时更新文档
5. **Proto 变更**: 修改 proto 文件后重新生成 SDK

## Makefile 命令

```bash
make build      # 构建二进制文件
make generate   # 生成 proto 代码
make test       # 运行测试
make lint       # 运行代码检查
make format     # 格式化代码
make clean      # 清理构建产物
```

## 故障排除

### Proto 生成问题

```bash
# 检查 buf 安装
buf --version

# 清理并重新生成
rm -rf sdk/go
make generate
```

### 构建错误

```bash
# 清理 go 缓存
go clean -cache

# 整理依赖
go mod tidy

# 重新构建
make build
```

### 测试失败

```bash
# 运行测试并显示详细输出
go test -v ./...

# 运行特定测试
go test -v -run TestName ./internal/domain/core/service/
```

## 贡献

1. 创建功能分支
2. 进行更改
3. 添加测试
4. 运行 `make lint` 和 `make test`
5. 提交 Pull Request

## 资源

- [Connect RPC 文档](https://connectrpc.com/docs/)
- [Protocol Buffers 指南](https://protobuf.dev/)
- [Go slog 包](https://pkg.go.dev/log/slog)
- [Buf CLI 文档](https://buf.build/docs/)
