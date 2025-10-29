package router

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/HJH0924/agent-sandbox/domain/core"
	coreservice "github.com/HJH0924/agent-sandbox/domain/core/service"
	"github.com/HJH0924/agent-sandbox/domain/file"
	fileservice "github.com/HJH0924/agent-sandbox/domain/file/service"
	"github.com/HJH0924/agent-sandbox/domain/shell"
	shellservice "github.com/HJH0924/agent-sandbox/domain/shell/service"
	"github.com/HJH0924/agent-sandbox/internal/middleware"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	// 创建测试依赖
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 创建 API Key Store
	apiKeyStore := coreservice.NewMemoryAPIKeyStore()

	// 创建 services
	coreService := coreservice.NewService(apiKeyStore)
	fileService := fileservice.NewService(1024*1024, "/tmp/test-sandbox")
	shellService := shellservice.NewService(30, "/tmp/test-sandbox")

	// 创建 handlers
	coreHandler := core.NewHandler(coreService, logger)
	fileHandler := file.NewHandler(fileService, logger)
	shellHandler := shell.NewHandler(shellService, logger)

	cfg := &Config{
		CoreHandler:  coreHandler,
		FileHandler:  fileHandler,
		ShellHandler: shellHandler,
		APIKeyStore:  apiKeyStore,
		Logger:       logger,
	}

	handler := Setup(cfg)

	assert.NotNil(t, handler)
}

func TestHealthCheckHandler(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := healthCheckHandler(logger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestHealthCheckHandler_Integration(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 创建完整的路由器
	apiKeyStore := coreservice.NewMemoryAPIKeyStore()
	coreService := coreservice.NewService(apiKeyStore)
	fileService := fileservice.NewService(1024*1024, "/tmp/test-sandbox")
	shellService := shellservice.NewService(30, "/tmp/test-sandbox")

	coreHandler := core.NewHandler(coreService, logger)
	fileHandler := file.NewHandler(fileService, logger)
	shellHandler := shell.NewHandler(shellService, logger)

	cfg := &Config{
		CoreHandler:  coreHandler,
		FileHandler:  fileHandler,
		ShellHandler: shellHandler,
		APIKeyStore:  apiKeyStore,
		Logger:       logger,
	}

	mux := Setup(cfg)

	// 测试健康检查端点
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestRegisterPublicRoutes(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	apiKeyStore := coreservice.NewMemoryAPIKeyStore()
	coreService := coreservice.NewService(apiKeyStore)
	fileService := fileservice.NewService(1024*1024, "/tmp/test-sandbox")
	shellService := shellservice.NewService(30, "/tmp/test-sandbox")

	coreHandler := core.NewHandler(coreService, logger)
	fileHandler := file.NewHandler(fileService, logger)
	shellHandler := shell.NewHandler(shellService, logger)

	cfg := &Config{
		CoreHandler:  coreHandler,
		FileHandler:  fileHandler,
		ShellHandler: shellHandler,
		APIKeyStore:  apiKeyStore,
		Logger:       logger,
	}

	mux := http.NewServeMux()
	registerPublicRoutes(mux, cfg)

	// 验证 /health 端点存在
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterProtectedRoutes(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	apiKeyStore := coreservice.NewMemoryAPIKeyStore()
	coreService := coreservice.NewService(apiKeyStore)
	fileService := fileservice.NewService(1024*1024, "/tmp/test-sandbox")
	shellService := shellservice.NewService(30, "/tmp/test-sandbox")

	coreHandler := core.NewHandler(coreService, logger)
	fileHandler := file.NewHandler(fileService, logger)
	shellHandler := shell.NewHandler(shellService, logger)

	cfg := &Config{
		CoreHandler:  coreHandler,
		FileHandler:  fileHandler,
		ShellHandler: shellHandler,
		APIKeyStore:  apiKeyStore,
		Logger:       logger,
	}

	mux := http.NewServeMux()

	// 注册受保护的路由
	authInterceptor := middleware.NewAuthInterceptor(apiKeyStore, logger)
	registerProtectedRoutes(mux, cfg, authInterceptor)

	// 验证不会 panic
	assert.NotNil(t, mux)
}
