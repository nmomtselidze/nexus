package agent

import (
	"context"
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {
	agent := NewAgent("test-id", "Test Agent", "test_type")

	if agent.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", agent.ID)
	}

	if agent.Name != "Test Agent" {
		t.Errorf("Expected Name 'Test Agent', got '%s'", agent.Name)
	}

	if agent.Type != "test_type" {
		t.Errorf("Expected Type 'test_type', got '%s'", agent.Type)
	}

	if agent.Status != "idle" {
		t.Errorf("Expected Status 'idle', got '%s'", agent.Status)
	}

	if agent.Config == nil {
		t.Error("Config should not be nil")
	}

	if agent.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestSetConfig(t *testing.T) {
	agent := NewAgent("test-id", "Test Agent", "test_type")

	agent.SetConfig("key1", "value1")
	agent.SetConfig("key2", 123)

	if val, exists := agent.GetConfig("key1"); !exists || val != "value1" {
		t.Error("Failed to set/get config key1")
	}

	if val, exists := agent.GetConfig("key2"); !exists || val != 123 {
		t.Error("Failed to set/get config key2")
	}
}

func TestGetConfig(t *testing.T) {
	agent := NewAgent("test-id", "Test Agent", "test_type")

	agent.SetConfig("existing_key", "value")

	// Test existing key
	if val, exists := agent.GetConfig("existing_key"); !exists {
		t.Error("Expected existing_key to exist")
	} else if val != "value" {
		t.Errorf("Expected 'value', got '%v'", val)
	}

	// Test non-existing key
	if _, exists := agent.GetConfig("non_existing_key"); exists {
		t.Error("Expected non_existing_key to not exist")
	}
}

func TestBaseAgentExecute(t *testing.T) {
	agent := NewAgent("test-id", "Test Agent", "test_type")

	ctx := context.Background()
	_, err := agent.Execute(ctx, "input")

	if err == nil {
		t.Error("Expected error for base agent Execute")
	}
}

func TestGetStatus(t *testing.T) {
	agent := NewAgent("test-id", "Test Agent", "test_type")

	if agent.GetStatus() != "idle" {
		t.Errorf("Expected status 'idle', got '%s'", agent.GetStatus())
	}

	agent.Status = "running"
	if agent.GetStatus() != "running" {
		t.Errorf("Expected status 'running', got '%s'", agent.GetStatus())
	}
}

func TestGetDuration(t *testing.T) {
	agent := NewAgent("test-id", "Test Agent", "test_type")

	// Test duration when not started
	if duration := agent.GetDuration(); duration != 0 {
		t.Errorf("Expected duration 0 when not started, got %v", duration)
	}

	// Test duration when started but not completed
	agent.StartedAt = time.Now()
	time.Sleep(10 * time.Millisecond)
	duration := agent.GetDuration()
	if duration < 10*time.Millisecond {
		t.Errorf("Expected duration >= 10ms, got %v", duration)
	}

	// Test duration when completed
	agent.CompletedAt = time.Now()
	completedDuration := agent.GetDuration()
	if completedDuration == 0 {
		t.Error("Expected non-zero duration when completed")
	}
}

func TestNewDataProcessorAgent(t *testing.T) {
	agent := NewDataProcessorAgent("test-id", "Test Processor")

	if agent.Agent.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", agent.Agent.ID)
	}

	if agent.Agent.Type != "data_processor" {
		t.Errorf("Expected Type 'data_processor', got '%s'", agent.Agent.Type)
	}
}

func TestDataProcessorAgentExecute(t *testing.T) {
	agent := NewDataProcessorAgent("test-id", "Test Processor")

	ctx := context.Background()
	input := map[string]interface{}{"data": "test"}

	output, err := agent.Execute(ctx, input)

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if agent.Agent.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", agent.Agent.Status)
	}

	result, ok := output.(map[string]interface{})
	if !ok {
		t.Fatal("Expected output to be map[string]interface{}")
	}

	if _, exists := result["processed_at"]; !exists {
		t.Error("Expected 'processed_at' in result")
	}

	if _, exists := result["processed_by"]; !exists {
		t.Error("Expected 'processed_by' in result")
	}
}

