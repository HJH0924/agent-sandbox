// Package middleware provides authentication and authorization interceptors.
package middleware

import (
	"context"
	"log/slog"
	"strings"

	"github.com/HJH0924/agent-sandbox/domain/core/service"

	"connectrpc.com/connect"
)

// contextKey 自定义上下文键类型.
type contextKey string

const (
	// APIKeyHeader API Key 请求头.
	APIKeyHeader = "X-Sandbox-Api-Key" // #nosec G101 -- This is a header name, not a credential
	// SandboxIDKey 上下文中的 Sandbox ID key.
	SandboxIDKey contextKey = "sandbox_id"
)

// AuthInterceptor 认证拦截器.
type AuthInterceptor struct {
	store  service.APIKeyStore
	logger *slog.Logger
}

// NewAuthInterceptor 创建认证拦截器.
func NewAuthInterceptor(store service.APIKeyStore, logger *slog.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		store:  store,
		logger: logger,
	}
}

// WrapUnary 拦截 Unary 调用.
func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		// 跳过 InitSandbox 接口的认证
		if strings.HasSuffix(req.Spec().Procedure, "/InitSandbox") {
			return next(ctx, req)
		}

		// 从请求头获取 API Key
		apiKey := req.Header().Get(APIKeyHeader)
		if apiKey == "" {
			i.logger.WarnContext(ctx, "missing api key",
				slog.String("procedure", req.Spec().Procedure))

			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				connect.NewError(connect.CodeUnauthenticated, nil),
			)
		}

		// 验证 API Key
		sandboxID, ok := i.store.Verify(apiKey)
		if !ok {
			i.logger.WarnContext(ctx, "invalid api key",
				slog.String("procedure", req.Spec().Procedure),
				slog.String("api_key_prefix", maskAPIKey(apiKey)))

			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				connect.NewError(connect.CodeUnauthenticated, nil),
			)
		}

		// 将 Sandbox ID 存入上下文
		ctx = context.WithValue(ctx, SandboxIDKey, sandboxID)

		i.logger.DebugContext(ctx, "api key verified",
			slog.String("procedure", req.Spec().Procedure),
			slog.String("sandbox_id", sandboxID))

		return next(ctx, req)
	}
}

// WrapStreamingClient 拦截流式客户端调用（暂不支持）.
func (i *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler 拦截流式服务端调用（暂不支持）.
func (i *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

// maskAPIKey 遮蔽 API Key，只显示前8个字符.
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "***"
	}

	return apiKey[:8] + "..."
}

// GetSandboxIDFromContext 从上下文中获取 Sandbox ID.
func GetSandboxIDFromContext(ctx context.Context) (string, bool) {
	sandboxID, ok := ctx.Value(SandboxIDKey).(string)
	return sandboxID, ok
}
