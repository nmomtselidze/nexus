package examples

import (
	"fmt"

	"multi-agent-graph/agent"
	"multi-agent-graph/graph"
	"multi-agent-graph/orchestrator"
)

// LinearWorkflow creates a simple linear workflow
func LinearWorkflow() (*orchestrator.Orchestrator, error) {
	// Build DAG
	builder := orchestrator.NewWorkflowBuilder()

	// Add nodes
	builder.AddNode("validate", "Validation", "validation")
	builder.AddNode("process", "Data Processing", "data_processor")
	builder.AddNode("transform", "Transformation", "transform")
	builder.AddNode("aggregate", "Aggregation", "aggregator")

	// Add edges (linear flow)
	builder.AddEdge("validate", "process", 1.0)
	builder.AddEdge("process", "transform", 1.0)
	builder.AddEdge("transform", "aggregate", 1.0)

	// Validate
	if err := builder.Validate(); err != nil {
		return nil, err
	}

	dag := builder.Build()

	// Create orchestrator
	orc := orchestrator.NewOrchestrator(dag)

	// Register agents
	orc.RegisterAgent("validate", agent.NewValidationAgent("agent-validate", "Validator"))
	orc.RegisterAgent("process", agent.NewDataProcessorAgent("agent-process", "Processor"))
	orc.RegisterAgent("transform", agent.NewTransformAgent("agent-transform", "Transformer"))
	orc.RegisterAgent("aggregate", agent.NewAggregatorAgent("agent-aggregate", "Aggregator"))

	return orc, nil
}

// ParallelWorkflow creates a workflow with parallel branches
func ParallelWorkflow() (*orchestrator.Orchestrator, error) {
	// Build DAG
	builder := orchestrator.NewWorkflowBuilder()

	// Add nodes
	builder.AddNode("input", "Input", "data_processor")
	builder.AddNode("branch1", "Branch 1 Processing", "data_processor")
	builder.AddNode("branch2", "Branch 2 Processing", "data_processor")
	builder.AddNode("branch3", "Branch 3 Processing", "data_processor")
	builder.AddNode("merge", "Merge Results", "aggregator")
	builder.AddNode("output", "Output", "transform")

	// Add edges (parallel flow)
	builder.AddEdge("input", "branch1", 1.0)
	builder.AddEdge("input", "branch2", 1.0)
	builder.AddEdge("input", "branch3", 1.0)
	builder.AddEdge("branch1", "merge", 1.0)
	builder.AddEdge("branch2", "merge", 1.0)
	builder.AddEdge("branch3", "merge", 1.0)
	builder.AddEdge("merge", "output", 1.0)

	// Validate
	if err := builder.Validate(); err != nil {
		return nil, err
	}

	dag := builder.Build()

	// Create orchestrator
	orc := orchestrator.NewOrchestrator(dag)

	// Register agents
	orc.RegisterAgent("input", agent.NewDataProcessorAgent("agent-input", "Input Processor"))
	orc.RegisterAgent("branch1", agent.NewDataProcessorAgent("agent-branch1", "Branch 1"))
	orc.RegisterAgent("branch2", agent.NewDataProcessorAgent("agent-branch2", "Branch 2"))
	orc.RegisterAgent("branch3", agent.NewDataProcessorAgent("agent-branch3", "Branch 3"))
	orc.RegisterAgent("merge", agent.NewAggregatorAgent("agent-merge", "Merger"))
	orc.RegisterAgent("output", agent.NewTransformAgent("agent-output", "Output Transformer"))

	return orc, nil
}

