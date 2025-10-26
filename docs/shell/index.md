# Shell 服务

Shell 服务提供沙箱内的安全 shell 命令执行。

## API

### Execute

在沙箱工作空间中执行 shell 命令。

**端点**: `/shell.v1.ShellService/Execute`

**认证**: 需要（X-Sandbox-Api-Key 请求头）

**请求**:
```json
{
  "command": "ls -la"
}
```

**响应**:
```json
{
  "output": "total 8\ndrwxr-xr-x  2 user  staff   64 Jan  1 00:00 .\ndrwxr-xr-x  3 user  staff   96 Jan  1 00:00 ..\n-rw-r--r--  1 user  staff   13 Jan  1 00:00 hello.txt\n"
}
```

## 使用示例

```bash
# 列出文件
curl -X POST http://localhost:8080/shell.v1.ShellService/Execute \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_key" \
  -d '{"command": "ls -la"}'

# 创建文件
curl -X POST http://localhost:8080/shell.v1.ShellService/Execute \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_key" \
  -d '{"command": "echo Hello > test.txt"}'

# 检查当前目录
curl -X POST http://localhost:8080/shell.v1.ShellService/Execute \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_key" \
  -d '{"command": "pwd"}'
```

## 功能

- 在沙箱工作空间目录中执行
- 捕获 stdout 和 stderr
- 可配置超时（默认 5 分钟）
- 返回合并的输出

## 限制

- 最大执行时间: 5 分钟（可配置）
- 不支持交互式命令
- 命令以服务器进程权限运行

## 安全性

- 命令在隔离的工作空间目录中执行
- 超时防止长时间运行的进程
- 输出大小限制防止内存问题
- 无法访问工作空间外的系统目录
