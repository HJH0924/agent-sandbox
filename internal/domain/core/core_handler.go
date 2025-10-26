package core

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	corev1 "github.com/HJH0924/agent-sandbox/sdk/go/core/v1"
	"github.com/HJH0924/agent-sandbox/internal/domain/core/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Handler 核心服务处理器
type Handler struct {
	coreService *service.Service
	logger      *slog.Logger
}

// NewHandler 创建核心服务处理器
func NewHandler(coreService *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		coreService: coreService,
		logger:      logger,
	}
}

// InitSandbox 初始化新沙箱并返回沙箱 ID 和 API 密钥
func (h *Handler) InitSandbox(
	ctx context.Context,
	req *connect.Request[corev1.InitSandboxRequest],
) (*connect.Response[corev1.InitSandboxResponse], error) {
	h.logger.InfoContext(ctx, "initializing sandbox")

	// 调用 service 层初始化沙箱
	result, err := h.coreService.InitSandbox()
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to initialize sandbox",
			slog.Any("error", err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	h.logger.InfoContext(ctx, "sandbox initialized successfully",
		slog.String("sandbox_id", result.SandboxID))

	// 返回响应
	return connect.NewResponse(&corev1.InitSandboxResponse{
		SandboxId: result.SandboxID,
		ApiKey:    result.APIKey,
		CreatedAt: timestamppb.New(result.CreatedAt),
	}), nil
}
