package graph

import (
	"testing"
)

func TestNewDAG(t *testing.T) {
	dag := NewDAG()
	if dag == nil {
		t.Fatal("NewDAG returned nil")
	}
	if dag.Nodes == nil {
		t.Error("Nodes map is nil")
	}
	if dag.Edges == nil {
		t.Error("Edges slice is nil")
	}
}

func TestAddNode(t *testing.T) {
	dag := NewDAG()
	node := NewNode("test", "Test Node", "test_type")

	err := dag.AddNode(node)
	if err != nil {
		t.Fatalf("AddNode failed: %v", err)
	}

	retrieved, err := dag.GetNode("test")
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}

	if retrieved.ID != node.ID {
		t.Errorf("Expected ID %s, got %s", node.ID, retrieved.ID)
	}
}

func TestAddDuplicateNode(t *testing.T) {
	dag := NewDAG()
	node := NewNode("test", "Test Node", "test_type")

	err := dag.AddNode(node)
	if err != nil {
		t.Fatalf("First AddNode failed: %v", err)
	}

	err = dag.AddNode(node)
	if err == nil {
		t.Error("Expected error when adding duplicate node")
	}
}

func TestAddEdge(t *testing.T) {
	dag := NewDAG()
	node1 := NewNode("node1", "Node 1", "type1")
	node2 := NewNode("node2", "Node 2", "type2")

	dag.AddNode(node1)
	dag.AddNode(node2)

	err := dag.AddEdge("node1", "node2", 1.0)
	if err != nil {
		t.Fatalf("AddEdge failed: %v", err)
	}

	edges := dag.GetEdges()
	if len(edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(edges))
	}

	if edges[0].From != "node1" || edges[0].To != "node2" {
		t.Error("Edge has incorrect endpoints")
	}
}

func TestCycleDetection(t *testing.T) {
	dag := NewDAG()
	node1 := NewNode("node1", "Node 1", "type1")
	node2 := NewNode("node2", "Node 2", "type2")
	node3 := NewNode("node3", "Node 3", "type3")

	dag.AddNode(node1)
	dag.AddNode(node2)
	dag.AddNode(node3)

	dag.AddEdge("node1", "node2", 1.0)
	dag.AddEdge("node2", "node3", 1.0)

	// This should create a cycle
	err := dag.AddEdge("node3", "node1", 1.0)
	if err == nil {
		t.Error("Expected error when creating cycle")
	}
}

func TestTopologicalSort(t *testing.T) {
	dag := NewDAG()
	node1 := NewNode("node1", "Node 1", "type1")
	node2 := NewNode("node2", "Node 2", "type2")
	node3 := NewNode("node3", "Node 3", "type3")

	dag.AddNode(node1)
	dag.AddNode(node2)
	dag.AddNode(node3)

	dag.AddEdge("node1", "node2", 1.0)
	dag.AddEdge("node2", "node3", 1.0)

	sorted, err := dag.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if len(sorted) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(sorted))
	}

	// Check that node1 comes before node2, and node2 before node3
	order := make(map[string]int)
	for i, node := range sorted {
		order[node.ID] = i
	}

	if order["node1"] >= order["node2"] {
		t.Error("node1 should come before node2")
	}
	if order["node2"] >= order["node3"] {
		t.Error("node2 should come before node3")
	}
}

func TestRootNodes(t *testing.T) {
	dag := NewDAG()
	node1 := NewNode("node1", "Node 1", "type1")
	node2 := NewNode("node2", "Node 2", "type2")
	node3 := NewNode("node3", "Node 3", "type3")

	dag.AddNode(node1)
	dag.AddNode(node2)
	dag.AddNode(node3)

	dag.AddEdge("node1", "node2", 1.0)
	dag.AddEdge("node1", "node3", 1.0)

	roots := dag.GetRootNodes()
	if len(roots) != 1 {
		t.Errorf("Expected 1 root node, got %d", len(roots))
	}

	if roots[0].ID != "node1" {
		t.Errorf("Expected root node1, got %s", roots[0].ID)
	}
}

func TestLeafNodes(t *testing.T) {
	dag := NewDAG()
	node1 := NewNode("node1", "Node 1", "type1")
	node2 := NewNode("node2", "Node 2", "type2")
	node3 := NewNode("node3", "Node 3", "type3")

	dag.AddNode(node1)
	dag.AddNode(node2)
	dag.AddNode(node3)

	dag.AddEdge("node1", "node2", 1.0)
	dag.AddEdge("node1", "node3", 1.0)

	leaves := dag.GetLeafNodes()
	if len(leaves) != 2 {
		t.Errorf("Expected 2 leaf nodes, got %d", len(leaves))
	}
}
