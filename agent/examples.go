package agent

import (
	"context"
	"fmt"
	"time"
)

// DataProcessorAgent processes data
type DataProcessorAgent struct {
	*Agent
	ProcessorFunc func(data interface{}) (interface{}, error)
}

// NewDataProcessorAgent creates a new data processor agent
func NewDataProcessorAgent(id, name string) *DataProcessorAgent {
	return &DataProcessorAgent{
		Agent: NewAgent(id, name, "data_processor"),
	}
}

// Execute implements the agent execution
func (a *DataProcessorAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	var output interface{}
	var err error

	if a.ProcessorFunc != nil {
		output, err = a.ProcessorFunc(input)
	} else {
		output, err = a.defaultProcess(ctx, input)
	}

	if err != nil {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, err
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return output, nil
}

// defaultProcess provides default data processing logic
func (a *DataProcessorAgent) defaultProcess(ctx context.Context, input interface{}) (interface{}, error) {
	data, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected map[string]interface{}")
	}

	data["processed_at"] = time.Now().Format(time.RFC3339)
	data["processed_by"] = a.Agent.Name

	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return data, nil
}

// TransformAgent transforms data using a custom function
type TransformAgent struct {
	*Agent
	TransformFunc func(interface{}) (interface{}, error)
}

// NewTransformAgent creates a new transform agent
func NewTransformAgent(id, name string) *TransformAgent {
	return &TransformAgent{
		Agent: NewAgent(id, name, "transform"),
	}
}

// Execute implements the agent execution
func (a *TransformAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	var output interface{}
	var err error

	if a.TransformFunc != nil {
		output, err = a.TransformFunc(input)
	} else {
		output, err = a.defaultTransform(ctx, input)
	}

	if err != nil {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, err
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return output, nil
}

// defaultTransform provides default transformation logic
func (a *TransformAgent) defaultTransform(ctx context.Context, input interface{}) (interface{}, error) {
	result := map[string]interface{}{
		"original":    input,
		"transformed": fmt.Sprintf("Transformed by %s at %s", a.Agent.Name, time.Now().Format(time.RFC3339)),
	}

	select {
	case <-time.After(50 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return result, nil
}

// ValidationAgent validates data
type ValidationAgent struct {
	*Agent
	ValidationFunc func(interface{}) (bool, error)
}

// NewValidationAgent creates a new validation agent
func NewValidationAgent(id, name string) *ValidationAgent {
	return &ValidationAgent{
		Agent: NewAgent(id, name, "validation"),
	}
}

// Execute implements the agent execution
func (a *ValidationAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	var output interface{}
	var err error

	if a.ValidationFunc != nil {
		valid, err := a.ValidationFunc(input)
		if err != nil {
			a.Agent.Status = "failed"
			a.Agent.CompletedAt = time.Now()
			return nil, err
		}
		output = map[string]interface{}{
			"valid":        valid,
			"validated_by": a.Agent.Name,
			"validated_at": time.Now().Format(time.RFC3339),
		}
	} else {
		output, err = a.defaultValidate(ctx, input)
	}

	if err != nil {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, err
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return output, nil
}

// defaultValidate provides default validation logic
func (a *ValidationAgent) defaultValidate(ctx context.Context, input interface{}) (interface{}, error) {
	result := map[string]interface{}{
		"valid":        true,
		"validated_by": a.Agent.Name,
		"validated_at": time.Now().Format(time.RFC3339),
	}

	select {
	case <-time.After(30 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return result, nil
}

// AggregatorAgent aggregates results from multiple sources
type AggregatorAgent struct {
	*Agent
	AggregateFunc func([]interface{}) (interface{}, error)
}

// NewAggregatorAgent creates a new aggregator agent
func NewAggregatorAgent(id, name string) *AggregatorAgent {
	return &AggregatorAgent{
		Agent: NewAgent(id, name, "aggregator"),
	}
}

// Execute implements the agent execution
func (a *AggregatorAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	var output interface{}
	var err error

	if a.AggregateFunc != nil {
		inputs, ok := input.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid input type for aggregator, expected []interface{}")
		}
		output, err = a.AggregateFunc(inputs)
	} else {
		output, err = a.defaultAggregate(ctx, input)
	}

	if err != nil {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, err
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return output, nil
}

// defaultAggregate provides default aggregation logic
func (a *AggregatorAgent) defaultAggregate(ctx context.Context, input interface{}) (interface{}, error) {
	result := map[string]interface{}{
		"aggregated":    true,
		"count":         1,
		"aggregated_by": a.Agent.Name,
		"aggregated_at": time.Now().Format(time.RFC3339),
	}

	select {
	case <-time.After(80 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return result, nil
}
