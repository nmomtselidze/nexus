package orchestrator

import (
	"testing"
)

func TestNewWorkflowBuilder(t *testing.T) {
	builder := NewWorkflowBuilder()

	if builder == nil {
		t.Fatal("NewWorkflowBuilder returned nil")
	}

	if builder.dag == nil {
		t.Error("DAG should not be nil")
	}
}

func TestAddNode(t *testing.T) {
	builder := NewWorkflowBuilder()

	err := builder.AddNode("node1", "Node 1", "data_processor")
	if err != nil {
		t.Fatalf("AddNode failed: %v", err)
	}

	dag := builder.Build()
	node, err := dag.GetNode("node1")
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}

	if node.ID != "node1" {
		t.Errorf("Expected ID 'node1', got '%s'", node.ID)
	}

	if node.Name != "Node 1" {
		t.Errorf("Expected Name 'Node 1', got '%s'", node.Name)
	}

	if node.Type != "data_processor" {
		t.Errorf("Expected Type 'data_processor', got '%s'", node.Type)
	}
}

func TestAddDuplicateNode(t *testing.T) {
	builder := NewWorkflowBuilder()

	err := builder.AddNode("node1", "Node 1", "data_processor")
	if err != nil {
		t.Fatalf("First AddNode failed: %v", err)
	}

	err = builder.AddNode("node1", "Node 1 Duplicate", "data_processor")
	if err == nil {
		t.Error("Expected error when adding duplicate node")
	}
}

func TestAddEdge(t *testing.T) {
	builder := NewWorkflowBuilder()

	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")

	err := builder.AddEdge("node1", "node2", 1.0)
	if err != nil {
		t.Fatalf("AddEdge failed: %v", err)
	}

	dag := builder.Build()
	edges := dag.GetEdges()

	if len(edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(edges))
	}

	if edges[0].From != "node1" {
		t.Errorf("Expected From 'node1', got '%s'", edges[0].From)
	}

	if edges[0].To != "node2" {
		t.Errorf("Expected To 'node2', got '%s'", edges[0].To)
	}

	if edges[0].Weight != 1.0 {
		t.Errorf("Expected Weight 1.0, got %f", edges[0].Weight)
	}
}

func TestAddEdgeNonExistentNode(t *testing.T) {
	builder := NewWorkflowBuilder()

	builder.AddNode("node1", "Node 1", "data_processor")

	err := builder.AddEdge("node1", "non-existent", 1.0)
	if err == nil {
		t.Error("Expected error when adding edge to non-existent node")
	}
}

func TestBuild(t *testing.T) {
	builder := NewWorkflowBuilder()

	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")
	builder.AddEdge("node1", "node2", 1.0)

	dag := builder.Build()

	if dag == nil {
		t.Fatal("Build returned nil")
	}

	nodes := dag.GetNodes()
	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}
}

func TestValidateEmptyWorkflow(t *testing.T) {
	builder := NewWorkflowBuilder()

	err := builder.Validate()
	if err == nil {
		t.Error("Expected error when validating empty workflow")
	}
}

func TestValidateValidWorkflow(t *testing.T) {
	builder := NewWorkflowBuilder()

	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")
	builder.AddEdge("node1", "node2", 1.0)

	err := builder.Validate()
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
}

func TestValidateWorkflowWithCycle(t *testing.T) {
	builder := NewWorkflowBuilder()

	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")
	builder.AddNode("node3", "Node 3", "aggregator")

	builder.AddEdge("node1", "node2", 1.0)
	builder.AddEdge("node2", "node3", 1.0)

	// Adding the cycle edge should fail at AddEdge level
	err := builder.AddEdge("node3", "node1", 1.0)
	if err == nil {
		t.Error("Expected error when adding edge that creates cycle")
	}
}

func TestValidateWorkflowWithDisconnectedNodes(t *testing.T) {
	builder := NewWorkflowBuilder()

	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")

	err := builder.Validate()
	if err != nil {
		t.Fatalf("Validate failed for disconnected nodes: %v", err)
	}
}

func TestBuildComplexWorkflow(t *testing.T) {
	builder := NewWorkflowBuilder()

	// Build a diamond pattern
	builder.AddNode("input", "Input", "data_processor")
	builder.AddNode("branch1", "Branch 1", "data_processor")
	builder.AddNode("branch2", "Branch 2", "data_processor")
	builder.AddNode("merge", "Merge", "aggregator")
	builder.AddNode("output", "Output", "transform")

	builder.AddEdge("input", "branch1", 1.0)
	builder.AddEdge("input", "branch2", 1.0)
	builder.AddEdge("branch1", "merge", 1.0)
	builder.AddEdge("branch2", "merge", 1.0)
	builder.AddEdge("merge", "output", 1.0)

	err := builder.Validate()
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	dag := builder.Build()

	nodes := dag.GetNodes()
	if len(nodes) != 5 {
		t.Errorf("Expected 5 nodes, got %d", len(nodes))
	}

	edges := dag.GetEdges()
	if len(edges) != 5 {
		t.Errorf("Expected 5 edges, got %d", len(edges))
	}
}

func TestAddMultipleEdges(t *testing.T) {
	builder := NewWorkflowBuilder()

	builder.AddNode("node1", "Node 1", "data_processor")
	builder.AddNode("node2", "Node 2", "transform")
	builder.AddNode("node3", "Node 3", "aggregator")

	builder.AddEdge("node1", "node2", 1.0)
	builder.AddEdge("node1", "node3", 1.0)
	builder.AddEdge("node2", "node3", 1.0)

	dag := builder.Build()
	edges := dag.GetEdges()

	if len(edges) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(edges))
	}
}

func TestWorkflowBuilderReusability(t *testing.T) {
	builder := NewWorkflowBuilder()

	// Build first workflow
	builder.AddNode("node1", "Node 1", "data_processor")
	dag1 := builder.Build()

	// Build second workflow with same builder
	builder2 := NewWorkflowBuilder()
	builder2.AddNode("node2", "Node 2", "transform")
	dag2 := builder2.Build()

	// Verify they are independent
	_, err := dag1.GetNode("node1")
	if err != nil {
		t.Error("node1 should exist in dag1")
	}

	_, err = dag2.GetNode("node1")
	if err == nil {
		t.Error("node1 should not exist in dag2")
	}

	_, err = dag2.GetNode("node2")
	if err != nil {
		t.Error("node2 should exist in dag2")
	}
}
