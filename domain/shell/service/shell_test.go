package service

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestShellService_Execute(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-sandbox-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	service := NewService(30, tmpDir) // 30 seconds timeout

	// Test simple command
	ctx := context.Background()

	result, err := service.Execute(ctx, "echo 'Hello, World!'")
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	expectedOutput := "Hello, World!"
	if !strings.Contains(result.Output, expectedOutput) {
		t.Fatalf("Expected output to contain %q, got %q", expectedOutput, result.Output)
	}
}

func TestShellService_WorkingDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-sandbox-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	service := NewService(30, tmpDir)

	// Execute pwd command to check working directory
	ctx := context.Background()

	result, err := service.Execute(ctx, "pwd")
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// The output should contain the tmpDir path
	if !strings.Contains(result.Output, tmpDir) {
		t.Fatalf("Expected output to contain %q, got %q", tmpDir, result.Output)
	}
}

func TestShellService_FailedCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-sandbox-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	service := NewService(30, tmpDir)

	// Execute a command that will fail
	ctx := context.Background()

	_, err = service.Execute(ctx, "exit 1")
	if err == nil {
		t.Fatal("Expected error for failed command")
	}
}

func TestShellService_Timeout(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-sandbox-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	service := NewService(1, tmpDir) // 1 second timeout

	// Execute a command that takes longer than timeout
	ctx := context.Background()

	_, err = service.Execute(ctx, "sleep 5")
	if err == nil {
		t.Fatal("Expected timeout error")
	}
}
