package orchestrator

import (
	"context"
	"testing"
	"time"

	"multi-agent-graph/agent"
	"multi-agent-graph/graph"
)

func TestNewOrchestrator(t *testing.T) {
	dag := graph.NewDAG()
	orc := NewOrchestrator(dag)

	if orc == nil {
		t.Fatal("NewOrchestrator returned nil")
	}

	if orc.DAG != dag {
		t.Error("DAG not set correctly")
	}

	if orc.Agents == nil {
		t.Error("Agents map should not be nil")
	}

	if orc.Registry == nil {
		t.Error("Registry should not be nil")
	}

	if orc.Context == nil {
		t.Error("Context should not be nil")
	}

	if orc.Cancel == nil {
		t.Error("Cancel function should not be nil")
	}

	if orc.Results == nil {
		t.Error("Results map should not be nil")
	}

	if orc.Errors == nil {
		t.Error("Errors map should not be nil")
	}
}

func TestRegisterAgent(t *testing.T) {
	dag := graph.NewDAG()
	node := graph.NewNode("test-node", "Test Node", "data_processor")
	dag.AddNode(node)

	orc := NewOrchestrator(dag)
	testAgent := agent.NewDataProcessorAgent("agent-1", "Test Agent")

	err := orc.RegisterAgent("test-node", testAgent)
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}

	_, exists := orc.Agents["test-node"]
	if !exists {
		t.Error("Agent not registered")
	}
}

func TestRegisterAgentNonExistentNode(t *testing.T) {
	dag := graph.NewDAG()
	orc := NewOrchestrator(dag)
	testAgent := agent.NewDataProcessorAgent("agent-1", "Test Agent")

	err := orc.RegisterAgent("non-existent-node", testAgent)
	if err == nil {
		t.Error("Expected error when registering agent for non-existent node")
	}
}

func TestExecuteSimpleWorkflow(t *testing.T) {
	builder := NewWorkflowBuilder()
	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")
	builder.AddEdge("node1", "node2", 1.0)

	dag := builder.Build()
	orc := NewOrchestrator(dag)

	agent1 := agent.NewDataProcessorAgent("agent-1", "Processor 1")
	agent2 := agent.NewTransformAgent("agent-2", "Transformer")

	orc.RegisterAgent("node1", agent1)
	orc.RegisterAgent("node2", agent2)

	input := map[string]interface{}{"data": "test"}
	err := orc.Execute(input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	result1, err := orc.GetResult("node1")
	if err != nil {
		t.Fatalf("GetResult for node1 failed: %v", err)
	}

	if result1 == nil {
		t.Error("Expected non-nil result for node1")
	}

	result2, err := orc.GetResult("node2")
	if err != nil {
		t.Fatalf("GetResult for node2 failed: %v", err)
	}

	if result2 == nil {
		t.Error("Expected non-nil result for node2")
	}
}

func TestExecuteParallelWorkflow(t *testing.T) {
	builder := NewWorkflowBuilder()
	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "data_processor")
	builder.AddNode("node3", "Node 3", "aggregator")
	builder.AddEdge("node1", "node3", 1.0)
	builder.AddEdge("node2", "node3", 1.0)

	dag := builder.Build()
	orc := NewOrchestrator(dag)

	agent1 := agent.NewDataProcessorAgent("agent-1", "Processor 1")
	agent2 := agent.NewDataProcessorAgent("agent-2", "Processor 2")
	agent3 := agent.NewAggregatorAgent("agent-3", "Aggregator")

	orc.RegisterAgent("node1", agent1)
	orc.RegisterAgent("node2", agent2)
	orc.RegisterAgent("node3", agent3)

	input := map[string]interface{}{"data": "test"}
	err := orc.ExecuteParallel(input)
	if err != nil {
		t.Fatalf("ExecuteParallel failed: %v", err)
	}

	_, err = orc.GetResult("node1")
	if err != nil {
		t.Fatalf("GetResult for node1 failed: %v", err)
	}

	_, err = orc.GetResult("node2")
	if err != nil {
		t.Fatalf("GetResult for node2 failed: %v", err)
	}

	_, err = orc.GetResult("node3")
	if err != nil {
		t.Fatalf("GetResult for node3 failed: %v", err)
	}
}

func TestGetResult(t *testing.T) {
	dag := graph.NewDAG()
	orc := NewOrchestrator(dag)

	// Test getting non-existent result
	_, err := orc.GetResult("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent result")
	}

	// Test getting existing result
	orc.Results["test-node"] = "test-result"
	result, err := orc.GetResult("test-node")
	if err != nil {
		t.Fatalf("GetResult failed: %v", err)
	}

	if result != "test-result" {
		t.Errorf("Expected 'test-result', got '%v'", result)
	}
}

func TestGetError(t *testing.T) {
	dag := graph.NewDAG()
	orc := NewOrchestrator(dag)

	// Test getting non-existent error
	_, exists := orc.GetError("non-existent")
	if exists {
		t.Error("Expected false for non-existent error")
	}

	// Test getting existing error
	testErr := testError("test error")
	orc.Errors["test-node"] = testErr
	retrievedErr, exists := orc.GetError("test-node")
	if !exists {
		t.Error("Expected true for existing error")
	}

	if retrievedErr != testErr {
		t.Error("Retrieved error doesn't match")
	}
}

func TestGetAllResults(t *testing.T) {
	dag := graph.NewDAG()
	orc := NewOrchestrator(dag)

	orc.Results["node1"] = "result1"
	orc.Results["node2"] = "result2"

	results := orc.GetAllResults()

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if results["node1"] != "result1" {
		t.Error("node1 result mismatch")
	}

	if results["node2"] != "result2" {
		t.Error("node2 result mismatch")
	}
}

