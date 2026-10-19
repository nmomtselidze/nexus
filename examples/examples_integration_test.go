package examples

import (
	"testing"

	"multi-agent-graph/orchestrator"
)

func TestLinearWorkflow(t *testing.T) {
	orc, err := LinearWorkflow()
	if err != nil {
		t.Fatalf("LinearWorkflow failed: %v", err)
	}

	input := map[string]interface{}{"data": "test"}
	err = orc.Execute(input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	results := orc.GetAllResults()
	if len(results) == 0 {
		t.Error("Expected some results")
	}

	status := orc.GetStatus()
	for nodeID, nodeStatus := range status {
		if nodeStatus != "completed" {
			t.Errorf("Node %s has status %s, expected completed", nodeID, nodeStatus)
		}
	}
}

func TestLinearWorkflowParallel(t *testing.T) {
	orc, err := LinearWorkflow()
	if err != nil {
		t.Fatalf("LinearWorkflow failed: %v", err)
	}

	input := map[string]interface{}{"data": "test"}
	err = orc.ExecuteParallel(input)
	if err != nil {
		t.Fatalf("ExecuteParallel failed: %v", err)
	}

	results := orc.GetAllResults()
	if len(results) == 0 {
		t.Error("Expected some results")
	}

	status := orc.GetStatus()
	for nodeID, nodeStatus := range status {
		if nodeStatus != "completed" {
			t.Errorf("Node %s has status %s, expected completed", nodeID, nodeStatus)
		}
	}
}

func TestParallelWorkflow(t *testing.T) {
	orc, err := ParallelWorkflow()
	if err != nil {
		t.Fatalf("ParallelWorkflow failed: %v", err)
	}

	input := map[string]interface{}{"data": "test"}
	err = orc.ExecuteParallel(input)
	if err != nil {
		t.Fatalf("ExecuteParallel failed: %v", err)
	}

	results := orc.GetAllResults()
	if len(results) == 0 {
		t.Error("Expected some results")
	}

	status := orc.GetStatus()
	for nodeID, nodeStatus := range status {
		if nodeStatus != "completed" {
			t.Errorf("Node %s has status %s, expected completed", nodeID, nodeStatus)
		}
	}
}

func TestComplexWorkflow(t *testing.T) {
	orc, err := ComplexWorkflow()
	if err != nil {
		t.Fatalf("ComplexWorkflow failed: %v", err)
	}

	input := map[string]interface{}{"data": "test"}
	err = orc.ExecuteParallel(input)
	if err != nil {
		t.Fatalf("ExecuteParallel failed: %v", err)
	}

	results := orc.GetAllResults()
	if len(results) == 0 {
		t.Error("Expected some results")
	}

	status := orc.GetStatus()
	for nodeID, nodeStatus := range status {
		if nodeStatus != "completed" {
			t.Errorf("Node %s has status %s, expected completed", nodeID, nodeStatus)
		}
	}
}

func TestETLPipeline(t *testing.T) {
	orc, err := ETLPipeline()
	if err != nil {
		t.Fatalf("ETLPipeline failed: %v", err)
	}

	input := map[string]interface{}{"data": "test"}
	err = orc.ExecuteParallel(input)
	if err != nil {
		t.Fatalf("ExecuteParallel failed: %v", err)
	}

	results := orc.GetAllResults()
	if len(results) == 0 {
		t.Error("Expected some results")
	}

	// Check specific ETL results
	_, err = orc.GetResult("db1_extract")
	if err != nil {
		t.Error("Expected db1_extract result")
	}

	_, err = orc.GetResult("warehouse_load")
	if err != nil {
		t.Error("Expected warehouse_load result")
	}

	status := orc.GetStatus()
	for nodeID, nodeStatus := range status {
		if nodeStatus != "completed" {
			t.Errorf("Node %s has status %s, expected completed", nodeID, nodeStatus)
		}
	}
}

func TestCICDPipeline(t *testing.T) {
	orc, err := CICDPipeline()
	if err != nil {
		t.Fatalf("CICDPipeline failed: %v", err)
	}

	input := map[string]interface{}{"data": "test"}
	err = orc.ExecuteParallel(input)
	if err != nil {
		t.Fatalf("ExecuteParallel failed: %v", err)
	}

	results := orc.GetAllResults()
	if len(results) == 0 {
		t.Error("Expected some results")
	}

	// Check specific CI/CD results
	_, err = orc.GetResult("build")
	if err != nil {
		t.Error("Expected build result")
	}

	_, err = orc.GetResult("deploy_staging")
	if err != nil {
		t.Error("Expected deploy_staging result")
	}

	status := orc.GetStatus()
	for nodeID, nodeStatus := range status {
		if nodeStatus != "completed" {
			t.Errorf("Node %s has status %s, expected completed", nodeID, nodeStatus)
		}
	}
}

func TestMLPipeline(t *testing.T) {
	orc, err := MLPipeline()
	if err != nil {
		t.Fatalf("MLPipeline failed: %v", err)
	}

	input := map[string]interface{}{"data": "test"}
	err = orc.ExecuteParallel(input)
	if err != nil {
		t.Fatalf("ExecuteParallel failed: %v", err)
	}

	results := orc.GetAllResults()
	if len(results) == 0 {
		t.Error("Expected some results")
	}

	// Check specific ML results
	_, err = orc.GetResult("data_collection")
	if err != nil {
		t.Error("Expected data_collection result")
	}

	_, err = orc.GetResult("deploy_model")
	if err != nil {
		t.Error("Expected deploy_model result")
	}

	status := orc.GetStatus()
	for nodeID, nodeStatus := range status {
		if nodeStatus != "completed" {
			t.Errorf("Node %s has status %s, expected completed", nodeID, nodeStatus)
		}
	}
}

func TestFlightBookingWorkflow(t *testing.T) {
	orc, err := FlightBookingWorkflow()
	if err != nil {
		t.Fatalf("FlightBookingWorkflow failed: %v", err)
	}

	input := map[string]interface{}{
		"origin_city":    "New York",
		"dest_city":      "London",
		"departure_date": "2024-01-15",
	}
	err = orc.ExecuteParallel(input)
	if err != nil {
		t.Fatalf("ExecuteParallel failed: %v", err)
	}

	results := orc.GetAllResults()
	if len(results) == 0 {
		t.Error("Expected some results")
	}

	status := orc.GetStatus()
	for nodeID, nodeStatus := range status {
		if nodeStatus != "completed" {
			t.Errorf("Node %s has status %s, expected completed", nodeID, nodeStatus)
		}
	}
}

func TestWorkflowBuilders(t *testing.T) {
	// Test that all workflow builders create valid orchestrators
	workflows := []func() (*orchestrator.Orchestrator, error){
		LinearWorkflow,
		ParallelWorkflow,
		ComplexWorkflow,
		ETLPipeline,
		CICDPipeline,
		MLPipeline,
		FlightBookingWorkflow,
	}

	for i, workflowBuilder := range workflows {
		orc, err := workflowBuilder()
		if err != nil {
			t.Errorf("Workflow %d builder failed: %v", i, err)
			continue
		}

		if orc == nil {
			t.Errorf("Workflow %d returned nil orchestrator", i)
		}

		if orc.DAG == nil {
			t.Errorf("Workflow %d has nil DAG", i)
		}
	}
}
