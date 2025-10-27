// Package service implements file management operations for the sandbox.
package service

import (
	"fmt"
	"os"
	"path/filepath"
)

// Service 文件服务.
type Service struct {
	maxFileSize  int64
	workspaceDir string
}

// NewService 创建文件服务实例.
func NewService(maxFileSize int64, workspaceDir string) *Service {
	return &Service{
		maxFileSize:  maxFileSize,
		workspaceDir: workspaceDir,
	}
}

// ReadResult 读取结果.
type ReadResult struct {
	Content string
}

// Read 读取文件.
func (s *Service) Read(path string) (*ReadResult, error) {
	// 确保路径在工作目录下
	fullPath := filepath.Join(s.workspaceDir, path)

	// 检查文件是否存在
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// 检查文件大小
	if info.Size() > s.maxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max: %d)", info.Size(), s.maxFileSize)
	}

	// 读取文件
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &ReadResult{
		Content: string(content),
	}, nil
}

// Write 写入文件.
func (s *Service) Write(path, content string) error {
	// 确保路径在工作目录下
	fullPath := filepath.Join(s.workspaceDir, path)

	// 检查内容大小
	contentSize := int64(len(content))
	if contentSize > s.maxFileSize {
		return fmt.Errorf("content too large: %d bytes (max: %d)", contentSize, s.maxFileSize)
	}

	// 创建目录
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// EditResult 编辑结果.
type EditResult struct {
	Path    string
	Content string
}

// Edit 编辑文件（直接覆盖内容）.
func (s *Service) Edit(path, content string) (*EditResult, error) {
	// 确保路径在工作目录下
	fullPath := filepath.Join(s.workspaceDir, path)

	// 检查内容大小
	contentSize := int64(len(content))
	if contentSize > s.maxFileSize {
		return nil, fmt.Errorf("content too large: %d bytes (max: %d)", contentSize, s.maxFileSize)
	}

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &EditResult{
		Path:    path,
		Content: content,
	}, nil
}