func TestGetAllErrors(t *testing.T) {
	dag := graph.NewDAG()
	orc := NewOrchestrator(dag)

	testErr1 := testError("error1")
	testErr2 := testError("error2")
	orc.Errors["node1"] = testErr1
	orc.Errors["node2"] = testErr2

	errors := orc.GetAllErrors()

	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}
}

func TestCancelExecution(t *testing.T) {
	dag := graph.NewDAG()
	orc := NewOrchestrator(dag)

	orc.CancelExecution()

	if orc.Context.Err() == nil {
		t.Error("Expected context to be cancelled")
	}
}

func TestGetStatus(t *testing.T) {
	builder := NewWorkflowBuilder()
	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")
	builder.AddEdge("node1", "node2", 1.0)

	dag := builder.Build()
	orc := NewOrchestrator(dag)

	status := orc.GetStatus()

	if len(status) != 2 {
		t.Errorf("Expected 2 status entries, got %d", len(status))
	}

	if _, exists := status["node1"]; !exists {
		t.Error("Expected status for node1")
	}

	if _, exists := status["node2"]; !exists {
		t.Error("Expected status for node2")
	}
}

func TestExecuteWithMissingAgent(t *testing.T) {
	builder := NewWorkflowBuilder()
	builder.AddNode("node1", "Node 1", "data_processor")

	dag := builder.Build()
	orc := NewOrchestrator(dag)

	// Don't register any agent
	input := map[string]interface{}{"data": "test"}
	err := orc.Execute(input)
	if err != nil {
		t.Fatalf("Execute should not fail even with missing agent: %v", err)
	}

	// Check that error was recorded
	_, exists := orc.GetError("node1")
	if !exists {
		t.Error("Expected error to be recorded for node with missing agent")
	}
}

func TestExecuteWithCancelledContext(t *testing.T) {
	builder := NewWorkflowBuilder()
	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")
	builder.AddEdge("node1", "node2", 1.0)

	dag := builder.Build()
	orc := NewOrchestrator(dag)

	agent1 := agent.NewDataProcessorAgent("agent-1", "Processor 1")
	agent2 := agent.NewTransformAgent("agent-2", "Transformer")

	orc.RegisterAgent("node1", agent1)
	orc.RegisterAgent("node2", agent2)

	// Cancel context before execution
	orc.CancelExecution()

	input := map[string]interface{}{"data": "test"}
	err := orc.Execute(input)
	if err == nil {
		t.Error("Expected error when executing with cancelled context")
	}
}

func TestExecuteParallelWithCancelledContext(t *testing.T) {
	builder := NewWorkflowBuilder()
	builder.AddNode("node1", "Node 1", "data_processor")

	dag := builder.Build()
	orc := NewOrchestrator(dag)

	agent1 := agent.NewDataProcessorAgent("agent-1", "Processor 1")
	orc.RegisterAgent("node1", agent1)

	// Cancel context before execution
	orc.CancelExecution()

	input := map[string]interface{}{"data": "test"}
	err := orc.ExecuteParallel(input)
	if err == nil {
		t.Error("Expected error when executing parallel with cancelled context")
	}
}

func TestPrepareNodeInput(t *testing.T) {
	builder := NewWorkflowBuilder()
	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")
	builder.AddEdge("node1", "node2", 1.0)

	dag := builder.Build()
	orc := NewOrchestrator(dag)

	agent1 := agent.NewDataProcessorAgent("agent-1", "Processor 1")
	orc.RegisterAgent("node1", agent1)

	// Execute node1 to get result
	agent1.Execute(orc.Context, map[string]interface{}{"data": "test"})
	orc.Results["node1"] = "node1-result"

	// Test input for node with dependencies
	node2, _ := dag.GetNode("node2")
	input := orc.prepareNodeInput(node2, nil)

	if input != "node1-result" {
		t.Errorf("Expected 'node1-result', got '%v'", input)
	}
}

func TestPrepareNodeInputNoDependencies(t *testing.T) {
	builder := NewWorkflowBuilder()
	builder.AddNode("node1", "Node 1", "data_processor")

	dag := builder.Build()
	orc := NewOrchestrator(dag)

	node1, _ := dag.GetNode("node1")
	input := orc.prepareNodeInput(node1, "initial-input")

	if input != "initial-input" {
		t.Errorf("Expected 'initial-input', got '%v'", input)
	}
}

func TestCustomAgentExecution(t *testing.T) {
	builder := NewWorkflowBuilder()
	builder.AddNode("node1", "Node 1", "custom")

	dag := builder.Build()
	orc := NewOrchestrator(dag)

	// Create a custom agent with Execute method
	customAgent := &CustomTestAgent{
		Agent: agent.NewAgent("agent-1", "Custom Agent", "custom"),
	}

	orc.RegisterAgent("node1", customAgent)

	input := "test-input"
	err := orc.Execute(input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	result, err := orc.GetResult("node1")
	if err != nil {
		t.Fatalf("GetResult failed: %v", err)
	}

	if result != "custom-executed" {
		t.Errorf("Expected 'custom-executed', got '%v'", result)
	}
}

// CustomTestAgent is a custom agent for testing
type CustomTestAgent struct {
	*agent.Agent
}

func (a *CustomTestAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	time.Sleep(10 * time.Millisecond)

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return "custom-executed", nil
}

// testError is a simple error for testing
type testError string

func (e testError) Error() string {
	return string(e)
}
