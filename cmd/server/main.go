// Package main provides the entry point for the agent-sandbox API server.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HJH0924/agent-sandbox/domain/core"
	coreService "github.com/HJH0924/agent-sandbox/domain/core/service"
	"github.com/HJH0924/agent-sandbox/domain/file"
	fileService "github.com/HJH0924/agent-sandbox/domain/file/service"
	"github.com/HJH0924/agent-sandbox/domain/shell"
	shellService "github.com/HJH0924/agent-sandbox/domain/shell/service"
	"github.com/HJH0924/agent-sandbox/internal/config"
	"github.com/HJH0924/agent-sandbox/internal/router"

	"github.com/spf13/cobra"
)

const (
	Version = "1.0.0"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:     "api-server",
		Short:   "Agent Sandbox API Server",
		Long:    "A sandbox service for Agent to execute shell commands and file operations",
		Run:     run,
		Version: Version,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "configs/config.toml", "config file path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(_ *cobra.Command, _ []string) {
	// 加载配置
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	logger := initLogger(cfg.Log)
	slog.SetDefault(logger)

	logger.Info("starting agent sandbox server",
		slog.String("version", Version),
		slog.String("config", configFile))

	// 确保工作目录存在
	if err := os.MkdirAll(cfg.Sandbox.WorkspaceDir, 0o750); err != nil {
		logger.Error("failed to create workspace directory",
			slog.String("dir", cfg.Sandbox.WorkspaceDir),
			slog.Any("error", err))
		os.Exit(1)
	}

	// 创建 API Key 存储
	apiKeyStore := coreService.NewMemoryAPIKeyStore()

	// 创建服务
	coreSvc := coreService.NewService(apiKeyStore)
	fileSvc := fileService.NewService(cfg.Sandbox.MaxFileSize, cfg.Sandbox.WorkspaceDir)
	shellSvc := shellService.NewService(cfg.Sandbox.ShellTimeout, cfg.Sandbox.WorkspaceDir)

	// 创建处理器
	coreHandler := core.NewHandler(coreSvc, logger)
	fileHandler := file.NewHandler(fileSvc, logger)
	shellHandler := shell.NewHandler(shellSvc, logger)

	// 设置路由
	handler := router.Setup(&router.Config{
		CoreHandler:  coreHandler,
		FileHandler:  fileHandler,
		ShellHandler: shellHandler,
		APIKeyStore:  apiKeyStore,
		Logger:       logger,
	})

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		logger.Info("server listening",
			slog.String("address", cfg.Server.Address()))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", slog.Any("error", err))
	}

	logger.Info("server stopped")
}

func initLogger(cfg config.LogConfig) *slog.Logger {
	// 解析日志级别
	levelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	level, ok := levelMap[cfg.Level]
	if !ok {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	// 根据格式创建 handler
	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
