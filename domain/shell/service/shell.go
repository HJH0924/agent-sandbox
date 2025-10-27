package service

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

// Service Shell 服务
type Service struct {
	defaultTimeout time.Duration
	workspaceDir   string
}

// NewService 创建 Shell 服务实例
func NewService(defaultTimeout int, workspaceDir string) *Service {
	return &Service{
		defaultTimeout: time.Duration(defaultTimeout) * time.Second,
		workspaceDir:   workspaceDir,
	}
}

// ExecuteResult 执行结果
type ExecuteResult struct {
	Output string
}

// Execute 执行 Shell 命令
func (s *Service) Execute(ctx context.Context, command string) (*ExecuteResult, error) {
	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(ctx, s.defaultTimeout)
	defer cancel()

	// 创建命令
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// 设置工作目录
	if s.workspaceDir != "" {
		// 确保工作目录存在
		absPath, err := filepath.Abs(s.workspaceDir)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path: %w", err)
		}
		cmd.Dir = absPath
	}

	// 捕获输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()

	// 合并 stdout 和 stderr
	output := stdout.String()
	if stderr.Len() > 0 {
		if len(output) > 0 {
			output += "\n"
		}
		output += stderr.String()
	}

	if err != nil {
		// 如果命令执行失败，返回带输出的错误
		if output != "" {
			return &ExecuteResult{Output: output}, fmt.Errorf("command execution failed: %w", err)
		}
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	return &ExecuteResult{
		Output: output,
	}, nil
}
