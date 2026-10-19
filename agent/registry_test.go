package agent

import (
	"testing"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	if registry == nil {
		t.Fatal("NewRegistry returned nil")
	}

	if registry.agents == nil {
		t.Error("agents map should not be nil")
	}
}

func TestRegister(t *testing.T) {
	registry := NewRegistry()

	factory := func(id, name string) *Agent {
		return NewAgent(id, name, "custom_type")
	}

	err := registry.Register("custom_type", factory)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Test duplicate registration
	err = registry.Register("custom_type", factory)
	if err == nil {
		t.Error("Expected error when registering duplicate agent type")
	}
}

func TestCreate(t *testing.T) {
	registry := NewRegistry()

	factory := func(id, name string) *Agent {
		return NewAgent(id, name, "custom_type")
	}

	registry.Register("custom_type", factory)

	agent, err := registry.Create("custom_type", "test-id", "Test Agent")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if agent.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", agent.ID)
	}

	if agent.Name != "Test Agent" {
		t.Errorf("Expected Name 'Test Agent', got '%s'", agent.Name)
	}

	if agent.Type != "custom_type" {
		t.Errorf("Expected Type 'custom_type', got '%s'", agent.Type)
	}
}

func TestCreateNonExistentType(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Create("non_existent_type", "test-id", "Test Agent")
	if err == nil {
		t.Error("Expected error when creating non-existent agent type")
	}
}

func TestGetTypes(t *testing.T) {
	registry := NewRegistry()

	factory1 := func(id, name string) *Agent {
		return NewAgent(id, name, "type1")
	}
	factory2 := func(id, name string) *Agent {
		return NewAgent(id, name, "type2")
	}

	registry.Register("type1", factory1)
	registry.Register("type2", factory2)

	types := registry.GetTypes()

	if len(types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(types))
	}

	typeSet := make(map[string]bool)
	for _, t := range types {
		typeSet[t] = true
	}

	if !typeSet["type1"] {
		t.Error("Expected type1 in types")
	}

	if !typeSet["type2"] {
		t.Error("Expected type2 in types")
	}
}

func TestGetTypesEmpty(t *testing.T) {
	registry := NewRegistry()

	types := registry.GetTypes()

	if len(types) != 0 {
		t.Errorf("Expected 0 types, got %d", len(types))
	}
}

func TestRegisterMultipleTypes(t *testing.T) {
	registry := NewRegistry()

	types := []string{"data_processor", "transform", "validation", "aggregator"}

	for _, agentType := range types {
		factory := func(id, name string) *Agent {
			return NewAgent(id, name, agentType)
		}
		err := registry.Register(agentType, factory)
		if err != nil {
			t.Fatalf("Register failed for type %s: %v", agentType, err)
		}
	}

	retrievedTypes := registry.GetTypes()
	if len(retrievedTypes) != len(types) {
		t.Errorf("Expected %d types, got %d", len(types), len(retrievedTypes))
	}
}

func TestCreateWithDifferentFactories(t *testing.T) {
	registry := NewRegistry()

	factory1 := func(id, name string) *Agent {
		agent := NewAgent(id, name, "type1")
		agent.SetConfig("source", "factory1")
		return agent
	}

	factory2 := func(id, name string) *Agent {
		agent := NewAgent(id, name, "type2")
		agent.SetConfig("source", "factory2")
		return agent
	}

	registry.Register("type1", factory1)
	registry.Register("type2", factory2)

	agent1, _ := registry.Create("type1", "id1", "Agent1")
	agent2, _ := registry.Create("type2", "id2", "Agent2")

	if val, _ := agent1.GetConfig("source"); val != "factory1" {
		t.Error("Expected agent1 config source to be factory1")
	}

	if val, _ := agent2.GetConfig("source"); val != "factory2" {
		t.Error("Expected agent2 config source to be factory2")
	}
}
