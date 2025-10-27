# Agent Sandbox Service

A secure sandbox service designed for AI Agents, providing isolated file operations and command execution environment. Built with Go and gRPC/Connect-RPC.

## Features

- **Secure Sandbox Environment**: Each sandbox instance has its own isolated workspace
- **API Key Authentication**: Secure access control using X-Sandbox-Api-Key
- **File Operations**: Support for file reading, writing, and editing
- **Shell Command Execution**: Command execution with timeout control
- **Structured Logging**: JSON format logging using slog
- **Cloud Deployment**: Support for E2B and Docker deployment

## Quick Start

### Prerequisites

- Go 1.24+
- Make
- Docker (optional, for container deployment)

### Local Development

```bash
# Clone the repository
git clone https://github.com/HJH0924/agent-sandbox.git
cd agent-sandbox

# Install development dependencies (recommended)
make install-deps

# Install Go module dependencies
go mod tidy

# Build and run
make build
./bin/api-server --config configs/config.yaml
```

The service will start at `http://localhost:8080`.

## Basic Usage

### 1. Initialize Sandbox

```bash
curl -X POST http://localhost:8080/core.v1.CoreService/InitSandbox \
  -H "Content-Type: application/json" \
  -d '{}'
```

Response example:
```json
{
  "sandboxId": "550e8400-e29b-41d4-a716-446655440000",
  "apiKey": "sk_0123456789abcdef...",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

### 2. Execute Commands

```bash
curl -X POST http://localhost:8080/shell.v1.ShellService/Execute \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"command": "ls -la"}'
```

### 3. File Operations

**Write File:**
```bash
curl -X POST http://localhost:8080/file.v1.FileService/Write \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"path": "example.txt", "content": "Hello, World!"}'
```

**Read File:**
```bash
curl -X POST http://localhost:8080/file.v1.FileService/Read \
  -H "Content-Type: application/json" \
  -H "X-Sandbox-Api-Key: sk_your_api_key" \
  -d '{"path": "example.txt"}'
```

## Local Testing

### Using Docker to Simulate E2B Environment

You can use Docker locally to simulate the E2B sandbox environment for development and testing:

```bash
# Start Docker container (using the same Dockerfile as E2B)
make docker

# View logs
docker-compose logs -f

# Enter container for testing
make shell
```

## Production Deployment

### E2B Sandbox Deployment

**E2B cloud sandbox is recommended for production environments**, specifically designed for AI Agents.

```bash
# Install E2B CLI
npm install -g @e2b/cli

# Set API Key (get from https://e2b.dev/dashboard)
export E2B_API_KEY=your_api_key

# Build E2B template
make e2b

# Create sandbox
e2b sandbox spawn agent-sandbox
```

For detailed E2B usage instructions, see the [Development Guide](docs/development.md#e2b-template-deployment).

## Configuration

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

## Development

See the [Development Guide](docs/development.md) for detailed development instructions.

**Common Commands:**

```bash
make install-deps      # Install development dependencies (includes pre-commit hook)
make format            # Format code
make lint              # Run code linting
make build             # Build project
make test              # Run tests
make generate          # Generate proto code
make docker            # Start Docker to simulate E2B environment (for testing)
make e2b               # Build E2B template (for production deployment)
make help              # View all commands
```

## Project Structure

```
agent-sandbox/
├── cmd/server/           # Main program entry
├── configs/              # Configuration files
├── docs/                 # Documentation
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   ├── domain/          # Business domains (core/file/shell)
│   └── middleware/      # Middleware (authentication, logging)
├── proto/               # Protocol Buffer definitions
├── scripts/             # Build scripts and tools
└── sdk/                 # Generated SDK code
```

## Documentation

- [Development Guide](docs/development.md) - Detailed development instructions, API development, testing and deployment
- [API Documentation](docs/api.md) - Detailed API interface documentation (coming soon)

## License

MIT License

---

[中文文档](README.zh-CN.md)
