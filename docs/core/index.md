# 核心服务

核心服务负责沙箱初始化和 API 密钥管理。

## API

### InitSandbox

初始化一个新的沙箱实例并生成 API 密钥。

**端点**: `/core.v1.CoreService/InitSandbox`

**认证**: 不需要

**请求**:
```json
{}
```

**响应**:
```json
{
  "sandboxId": "550e8400-e29b-41d4-a716-446655440000",
  "apiKey": "sk_0123456789abcdef...",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

## 使用示例

```bash
curl -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}'
```

## 实现细节

- 为 sandbox ID 生成 UUID
- 创建一个带有 `sk_` 前缀的 32 字节随机 API key
- 将映射关系存储在内存中（MemoryAPIKeyStore）
- 返回创建时间戳

## 安全性

- 每个沙箱都有唯一的 API key
- 其他所有服务调用都需要 API key
- 密钥由认证中间件验证
