package file

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	filev1 "github.com/HJH0924/agent-sandbox/sdk/go/file/v1"
	"github.com/HJH0924/agent-sandbox/internal/domain/file/service"
)

// Handler 文件服务处理器
type Handler struct {
	fileService *service.Service
	logger      *slog.Logger
}

// NewHandler 创建文件服务处理器
func NewHandler(fileService *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		fileService: fileService,
		logger:      logger,
	}
}

// Read 读取文件
func (h *Handler) Read(
	ctx context.Context,
	req *connect.Request[filev1.ReadRequest],
) (*connect.Response[filev1.ReadResponse], error) {
	path := req.Msg.GetPath()

	h.logger.InfoContext(ctx, "reading file",
		slog.String("path", path))

	// 调用 service 层读取文件
	result, err := h.fileService.Read(path)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to read file",
			slog.String("path", path),
			slog.Any("error", err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	h.logger.InfoContext(ctx, "file read successfully",
		slog.String("path", path),
		slog.Int("content_length", len(result.Content)))

	// 返回响应
	return connect.NewResponse(&filev1.ReadResponse{
		Content: result.Content,
	}), nil
}

// Write 写入文件
func (h *Handler) Write(
	ctx context.Context,
	req *connect.Request[filev1.WriteRequest],
) (*connect.Response[filev1.WriteResponse], error) {
	path := req.Msg.GetPath()
	content := req.Msg.GetContent()

	h.logger.InfoContext(ctx, "writing file",
		slog.String("path", path),
		slog.Int("content_length", len(content)))

	// 调用 service 层写入文件
	if err := h.fileService.Write(path, content); err != nil {
		h.logger.ErrorContext(ctx, "failed to write file",
			slog.String("path", path),
			slog.Any("error", err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	h.logger.InfoContext(ctx, "file written successfully",
		slog.String("path", path))

	// 返回响应
	return connect.NewResponse(&filev1.WriteResponse{}), nil
}

// Edit 编辑文件
func (h *Handler) Edit(
	ctx context.Context,
	req *connect.Request[filev1.EditRequest],
) (*connect.Response[filev1.EditResponse], error) {
	path := req.Msg.GetPath()
	content := req.Msg.GetContent()

	h.logger.InfoContext(ctx, "editing file",
		slog.String("path", path),
		slog.Int("content_length", len(content)))

	// 调用 service 层编辑文件
	result, err := h.fileService.Edit(path, content)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to edit file",
			slog.String("path", path),
			slog.Any("error", err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	h.logger.InfoContext(ctx, "file edited successfully",
		slog.String("path", path))

	// 返回响应
	return connect.NewResponse(&filev1.EditResponse{
		Path:    result.Path,
		Content: result.Content,
	}), nil
}
