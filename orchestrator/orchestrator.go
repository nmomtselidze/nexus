package orchestrator

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"multi-agent-graph/agent"
	"multi-agent-graph/graph"
)

// Orchestrator manages workflow execution
type Orchestrator struct {
	DAG      *graph.DAG
	Agents   map[string]interface{} // Can hold various agent types
	Registry *agent.Registry
	Context  context.Context
	Cancel   context.CancelFunc
	Results  map[string]interface{}
	Errors   map[string]error
	mu       sync.RWMutex
}

// NewOrchestrator creates a new orchestrator
func NewOrchestrator(dag *graph.DAG) *Orchestrator {
	ctx, cancel := context.WithCancel(context.Background())
	return &Orchestrator{
		DAG:      dag,
		Agents:   make(map[string]interface{}),
		Registry: agent.NewRegistry(),
		Context:  ctx,
		Cancel:   cancel,
		Results:  make(map[string]interface{}),
		Errors:   make(map[string]error),
	}
}

// RegisterAgent registers an agent with a node
func (o *Orchestrator) RegisterAgent(nodeID string, agentInstance interface{}) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	node, err := o.DAG.GetNode(nodeID)
	if err != nil {
		return fmt.Errorf("node %s not found: %w", nodeID, err)
	}

	o.Agents[nodeID] = agentInstance
	node.SetData(agentInstance)

	return nil
}

// Execute executes the workflow
func (o *Orchestrator) Execute(input interface{}) error {
	// Get topological order
	orderedNodes, err := o.DAG.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to get topological order: %w", err)
	}

	log.Printf("Starting workflow execution with %d nodes", len(orderedNodes))

	// Execute nodes in topological order
	for _, node := range orderedNodes {
		log.Printf("Executing node: %s (%s)", node.ID, node.Name)

		// Check if context is cancelled
		if err := o.Context.Err(); err != nil {
			return fmt.Errorf("workflow cancelled: %w", err)
		}

		// Get agent for this node
		agentInstance, exists := o.Agents[node.ID]
		if !exists {
			o.mu.Lock()
			o.Errors[node.ID] = fmt.Errorf("no agent registered for node %s", node.ID)
			o.mu.Unlock()
			continue
		}

		// Prepare input for this node
		nodeInput := o.prepareNodeInput(node, input)

		// Execute the agent
		output, err := o.executeAgent(node.ID, agentInstance, nodeInput)
		if err != nil {
			o.mu.Lock()
			o.Errors[node.ID] = err
			node.SetStatus("failed")
			o.mu.Unlock()
			log.Printf("Node %s failed: %v", node.ID, err)
			continue
		}

		// Store result
		o.mu.Lock()
		o.Results[node.ID] = output
		node.SetStatus("completed")
		node.SetData(output)
		o.mu.Unlock()

		log.Printf("Node %s completed successfully", node.ID)
	}

	log.Printf("Workflow execution completed")
	return nil
}

