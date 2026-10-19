package graph

import (
	"sync"
)

// Node represents a node in the graph
type Node struct {
	ID      string
	Name    string
	Type    string
	Data    interface{}
	Status  string // "pending", "running", "completed", "failed"
	Inputs  []string
	Outputs []string
	mu      sync.RWMutex
}

// NewNode creates a new node
func NewNode(id, name, nodeType string) *Node {
	return &Node{
		ID:     id,
		Name:   name,
		Type:   nodeType,
		Status: "pending",
	}
}

// SetStatus sets the node status
func (n *Node) SetStatus(status string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Status = status
}

// GetStatus returns the node status
func (n *Node) GetStatus() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Status
}

// SetData sets the node data
func (n *Node) SetData(data interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Data = data
}

// GetData returns the node data
func (n *Node) GetData() interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Data
}
