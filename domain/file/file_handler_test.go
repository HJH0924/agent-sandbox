package file

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/HJH0924/agent-sandbox/domain/file/service"
	filev1 "github.com/HJH0924/agent-sandbox/sdk/go/file/v1"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testFileName = "test.txt"

func TestNewHandler(t *testing.T) {
	fileService := service.NewService(1024*1024, "/tmp")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	handler := NewHandler(fileService, logger)

	assert.NotNil(t, handler)
	assert.Equal(t, fileService, handler.fileService)
	assert.Equal(t, logger, handler.logger)
}

// ==================== Read Tests ====================

func TestHandler_Read_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := testFileName
	testContent := "Hello, World!"

	// 创建测试文件
	fullPath := filepath.Join(tmpDir, testFile)
	err := os.WriteFile(fullPath, []byte(testContent), 0o600)
	require.NoError(t, err)

	fileService := service.NewService(1024*1024, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(fileService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&filev1.ReadRequest{
		Path: testFile,
	})

	resp, err := handler.Read(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, testContent, resp.Msg.GetContent())
}

func TestHandler_Read_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	fileService := service.NewService(1024*1024, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(fileService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&filev1.ReadRequest{
		Path: "nonexistent.txt",
	})

	resp, err := handler.Read(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	var connectErr *connect.Error
	assert.True(t, errors.As(err, &connectErr))
	assert.Equal(t, connect.CodeInternal, connectErr.Code())
}

func TestHandler_Read_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := "empty.txt"

	// 创建空文件
	fullPath := filepath.Join(tmpDir, testFile)
	err := os.WriteFile(fullPath, []byte(""), 0o600)
	require.NoError(t, err)

	fileService := service.NewService(1024*1024, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(fileService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&filev1.ReadRequest{
		Path: testFile,
	})

	resp, err := handler.Read(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "", resp.Msg.GetContent())
}

// ==================== Write Tests ====================

func TestHandler_Write_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := testFileName
	testContent := "Hello, World!"

	fileService := service.NewService(1024*1024, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(fileService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&filev1.WriteRequest{
		Path:    testFile,
		Content: testContent,
	})

	resp, err := handler.Write(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// 验证文件已创建
	fullPath := filepath.Join(tmpDir, testFile)
	content, readErr := os.ReadFile(fullPath) //nolint:gosec
	require.NoError(t, readErr)
	assert.Equal(t, testContent, string(content))
}

func TestHandler_Write_EmptyContent(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := testFileName

	fileService := service.NewService(1024*1024, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(fileService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&filev1.WriteRequest{
		Path:    testFile,
		Content: "",
	})

	resp, err := handler.Write(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestHandler_Write_CreateNestedDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := "subdir/nested/test.txt"
	testContent := "content"

	fileService := service.NewService(1024*1024, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(fileService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&filev1.WriteRequest{
		Path:    testFile,
		Content: testContent,
	})

	resp, err := handler.Write(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// 验证文件已创建
	fullPath := filepath.Join(tmpDir, testFile)
	content, readErr := os.ReadFile(fullPath) //nolint:gosec
	require.NoError(t, readErr)
	assert.Equal(t, testContent, string(content))
}

// ==================== Edit Tests ====================

func TestHandler_Edit_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := testFileName
	originalContent := "original content"
	newContent := "updated content"

	// 创建原始文件
	fullPath := filepath.Join(tmpDir, testFile)
	err := os.WriteFile(fullPath, []byte(originalContent), 0o600)
	require.NoError(t, err)

	fileService := service.NewService(1024*1024, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(fileService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&filev1.EditRequest{
		Path:    testFile,
		Content: newContent,
	})

	resp, err := handler.Edit(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, testFile, resp.Msg.GetPath())
	assert.Equal(t, newContent, resp.Msg.GetContent())

	// 验证文件已更新
	content, readErr := os.ReadFile(fullPath) //nolint:gosec //nolint:gosec
	require.NoError(t, readErr)
	assert.Equal(t, newContent, string(content))
}

func TestHandler_Edit_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	fileService := service.NewService(1024*1024, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(fileService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&filev1.EditRequest{
		Path:    "nonexistent.txt",
		Content: "content",
	})

	resp, err := handler.Edit(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	var connectErr *connect.Error
	assert.True(t, errors.As(err, &connectErr))
	assert.Equal(t, connect.CodeInternal, connectErr.Code())
}

func TestHandler_Edit_EmptyContent(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := testFileName
	originalContent := "original content"

	// 创建原始文件
	fullPath := filepath.Join(tmpDir, testFile)
	err := os.WriteFile(fullPath, []byte(originalContent), 0o600)
	require.NoError(t, err)

	fileService := service.NewService(1024*1024, tmpDir)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(fileService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&filev1.EditRequest{
		Path:    testFile,
		Content: "",
	})

	resp, err := handler.Edit(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// 验证文件已清空
	content, readErr := os.ReadFile(fullPath) //nolint:gosec //nolint:gosec
	require.NoError(t, readErr)
	assert.Equal(t, "", string(content))
}
