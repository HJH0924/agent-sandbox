package middleware

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/HJH0924/agent-sandbox/domain/core/service"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthInterceptor(t *testing.T) {
	store := service.NewMemoryAPIKeyStore()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	interceptor := NewAuthInterceptor(store, logger)

	assert.NotNil(t, interceptor)
	assert.Equal(t, store, interceptor.store)
	assert.Equal(t, logger, interceptor.logger)
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		expected string
	}{
		{
			name:     "Long API key",
			apiKey:   "1234567890abcdef",
			expected: "12345678...",
		},
		{
			name:     "Short API key",
			apiKey:   "short",
			expected: "***",
		},
		{
			name:     "Exactly 8 characters",
			apiKey:   "12345678",
			expected: "***",
		},
		{
			name:     "Empty API key",
			apiKey:   "",
			expected: "***",
		},
		{
			name:     "9 characters",
			apiKey:   "123456789",
			expected: "12345678...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskAPIKey(tt.apiKey)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetSandboxIDFromContext(t *testing.T) {
	t.Run("Sandbox ID exists in context", func(t *testing.T) {
		expectedID := "sandbox-12345"
		ctx := context.WithValue(context.Background(), SandboxIDKey, expectedID)

		id, ok := GetSandboxIDFromContext(ctx)

		assert.True(t, ok)
		assert.Equal(t, expectedID, id)
	})

	t.Run("Sandbox ID not in context", func(t *testing.T) {
		ctx := context.Background()

		id, ok := GetSandboxIDFromContext(ctx)

		assert.False(t, ok)
		assert.Equal(t, "", id)
	})

	t.Run("Wrong type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), SandboxIDKey, 12345) // 存入 int 而不是 string

		id, ok := GetSandboxIDFromContext(ctx)

		assert.False(t, ok)
		assert.Equal(t, "", id)
	})
}

func TestAPIKeyStore_Integration(t *testing.T) {
	// 测试 APIKeyStore 的集成
	store := service.NewMemoryAPIKeyStore()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	interceptor := NewAuthInterceptor(store, logger)

	assert.NotNil(t, interceptor)

	// 测试存储和验证
	apiKey := "test-api-key-12345"
	sandboxID := "sandbox-67890"
	err := store.Store(sandboxID, apiKey)
	assert.NoError(t, err)

	// 验证有效的API密钥
	verifiedID, ok := store.Verify(apiKey)
	assert.True(t, ok)
	assert.Equal(t, sandboxID, verifiedID)

	// 验证无效的API密钥
	_, ok = store.Verify("invalid-key")
	assert.False(t, ok)

	// 删除API密钥
	err = store.Delete(sandboxID)
	assert.NoError(t, err)

	// 验证删除后的API密钥
	_, ok = store.Verify(apiKey)
	assert.False(t, ok)
}

func TestErrorMessages(t *testing.T) {
	// 测试错误消息定义
	assert.NotNil(t, errMissingAPIKey)
	assert.NotNil(t, errInvalidAPIKey)

	assert.Contains(t, errMissingAPIKey.Message(), APIKeyHeader)
	assert.Contains(t, errInvalidAPIKey.Message(), "invalid")
}

func TestSkipAuthSuffixes(t *testing.T) {
	// 验证跳过认证的后缀列表
	assert.NotEmpty(t, skipAuthSuffixes)
	assert.Contains(t, skipAuthSuffixes, "/InitSandbox")
}

func TestContextKey(t *testing.T) {
	// 验证上下文键定义
	assert.Equal(t, contextKey("sandbox_id"), SandboxIDKey)
}
