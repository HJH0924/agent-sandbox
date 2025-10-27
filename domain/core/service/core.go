// Package service implements the core business logic for sandbox management.
package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// APIKeyStore API 密钥存储接口.
type APIKeyStore interface {
	// Store 存储 API 密钥.
	Store(sandboxID, apiKey string) error
	// Verify 验证 API 密钥并返回沙箱 ID.
	Verify(apiKey string) (sandboxID string, ok bool)
	// Delete 删除 API 密钥.
	Delete(sandboxID string) error
}

// MemoryAPIKeyStore API 密钥的内存存储实现.
type MemoryAPIKeyStore struct {
	mu         sync.RWMutex
	keys       map[string]string // apiKey -> sandboxID
	sandboxes  map[string]string // sandboxID -> apiKey
	timestamps map[string]time.Time
}

// NewMemoryAPIKeyStore 创建基于内存的 API 密钥存储.
func NewMemoryAPIKeyStore() *MemoryAPIKeyStore {
	return &MemoryAPIKeyStore{
		keys:       make(map[string]string),
		sandboxes:  make(map[string]string),
		timestamps: make(map[string]time.Time),
	}
}

// Store 存储 API 密钥.
func (s *MemoryAPIKeyStore) Store(sandboxID, apiKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.keys[apiKey] = sandboxID
	s.sandboxes[sandboxID] = apiKey
	s.timestamps[sandboxID] = time.Now()

	return nil
}

// Verify 验证 API 密钥.
func (s *MemoryAPIKeyStore) Verify(apiKey string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sandboxID, ok := s.keys[apiKey]
	return sandboxID, ok
}

// Delete 删除 API 密钥.
func (s *MemoryAPIKeyStore) Delete(sandboxID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if apiKey, ok := s.sandboxes[sandboxID]; ok {
		delete(s.keys, apiKey)
		delete(s.sandboxes, sandboxID)
		delete(s.timestamps, sandboxID)
	}

	return nil
}

// Service 核心服务.
type Service struct {
	store APIKeyStore
}

// NewService 创建核心服务实例.
func NewService(store APIKeyStore) *Service {
	return &Service{
		store: store,
	}
}

// InitSandboxResult 沙箱初始化结果.
type InitSandboxResult struct {
	SandboxID string
	APIKey    string
	CreatedAt time.Time
}

// InitSandbox 初始化沙箱，生成沙箱 ID 和 API 密钥.
func (s *Service) InitSandbox() (*InitSandboxResult, error) {
	// 生成沙箱 ID
	sandboxID := uuid.New().String()

	// 生成 API 密钥（32字节随机数，hex编码）
	apiKeyBytes := make([]byte, 32)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return nil, fmt.Errorf("failed to generate api key: %w", err)
	}
	apiKey := "sk_" + hex.EncodeToString(apiKeyBytes)

	// 存储 API 密钥
	if err := s.store.Store(sandboxID, apiKey); err != nil {
		return nil, fmt.Errorf("failed to store api key: %w", err)
	}

	return &InitSandboxResult{
		SandboxID: sandboxID,
		APIKey:    apiKey,
		CreatedAt: time.Now(),
	}, nil
}
