package orchestrator

import (
	"fmt"

	"multi-agent-graph/graph"
)

// WorkflowBuilder helps build workflows
type WorkflowBuilder struct {
	dag *graph.DAG
}

// NewWorkflowBuilder creates a new workflow builder
func NewWorkflowBuilder() *WorkflowBuilder {
	return &WorkflowBuilder{
		dag: graph.NewDAG(),
	}
}

// AddNode adds a node to the workflow
func (wb *WorkflowBuilder) AddNode(id, name, nodeType string) error {
	node := graph.NewNode(id, name, nodeType)
	return wb.dag.AddNode(node)
}

// AddEdge adds an edge between two nodes
func (wb *WorkflowBuilder) AddEdge(from, to string, weight float64) error {
	return wb.dag.AddEdge(from, to, weight)
}

// Build returns the DAG
func (wb *WorkflowBuilder) Build() *graph.DAG {
	return wb.dag
}

// Validate validates the workflow
func (wb *WorkflowBuilder) Validate() error {
	nodes := wb.dag.GetNodes()
	if len(nodes) == 0 {
		return fmt.Errorf("workflow has no nodes")
	}

	_, err := wb.dag.TopologicalSort()
	if err != nil {
		return fmt.Errorf("workflow validation failed: %w", err)
	}

	return nil
}