// ComplexWorkflow creates a complex workflow with multiple dependencies
func ComplexWorkflow() (*orchestrator.Orchestrator, error) {
	// Build DAG
	builder := orchestrator.NewWorkflowBuilder()

	// Add nodes
	builder.AddNode("data_source1", "Data Source 1", "data_processor")
	builder.AddNode("data_source2", "Data Source 2", "data_processor")
	builder.AddNode("data_source3", "Data Source 3", "data_processor")
	builder.AddNode("validate1", "Validate 1", "validation")
	builder.AddNode("validate2", "Validate 2", "validation")
	builder.AddNode("process1", "Process 1", "data_processor")
	builder.AddNode("process2", "Process 2", "data_processor")
	builder.AddNode("transform1", "Transform 1", "transform")
	builder.AddNode("transform2", "Transform 2", "transform")
	builder.AddNode("merge", "Merge", "aggregator")
	builder.AddNode("final", "Final Output", "transform")

	// Add edges (complex dependencies)
	builder.AddEdge("data_source1", "validate1", 1.0)
	builder.AddEdge("data_source2", "validate1", 1.0)
	builder.AddEdge("data_source3", "validate2", 1.0)
	builder.AddEdge("validate1", "process1", 1.0)
	builder.AddEdge("validate2", "process2", 1.0)
	builder.AddEdge("process1", "transform1", 1.0)
	builder.AddEdge("process2", "transform2", 1.0)
	builder.AddEdge("transform1", "merge", 1.0)
	builder.AddEdge("transform2", "merge", 1.0)
	builder.AddEdge("merge", "final", 1.0)

	// Validate
	if err := builder.Validate(); err != nil {
		return nil, err
	}

	dag := builder.Build()

	// Create orchestrator
	orc := orchestrator.NewOrchestrator(dag)

	// Register agents
	orc.RegisterAgent("data_source1", agent.NewDataProcessorAgent("agent-ds1", "Data Source 1"))
	orc.RegisterAgent("data_source2", agent.NewDataProcessorAgent("agent-ds2", "Data Source 2"))
	orc.RegisterAgent("data_source3", agent.NewDataProcessorAgent("agent-ds3", "Data Source 3"))
	orc.RegisterAgent("validate1", agent.NewValidationAgent("agent-val1", "Validator 1"))
	orc.RegisterAgent("validate2", agent.NewValidationAgent("agent-val2", "Validator 2"))
	orc.RegisterAgent("process1", agent.NewDataProcessorAgent("agent-proc1", "Processor 1"))
	orc.RegisterAgent("process2", agent.NewDataProcessorAgent("agent-proc2", "Processor 2"))
	orc.RegisterAgent("transform1", agent.NewTransformAgent("agent-trans1", "Transformer 1"))
	orc.RegisterAgent("transform2", agent.NewTransformAgent("agent-trans2", "Transformer 2"))
	orc.RegisterAgent("merge", agent.NewAggregatorAgent("agent-merge", "Merger"))
	orc.RegisterAgent("final", agent.NewTransformAgent("agent-final", "Final"))

	return orc, nil
}

// PrintDAG prints the DAG structure
func PrintDAG(dag *graph.DAG) {
	fmt.Println("\n=== DAG Structure ===")
	fmt.Printf("Nodes: %d\n", len(dag.GetNodes()))
	fmt.Printf("Edges: %d\n", len(dag.GetEdges()))

	fmt.Println("\nNodes:")
	for _, node := range dag.GetNodes() {
		fmt.Printf("  - %s (%s) [%s]\n", node.ID, node.Name, node.Type)
	}

	fmt.Println("\nEdges:")
	for _, edge := range dag.GetEdges() {
		fmt.Printf("  - %s -> %s (weight: %.1f)\n", edge.From, edge.To, edge.Weight)
	}

	fmt.Println("\nRoot Nodes:")
	for _, node := range dag.GetRootNodes() {
		fmt.Printf("  - %s\n", node.ID)
	}

	fmt.Println("\nLeaf Nodes:")
	for _, node := range dag.GetLeafNodes() {
		fmt.Printf("  - %s\n", node.ID)
	}

	fmt.Println("\nTopological Order:")
	ordered, err := dag.TopologicalSort()
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		for i, node := range ordered {
			fmt.Printf("  %d. %s\n", i+1, node.ID)
		}
	}
	fmt.Println("====================")
}