// ExecuteParallel executes the workflow with parallel execution where possible
func (o *Orchestrator) ExecuteParallel(input interface{}) error {
	// Get topological order
	orderedNodes, err := o.DAG.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to get topological order: %w", err)
	}

	log.Printf("Starting parallel workflow execution with %d nodes", len(orderedNodes))

	// Track completion status
	completed := make(map[string]bool)
	var mu sync.Mutex

	// Execute nodes in topological order, parallelizing independent nodes
	for _, node := range orderedNodes {
		// Check if context is cancelled
		if err := o.Context.Err(); err != nil {
			return fmt.Errorf("workflow cancelled: %w", err)
		}

		// Wait for all dependencies to complete
		for {
			allDepsDone := true
			deps := o.DAG.GetDependencies(node.ID)
			mu.Lock()
			for _, dep := range deps {
				if !completed[dep.ID] {
					allDepsDone = false
				}
			}
			mu.Unlock()

			if allDepsDone {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		// Execute this node in a goroutine
		go func(n *graph.Node) {
			log.Printf("Executing node: %s (%s)", n.ID, n.Name)

			// Get agent for this node
			agentInstance, exists := o.Agents[n.ID]
			if !exists {
				o.mu.Lock()
				o.Errors[n.ID] = fmt.Errorf("no agent registered for node %s", n.ID)
				o.mu.Unlock()
				return
			}

			// Prepare input for this node
			nodeInput := o.prepareNodeInput(n, input)

			// Execute the agent
			output, err := o.executeAgent(n.ID, agentInstance, nodeInput)
			if err != nil {
				o.mu.Lock()
				o.Errors[n.ID] = err
				n.SetStatus("failed")
				mu.Lock()
				completed[n.ID] = true
				mu.Unlock()
				o.mu.Unlock()
				log.Printf("Node %s failed: %v", n.ID, err)
				return
			}

			// Store result
			o.mu.Lock()
			o.Results[n.ID] = output
			n.SetStatus("completed")
			n.SetData(output)
			mu.Lock()
			completed[n.ID] = true
			mu.Unlock()
			o.mu.Unlock()

			log.Printf("Node %s completed successfully", n.ID)
		}(node)
	}

	// Wait for all nodes to complete
	for {
		allDone := true
		mu.Lock()
		for _, node := range orderedNodes {
			if !completed[node.ID] {
				allDone = false
				break
			}
		}
		mu.Unlock()

		if allDone {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	log.Printf("Parallel workflow execution completed")
	return nil
}

// prepareNodeInput prepares input for a node based on its dependencies
func (o *Orchestrator) prepareNodeInput(node *graph.Node, initialInput interface{}) interface{} {
	deps := o.DAG.GetDependencies(node.ID)

	if len(deps) == 0 {
		return initialInput
	}

	// If node has dependencies, aggregate their results
	inputs := make([]interface{}, 0, len(deps))
	for _, dep := range deps {
		o.mu.RLock()
		if result, exists := o.Results[dep.ID]; exists {
			inputs = append(inputs, result)
		}
		o.mu.RUnlock()
	}

	if len(inputs) == 1 {
		return inputs[0]
	}

	return inputs
}

// executeAgent executes an agent (type assertion)
func (o *Orchestrator) executeAgent(nodeID string, agentInstance interface{}, input interface{}) (interface{}, error) {
	// Handle built-in agent types
	switch a := agentInstance.(type) {
	case *agent.DataProcessorAgent:
		return a.Execute(o.Context, input)
	case *agent.TransformAgent:
		return a.Execute(o.Context, input)
	case *agent.ValidationAgent:
		return a.Execute(o.Context, input)
	case *agent.AggregatorAgent:
		return a.Execute(o.Context, input)
	}

	// Handle custom agent types with Execute method
	if executor, ok := agentInstance.(interface {
		Execute(context.Context, interface{}) (interface{}, error)
	}); ok {
		return executor.Execute(o.Context, input)
	}

	return nil, fmt.Errorf("unknown agent type for node %s", nodeID)
}

// GetResult returns the result for a specific node
func (o *Orchestrator) GetResult(nodeID string) (interface{}, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	result, exists := o.Results[nodeID]
	if !exists {
		return nil, fmt.Errorf("no result for node %s", nodeID)
	}

	return result, nil
}

// GetError returns the error for a specific node
func (o *Orchestrator) GetError(nodeID string) (error, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	err, exists := o.Errors[nodeID]
	return err, exists
}

// GetAllResults returns all results
func (o *Orchestrator) GetAllResults() map[string]interface{} {
	o.mu.RLock()
	defer o.mu.RUnlock()

	results := make(map[string]interface{})
	for k, v := range o.Results {
		results[k] = v
	}

	return results
}

// GetAllErrors returns all errors
func (o *Orchestrator) GetAllErrors() map[string]error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	errors := make(map[string]error)
	for k, v := range o.Errors {
		errors[k] = v
	}

	return errors
}

// CancelExecution cancels the workflow execution
func (o *Orchestrator) CancelExecution() {
	o.Cancel()
}

// GetStatus returns the status of all nodes
func (o *Orchestrator) GetStatus() map[string]string {
	o.mu.RLock()
	defer o.mu.RUnlock()

	status := make(map[string]string)
	for _, node := range o.DAG.GetNodes() {
		status[node.ID] = node.GetStatus()
	}

	return status
}
