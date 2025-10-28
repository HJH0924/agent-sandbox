// Package router provides HTTP routing configuration for the agent-sandbox API server.
package router

import (
	"log/slog"
	"net/http"

	"github.com/HJH0924/agent-sandbox/domain/core"
	"github.com/HJH0924/agent-sandbox/domain/core/service"
	"github.com/HJH0924/agent-sandbox/domain/file"
	"github.com/HJH0924/agent-sandbox/domain/shell"
	"github.com/HJH0924/agent-sandbox/internal/middleware"
	corev1connect "github.com/HJH0924/agent-sandbox/sdk/go/core/v1/corev1connect"
	filev1connect "github.com/HJH0924/agent-sandbox/sdk/go/file/v1/filev1connect"
	shellv1connect "github.com/HJH0924/agent-sandbox/sdk/go/shell/v1/shellv1connect"

	"connectrpc.com/connect"
)

// Config 路由配置.
type Config struct {
	CoreHandler  *core.Handler
	FileHandler  *file.Handler
	ShellHandler *shell.Handler
	APIKeyStore  service.APIKeyStore
	Logger       *slog.Logger
}

// Setup 设置路由.
func Setup(cfg *Config) http.Handler {
	mux := http.NewServeMux()

	// 创建认证拦截器
	authInterceptor := middleware.NewAuthInterceptor(cfg.APIKeyStore, cfg.Logger)

	// 注册公开路由（不需要认证）
	registerPublicRoutes(mux, cfg)

	// 注册受保护路由（需要认证）
	registerProtectedRoutes(mux, cfg, authInterceptor)

	return mux
}

// registerPublicRoutes 注册不需要认证的路由.
func registerPublicRoutes(mux *http.ServeMux, cfg *Config) {
	// CoreService - 不需要认证
	corePath, coreHandler := corev1connect.NewCoreServiceHandler(cfg.CoreHandler)
	mux.Handle(corePath, coreHandler)

	// 健康检查
	mux.HandleFunc("/health", healthCheckHandler(cfg.Logger))
}

// registerProtectedRoutes 注册需要认证的路由.
func registerProtectedRoutes(mux *http.ServeMux, cfg *Config, authInterceptor connect.Interceptor) {
	// FileService - 需要认证
	filePath, fileHandler := filev1connect.NewFileServiceHandler(
		cfg.FileHandler,
		connect.WithInterceptors(authInterceptor),
	)
	mux.Handle(filePath, fileHandler)

	// ShellService - 需要认证
	shellPath, shellHandler := shellv1connect.NewShellServiceHandler(
		cfg.ShellHandler,
		connect.WithInterceptors(authInterceptor),
	)
	mux.Handle(shellPath, shellHandler)
}

// healthCheckHandler 健康检查处理器.
func healthCheckHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("OK")); err != nil {
			logger.Error("failed to write health check response", slog.Any("error", err))
		}
	}
}
