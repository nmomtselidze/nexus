package agent

import (
	"context"
	"fmt"
	"time"
)

// Agent represents an autonomous agent that can execute tasks
type Agent struct {
	ID          string
	Name        string
	Type        string
	Config      map[string]interface{}
	Status      string
	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time
}

// NewAgent creates a new agent
func NewAgent(id, name, agentType string) *Agent {
	return &Agent{
		ID:        id,
		Name:      name,
		Type:      agentType,
		Config:    make(map[string]interface{}),
		Status:    "idle",
		CreatedAt: time.Now(),
	}
}

// SetConfig sets a configuration value
func (a *Agent) SetConfig(key string, value interface{}) {
	a.Config[key] = value
}

// GetConfig gets a configuration value
func (a *Agent) GetConfig(key string) (interface{}, bool) {
	value, exists := a.Config[key]
	return value, exists
}

// Execute executes the agent's task - to be overridden by specific agent types
func (a *Agent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Execute not implemented for agent type: %s", a.Type)
}

// GetStatus returns the agent status
func (a *Agent) GetStatus() string {
	return a.Status
}

// GetDuration returns the execution duration
func (a *Agent) GetDuration() time.Duration {
	if a.StartedAt.IsZero() {
		return 0
	}
	if a.CompletedAt.IsZero() {
		return time.Since(a.StartedAt)
	}
	return a.CompletedAt.Sub(a.StartedAt)
}
