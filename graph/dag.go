package graph

import (
	"errors"
	"fmt"
	"sync"
)

// DAG represents a Directed Acyclic Graph
type DAG struct {
	Nodes map[string]*Node
	Edges []*Edge
	mu    sync.RWMutex
}

// NewDAG creates a new DAG
func NewDAG() *DAG {
	return &DAG{
		Nodes: make(map[string]*Node),
		Edges: make([]*Edge, 0),
	}
}

// AddNode adds a node to the DAG
func (d *DAG) AddNode(node *Node) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.Nodes[node.ID]; exists {
		return fmt.Errorf("node with ID %s already exists", node.ID)
	}

	d.Nodes[node.ID] = node
	return nil
}

// AddEdge adds an edge to the DAG
func (d *DAG) AddEdge(from, to string, weight float64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.Nodes[from]; !exists {
		return fmt.Errorf("node %s does not exist", from)
	}

	if _, exists := d.Nodes[to]; !exists {
		return fmt.Errorf("node %s does not exist", to)
	}

	// Check if edge already exists
	for _, edge := range d.Edges {
		if edge.From == from && edge.To == to {
			return fmt.Errorf("edge from %s to %s already exists", from, to)
		}
	}

	// Add edge
	d.Edges = append(d.Edges, NewEdge(from, to, weight))

	// Update node inputs/outputs
	d.Nodes[from].Outputs = append(d.Nodes[from].Outputs, to)
	d.Nodes[to].Inputs = append(d.Nodes[to].Inputs, from)

	// Check for cycles
	if d.hasCycle() {
		// Rollback edge addition
		d.Edges = d.Edges[:len(d.Edges)-1]
		d.Nodes[from].Outputs = d.Nodes[from].Outputs[:len(d.Nodes[from].Outputs)-1]
		d.Nodes[to].Inputs = d.Nodes[to].Inputs[:len(d.Nodes[to].Inputs)-1]
		return errors.New("adding this edge would create a cycle")
	}

	return nil
}

// GetNode returns a node by ID
func (d *DAG) GetNode(id string) (*Node, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	node, exists := d.Nodes[id]
	if !exists {
		return nil, fmt.Errorf("node %s does not exist", id)
	}

	return node, nil
}

// GetNodes returns all nodes
func (d *DAG) GetNodes() []*Node {
	d.mu.RLock()
	defer d.mu.RUnlock()

	nodes := make([]*Node, 0, len(d.Nodes))
	for _, node := range d.Nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// GetEdges returns all edges
func (d *DAG) GetEdges() []*Edge {
	d.mu.RLock()
	defer d.mu.RUnlock()

	edges := make([]*Edge, len(d.Edges))
	copy(edges, d.Edges)

	return edges
}

// GetRootNodes returns nodes with no inputs (root nodes)
func (d *DAG) GetRootNodes() []*Node {
	d.mu.RLock()
	defer d.mu.RUnlock()

	roots := make([]*Node, 0)
	for _, node := range d.Nodes {
		if len(node.Inputs) == 0 {
			roots = append(roots, node)
		}
	}

	return roots
}

// GetLeafNodes returns nodes with no outputs (leaf nodes)
func (d *DAG) GetLeafNodes() []*Node {
	d.mu.RLock()
	defer d.mu.RUnlock()

	leaves := make([]*Node, 0)
	for _, node := range d.Nodes {
		if len(node.Outputs) == 0 {
			leaves = append(leaves, node)
		}
	}

	return leaves
}

// GetDependents returns nodes that depend on the given node
func (d *DAG) GetDependents(nodeID string) []*Node {
	d.mu.RLock()
	defer d.mu.RUnlock()

	dependents := make([]*Node, 0)
	for _, edge := range d.Edges {
		if edge.From == nodeID {
			if node, exists := d.Nodes[edge.To]; exists {
				dependents = append(dependents, node)
			}
		}
	}

	return dependents
}

// GetDependencies returns nodes that the given node depends on
func (d *DAG) GetDependencies(nodeID string) []*Node {
	d.mu.RLock()
	defer d.mu.RUnlock()

	dependencies := make([]*Node, 0)
	for _, edge := range d.Edges {
		if edge.To == nodeID {
			if node, exists := d.Nodes[edge.From]; exists {
				dependencies = append(dependencies, node)
			}
		}
	}

	return dependencies
}

// hasCycle checks if the DAG has a cycle using DFS
func (d *DAG) hasCycle() bool {
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	for nodeID := range d.Nodes {
		if !visited[nodeID] {
			if d.dfsCycleCheck(nodeID, visited, recursionStack) {
				return true
			}
		}
	}

	return false
}

// dfsCycleCheck performs DFS to detect cycles
func (d *DAG) dfsCycleCheck(nodeID string, visited, recursionStack map[string]bool) bool {
	visited[nodeID] = true
	recursionStack[nodeID] = true

	for _, edge := range d.Edges {
		if edge.From == nodeID {
			if !visited[edge.To] {
				if d.dfsCycleCheck(edge.To, visited, recursionStack) {
					return true
				}
			} else if recursionStack[edge.To] {
				return true
			}
		}
	}

	recursionStack[nodeID] = false
	return false
}

// TopologicalSort returns nodes in topological order
func (d *DAG) TopologicalSort() ([]*Node, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.hasCycle() {
		return nil, errors.New("cannot perform topological sort on graph with cycles")
	}

	// Kahn's algorithm
	inDegree := make(map[string]int)
	for nodeID := range d.Nodes {
		inDegree[nodeID] = 0
	}

	for _, edge := range d.Edges {
		inDegree[edge.To]++
	}

	queue := make([]*Node, 0)
	for nodeID, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, d.Nodes[nodeID])
		}
	}

	result := make([]*Node, 0)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		for _, edge := range d.Edges {
			if edge.From == node.ID {
				inDegree[edge.To]--
				if inDegree[edge.To] == 0 {
					queue = append(queue, d.Nodes[edge.To])
				}
			}
		}
	}

	if len(result) != len(d.Nodes) {
		return nil, errors.New("topological sort failed - graph has cycles")
	}

	return result, nil
}
