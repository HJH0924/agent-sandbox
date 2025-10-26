package service

import (
	"testing"
)

func TestMemoryAPIKeyStore(t *testing.T) {
	store := NewMemoryAPIKeyStore()

	// Test Store and Verify
	sandboxID := "test-sandbox-123"
	apiKey := "test-api-key-456"

	err := store.Store(sandboxID, apiKey)
	if err != nil {
		t.Fatalf("Failed to store API key: %v", err)
	}

	// Verify the key
	retrievedID, ok := store.Verify(apiKey)
	if !ok {
		t.Fatal("Failed to verify API key")
	}

	if retrievedID != sandboxID {
		t.Fatalf("Expected sandbox ID %s, got %s", sandboxID, retrievedID)
	}

	// Test Delete
	err = store.Delete(sandboxID)
	if err != nil {
		t.Fatalf("Failed to delete API key: %v", err)
	}

	// Verify key is deleted
	_, ok = store.Verify(apiKey)
	if ok {
		t.Fatal("API key should have been deleted")
	}
}

func TestInitSandbox(t *testing.T) {
	store := NewMemoryAPIKeyStore()
	service := NewService(store)

	// Initialize sandbox
	result, err := service.InitSandbox()
	if err != nil {
		t.Fatalf("Failed to initialize sandbox: %v", err)
	}

	// Check sandbox ID is not empty
	if result.SandboxID == "" {
		t.Fatal("Sandbox ID should not be empty")
	}

	// Check API key is not empty and has correct prefix
	if result.APIKey == "" {
		t.Fatal("API key should not be empty")
	}

	if result.APIKey[:3] != "sk_" {
		t.Fatalf("API key should start with 'sk_', got: %s", result.APIKey[:3])
	}

	// Verify the API key is stored
	retrievedID, ok := store.Verify(result.APIKey)
	if !ok {
		t.Fatal("API key should be stored in the store")
	}

	if retrievedID != result.SandboxID {
		t.Fatalf("Expected sandbox ID %s, got %s", result.SandboxID, retrievedID)
	}
}