func TestDataProcessorAgentExecuteWithCustomFunc(t *testing.T) {
	agent := NewDataProcessorAgent("test-id", "Test Processor")

	customFunc := func(data interface{}) (interface{}, error) {
		return "custom_processed", nil
	}
	agent.ProcessorFunc = customFunc

	ctx := context.Background()
	output, err := agent.Execute(ctx, "input")

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if output != "custom_processed" {
		t.Errorf("Expected 'custom_processed', got '%v'", output)
	}
}

func TestNewTransformAgent(t *testing.T) {
	agent := NewTransformAgent("test-id", "Test Transformer")

	if agent.Agent.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", agent.Agent.ID)
	}

	if agent.Agent.Type != "transform" {
		t.Errorf("Expected Type 'transform', got '%s'", agent.Agent.Type)
	}
}

func TestTransformAgentExecute(t *testing.T) {
	agent := NewTransformAgent("test-id", "Test Transformer")

	ctx := context.Background()
	output, err := agent.Execute(ctx, "input")

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if agent.Agent.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", agent.Agent.Status)
	}

	result, ok := output.(map[string]interface{})
	if !ok {
		t.Fatal("Expected output to be map[string]interface{}")
	}

	if _, exists := result["original"]; !exists {
		t.Error("Expected 'original' in result")
	}

	if _, exists := result["transformed"]; !exists {
		t.Error("Expected 'transformed' in result")
	}
}

func TestNewValidationAgent(t *testing.T) {
	agent := NewValidationAgent("test-id", "Test Validator")

	if agent.Agent.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", agent.Agent.ID)
	}

	if agent.Agent.Type != "validation" {
		t.Errorf("Expected Type 'validation', got '%s'", agent.Agent.Type)
	}
}

func TestValidationAgentExecute(t *testing.T) {
	agent := NewValidationAgent("test-id", "Test Validator")

	ctx := context.Background()
	output, err := agent.Execute(ctx, "input")

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if agent.Agent.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", agent.Agent.Status)
	}

	result, ok := output.(map[string]interface{})
	if !ok {
		t.Fatal("Expected output to be map[string]interface{}")
	}

	if _, exists := result["valid"]; !exists {
		t.Error("Expected 'valid' in result")
	}

	if _, exists := result["validated_by"]; !exists {
		t.Error("Expected 'validated_by' in result")
	}
}

func TestValidationAgentExecuteWithCustomFunc(t *testing.T) {
	agent := NewValidationAgent("test-id", "Test Validator")

	customFunc := func(data interface{}) (bool, error) {
		return true, nil
	}
	agent.ValidationFunc = customFunc

	ctx := context.Background()
	output, err := agent.Execute(ctx, "input")

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	result, ok := output.(map[string]interface{})
	if !ok {
		t.Fatal("Expected output to be map[string]interface{}")
	}

	if result["valid"] != true {
		t.Error("Expected valid to be true")
	}
}

func TestNewAggregatorAgent(t *testing.T) {
	agent := NewAggregatorAgent("test-id", "Test Aggregator")

	if agent.Agent.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", agent.Agent.ID)
	}

	if agent.Agent.Type != "aggregator" {
		t.Errorf("Expected Type 'aggregator', got '%s'", agent.Agent.Type)
	}
}

func TestAggregatorAgentExecute(t *testing.T) {
	agent := NewAggregatorAgent("test-id", "Test Aggregator")

	ctx := context.Background()
	output, err := agent.Execute(ctx, []interface{}{"item1", "item2"})

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if agent.Agent.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", agent.Agent.Status)
	}

	result, ok := output.(map[string]interface{})
	if !ok {
		t.Fatal("Expected output to be map[string]interface{}")
	}

	if _, exists := result["aggregated"]; !exists {
		t.Error("Expected 'aggregated' in result")
	}

	if _, exists := result["aggregated_by"]; !exists {
		t.Error("Expected 'aggregated_by' in result")
	}
}

func TestAgentContextCancellation(t *testing.T) {
	agent := NewDataProcessorAgent("test-id", "Test Processor")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := agent.Execute(ctx, map[string]interface{}{"data": "test"})

	if err == nil {
		t.Error("Expected error due to context cancellation")
	}

	if agent.Agent.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", agent.Agent.Status)
	}
}
