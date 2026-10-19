package agent

import (
	"fmt"
	"sync"
)

// Registry manages available agent types
type Registry struct {
	agents map[string]func(id, name string) *Agent
	mu     sync.RWMutex
}

// NewRegistry creates a new agent registry
func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]func(id, name string) *Agent),
	}
}

// Register registers an agent type
func (r *Registry) Register(agentType string, factory func(id, name string) *Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[agentType]; exists {
		return fmt.Errorf("agent type %s already registered", agentType)
	}

	r.agents[agentType] = factory
	return nil
}

// Create creates an agent of the specified type
func (r *Registry) Create(agentType, id, name string) (*Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.agents[agentType]
	if !exists {
		return nil, fmt.Errorf("agent type %s not registered", agentType)
	}

	return factory(id, name), nil
}

// GetTypes returns all registered agent types
func (r *Registry) GetTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.agents))
	for agentType := range r.agents {
		types = append(types, agentType)
	}

	return types
}
