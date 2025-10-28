// Package middleware provides authentication and authorization interceptors.
package middleware

import (
	"context"
	"fmt"
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

var (
	// 跳过认证的路由后缀列表.
	skipAuthSuffixes = []string{
		"/InitSandbox",
	}

	// 错误定义.
	errMissingAPIKey = connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing API key in header %s", APIKeyHeader))
	errInvalidAPIKey = connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid API key"))
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
		// 跳过特定接口的认证
		if i.shouldSkipAuth(req) {
			return next(ctx, req)
		}

		// 执行认证
		sandboxID, err := i.authenticate(ctx, req)
		if err != nil {
			return nil, err
		}

		// 将 Sandbox ID 存入上下文
		ctx = context.WithValue(ctx, SandboxIDKey, sandboxID)

		i.logger.DebugContext(ctx, "authentication successful",
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

// shouldSkipAuth 判断是否需要跳过认证.
func (i *AuthInterceptor) shouldSkipAuth(req connect.AnyRequest) bool {
	procedure := req.Spec().Procedure
	for _, suffix := range skipAuthSuffixes {
		if strings.HasSuffix(procedure, suffix) {
			return true
		}
	}

	return false
}

// authenticate 执行认证逻辑.
func (i *AuthInterceptor) authenticate(ctx context.Context, req connect.AnyRequest) (string, error) {
	procedure := req.Spec().Procedure

	// 从请求头获取 API Key
	apiKey := req.Header().Get(APIKeyHeader)
	if apiKey == "" {
		i.logger.WarnContext(ctx, "authentication failed: missing API key",
			slog.String("procedure", procedure))

		return "", errMissingAPIKey
	}

	// 验证 API Key
	sandboxID, ok := i.store.Verify(apiKey)
	if !ok {
		i.logger.WarnContext(ctx, "authentication failed: invalid API key",
			slog.String("procedure", procedure),
			slog.String("api_key_prefix", maskAPIKey(apiKey)))

		return "", errInvalidAPIKey
	}

	return sandboxID, nil
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
