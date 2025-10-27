package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileService(t *testing.T) {
	// Create temporary workspace directory
	tmpDir, err := os.MkdirTemp("", "agent-sandbox-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	maxFileSize := int64(1024 * 1024) // 1MB
	service := NewService(maxFileSize, tmpDir)

	// Test Write
	testPath := "test/example.txt"
	testContent := "Hello, World!"

	err = service.Write(testPath, testContent)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Verify file exists
	fullPath := filepath.Join(tmpDir, testPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Fatal("File should exist after write")
	}

	// Test Read
	result, err := service.Read(testPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if result.Content != testContent {
		t.Fatalf("Expected content %q, got %q", testContent, result.Content)
	}

	// Test Edit
	newContent := "Hello, Updated World!"

	editResult, err := service.Edit(testPath, newContent)
	if err != nil {
		t.Fatalf("Failed to edit file: %v", err)
	}

	if editResult.Content != newContent {
		t.Fatalf("Expected content %q, got %q", newContent, editResult.Content)
	}

	// Verify edited content
	result, err = service.Read(testPath)
	if err != nil {
		t.Fatalf("Failed to read file after edit: %v", err)
	}

	if result.Content != newContent {
		t.Fatalf("Expected content %q after edit, got %q", newContent, result.Content)
	}
}

func TestFileService_ReadNonExistent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-sandbox-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	service := NewService(1024*1024, tmpDir)

	// Try to read non-existent file
	_, err = service.Read("nonexistent.txt")
	if err == nil {
		t.Fatal("Expected error when reading non-existent file")
	}
}

func TestFileService_TooLarge(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-sandbox-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	maxSize := int64(100) // Very small max size
	service := NewService(maxSize, tmpDir)

	// Try to write content larger than max size
	largeContent := string(make([]byte, 200))

	err = service.Write("large.txt", largeContent)
	if err == nil {
		t.Fatal("Expected error when writing content larger than max size")
	}
}
