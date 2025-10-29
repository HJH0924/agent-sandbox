package core

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/HJH0924/agent-sandbox/domain/core/service"
	corev1 "github.com/HJH0924/agent-sandbox/sdk/go/core/v1"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	apiKeyStore := service.NewMemoryAPIKeyStore()
	coreService := service.NewService(apiKeyStore)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	handler := NewHandler(coreService, logger)

	assert.NotNil(t, handler)
	assert.Equal(t, coreService, handler.coreService)
	assert.Equal(t, logger, handler.logger)
}

func TestHandler_InitSandbox_Success(t *testing.T) {
	apiKeyStore := service.NewMemoryAPIKeyStore()
	coreService := service.NewService(apiKeyStore)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(coreService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&corev1.InitSandboxRequest{})

	resp, err := handler.InitSandbox(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Msg.GetSandboxId())
	assert.NotEmpty(t, resp.Msg.GetApiKey())
	assert.NotNil(t, resp.Msg.GetCreatedAt())

	// 验证API密钥格式
	assert.Contains(t, resp.Msg.GetApiKey(), "sk_")

	// 验证API密钥可以被验证
	sandboxID, ok := apiKeyStore.Verify(resp.Msg.GetApiKey())
	assert.True(t, ok)
	assert.Equal(t, resp.Msg.GetSandboxId(), sandboxID)
}

func TestHandler_InitSandbox_MultipleInvocations(t *testing.T) {
	apiKeyStore := service.NewMemoryAPIKeyStore()
	coreService := service.NewService(apiKeyStore)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(coreService, logger)

	ctx := context.Background()

	// 第一次调用
	req1 := connect.NewRequest(&corev1.InitSandboxRequest{})
	resp1, err1 := handler.InitSandbox(ctx, req1)

	assert.NoError(t, err1)
	assert.NotNil(t, resp1)

	// 第二次调用
	req2 := connect.NewRequest(&corev1.InitSandboxRequest{})
	resp2, err2 := handler.InitSandbox(ctx, req2)

	assert.NoError(t, err2)
	assert.NotNil(t, resp2)

	// 确保返回的是不同的沙箱
	assert.NotEqual(t, resp1.Msg.GetSandboxId(), resp2.Msg.GetSandboxId())
	assert.NotEqual(t, resp1.Msg.GetApiKey(), resp2.Msg.GetApiKey())

	// 验证两个API密钥都有效
	sandboxID1, ok1 := apiKeyStore.Verify(resp1.Msg.GetApiKey())
	assert.True(t, ok1)
	assert.Equal(t, resp1.Msg.GetSandboxId(), sandboxID1)

	sandboxID2, ok2 := apiKeyStore.Verify(resp2.Msg.GetApiKey())
	assert.True(t, ok2)
	assert.Equal(t, resp2.Msg.GetSandboxId(), sandboxID2)
}

func TestHandler_InitSandbox_TimestampIsSet(t *testing.T) {
	apiKeyStore := service.NewMemoryAPIKeyStore()
	coreService := service.NewService(apiKeyStore)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := NewHandler(coreService, logger)

	ctx := context.Background()
	req := connect.NewRequest(&corev1.InitSandboxRequest{})

	resp, err := handler.InitSandbox(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg.GetCreatedAt())

	// 验证时间戳不是零值
	assert.False(t, resp.Msg.GetCreatedAt().AsTime().IsZero())
}
