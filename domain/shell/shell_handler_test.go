package shell

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/HJH0924/agent-sandbox/domain/shell/service"
	shellv1 "github.com/HJH0924/agent-sandbox/sdk/go/shell/v1"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	shellService := service.NewService(30, "/tmp")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	handler := NewHandler(shellService, logger)

	assert.NotNil(t, handler)
	assert.Equal(t, shellService, handler.shellService)
	assert.Equal(t, logger, handler.logger)
}

func TestHandler_Execute_Success(t *testing.T) {
	tmpDir := t.TempDir()
	shellService := service.NewService(30, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(shellService, logger)

	ctx := context.Background()
	command := "echo hello"

	req := connect.NewRequest(&shellv1.ExecuteRequest{
		Command: command,
	})

	resp, err := handler.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.Msg.GetOutput(), "hello")
}

func TestHandler_Execute_CommandFailed_WithOutput(t *testing.T) {
	tmpDir := t.TempDir()
	shellService := service.NewService(30, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(shellService, logger)

	ctx := context.Background()
	command := "ls /nonexistent_path_12345"

	req := connect.NewRequest(&shellv1.ExecuteRequest{
		Command: command,
	})

	resp, err := handler.Execute(ctx, req)

	// 应该返回输出和错误
	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Msg.GetOutput())

	// 验证错误类型
	var connectErr *connect.Error
	assert.True(t, errors.As(err, &connectErr))
	assert.Equal(t, connect.CodeInternal, connectErr.Code())
}

func TestHandler_Execute_CommandFailed_NoOutput(t *testing.T) {
	tmpDir := t.TempDir()
	shellService := service.NewService(1, tmpDir) // 1秒超时
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(shellService, logger)

	ctx := context.Background()
	command := "sleep 5" // 会超时

	req := connect.NewRequest(&shellv1.ExecuteRequest{
		Command: command,
	})

	_, err := handler.Execute(ctx, req)

	// 超时错误，可能没有输出
	assert.Error(t, err)

	var connectErr *connect.Error
	assert.True(t, errors.As(err, &connectErr))
	assert.Equal(t, connect.CodeInternal, connectErr.Code())
}

func TestHandler_Execute_EmptyCommand(t *testing.T) {
	tmpDir := t.TempDir()
	shellService := service.NewService(30, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(shellService, logger)

	ctx := context.Background()
	command := ""

	req := connect.NewRequest(&shellv1.ExecuteRequest{
		Command: command,
	})

	resp, err := handler.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestHandler_Execute_WorkingDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// 在临时目录创建测试文件
	testFile := "test.txt"
	testContent := "test content"
	err := os.WriteFile(tmpDir+"/"+testFile, []byte(testContent), 0o600)
	require.NoError(t, err)

	shellService := service.NewService(30, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(shellService, logger)

	ctx := context.Background()
	command := "cat " + testFile

	req := connect.NewRequest(&shellv1.ExecuteRequest{
		Command: command,
	})

	resp, err := handler.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.Msg.GetOutput(), testContent)
}
