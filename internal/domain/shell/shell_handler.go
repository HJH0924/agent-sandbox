package shell

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	shellv1 "github.com/HJH0924/agent-sandbox/sdk/go/shell/v1"
	"github.com/HJH0924/agent-sandbox/internal/domain/shell/serivce"
)

// Handler Shell 服务处理器
type Handler struct {
	shellService *service.Service
	logger       *slog.Logger
}

// NewHandler 创建 Shell 服务处理器
func NewHandler(shellService *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		shellService: shellService,
		logger:       logger,
	}
}

// Execute 执行 Shell 命令
func (h *Handler) Execute(
	ctx context.Context,
	req *connect.Request[shellv1.ExecuteRequest],
) (*connect.Response[shellv1.ExecuteResponse], error) {
	command := req.Msg.GetCommand()

	h.logger.InfoContext(ctx, "executing shell command",
		slog.String("command", command))

	// 调用 service 层执行命令
	result, err := h.shellService.Execute(ctx, command)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to execute shell command",
			slog.String("command", command),
			slog.Any("error", err))

		// 即使命令执行失败，也返回输出（如果有的话）
		if result != nil && result.Output != "" {
			return connect.NewResponse(&shellv1.ExecuteResponse{
				Output: result.Output,
			}), connect.NewError(connect.CodeInternal, err)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	h.logger.InfoContext(ctx, "shell command executed successfully",
		slog.String("command", command),
		slog.Int("output_length", len(result.Output)))

	// 返回响应
	return connect.NewResponse(&shellv1.ExecuteResponse{
		Output: result.Output,
	}), nil
}
