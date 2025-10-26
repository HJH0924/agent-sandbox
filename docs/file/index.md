# 文件服务

文件服务提供沙箱工作空间内的安全文件操作。

## API

### Read

从沙箱工作空间读取文件内容。

**端点**: `/file.v1.FileService/Read`

**认证**: 需要（X-Sandbox-Api-Key 请求头）

**请求**:
```json
{
  "path": "example.txt"
}
```

**响应**:
```json
{
  "content": "文件内容"
}
```

### Write

向文件写入内容。如需要会自动创建目录。

**端点**: `/file.v1.FileService/Write`

**认证**: 需要

**请求**:
```json
{
  "path": "folder/example.txt",
  "content": "Hello, World!"
}
```

**响应**:
```json
{}
```

### Edit

通过替换内容编辑现有文件。

**端点**: `/file.v1.FileService/Edit`

**认证**: 需要

**请求**:
```json
{
  "path": "example.txt",
  "content": "更新后的内容"
}
```

**响应**:
```json
{
  "path": "example.txt",
  "content": "更新后的内容"
}
```

## 使用示例

```bash
# 写入文件
curl -X POST http://localhost:8080/file.v1.FileService/Write \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_key" \
  -d '{"path": "hello.txt", "content": "Hello"}'

# 读取文件
curl -X POST http://localhost:8080/file.v1.FileService/Read \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_key" \
  -d '{"path": "hello.txt"}'

# 编辑文件
curl -X POST http://localhost:8080/file.v1.FileService/Edit \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_key" \
  -d '{"path": "hello.txt", "content": "Hello, Updated!"}'
```

## 限制

- 最大文件大小: 100MB（可配置）
- 文件限定在沙箱工作空间目录内
- 支持二进制文件，但以 base64 格式返回

## 安全性

- 所有路径都相对于工作空间目录
- 防止路径遍历攻击
- 强制执行文件大小限制
