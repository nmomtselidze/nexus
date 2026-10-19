package examples

import (
	"context"
	"fmt"
	"time"

	"multi-agent-graph/agent"
	"multi-agent-graph/orchestrator"
)

// DatabaseData represents data from a database
type DatabaseData struct {
	Source      string
	Records     []map[string]interface{}
	Count       int
	ExtractedAt time.Time
}

// ValidatedData represents validated data
type ValidatedData struct {
	Source      string
	Data        []map[string]interface{}
	Valid       bool
	Issues      []string
	ValidatedAt time.Time
}

// MergedData represents merged data from multiple sources
type MergedData struct {
	Sources    []string
	Records    []map[string]interface{}
	TotalCount int
	MergedAt   time.Time
}

// WarehouseLoadResult represents warehouse load result
type WarehouseLoadResult struct {
	Success  bool
	RowCount int
	LoadTime time.Duration
	LoadedAt time.Time
}

// DatabaseExtractor extracts data from a database
type DatabaseExtractor struct {
	*agent.Agent
	ConnectionString string
	Query            string
}

// NewDatabaseExtractor creates a new database extractor agent
func NewDatabaseExtractor(id, name, connectionString, query string) *DatabaseExtractor {
	return &DatabaseExtractor{
		Agent:            agent.NewAgent(id, name, "database_extractor"),
		ConnectionString: connectionString,
		Query:            query,
	}
}

// Execute extracts data from the database
func (a *DatabaseExtractor) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Simulate database extraction
	select {
	case <-time.After(150 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	// Simulate extracted data
	records := make([]map[string]interface{}, 0)
	for i := 0; i < 10; i++ {
		records = append(records, map[string]interface{}{
			"id":        i + 1,
			"name":      fmt.Sprintf("Record %d from %s", i+1, a.Agent.Name),
			"value":     (i + 1) * 100,
			"source":    a.Agent.Name,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	result := &DatabaseData{
		Source:      a.Agent.Name,
		Records:     records,
		Count:       len(records),
		ExtractedAt: time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// DataValidator validates data
type DataValidator struct {
	*agent.Agent
	Rules []ValidationRule
}

type ValidationRule struct {
	Name  string
	Check func(map[string]interface{}) bool
}

// NewDataValidator creates a new data validator agent
func NewDataValidator(id, name string) *DataValidator {
	return &DataValidator{
		Agent: agent.NewAgent(id, name, "data_validator"),
		Rules: []ValidationRule{
			{
				Name: "has_id",
				Check: func(record map[string]interface{}) bool {
					_, ok := record["id"]
					return ok
				},
			},
			{
				Name: "has_name",
				Check: func(record map[string]interface{}) bool {
					_, ok := record["name"]
					return ok
				},
			},
		},
	}
}

// Execute validates the data
func (a *DataValidator) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	dbData, ok := input.(*DatabaseData)
	if !ok {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *DatabaseData")
	}

	select {
	case <-time.After(50 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	// Validate each record
	validRecords := make([]map[string]interface{}, 0)
	issues := make([]string, 0)

	for _, record := range dbData.Records {
		allValid := true
		for _, rule := range a.Rules {
			if !rule.Check(record) {
				allValid = false
				issues = append(issues, fmt.Sprintf("Record %d failed rule: %s", record["id"], rule.Name))
			}
		}
		if allValid {
			validRecords = append(validRecords, record)
		}
	}

	result := &ValidatedData{
		Source:      dbData.Source,
		Data:        validRecords,
		Valid:       len(issues) == 0,
		Issues:      issues,
		ValidatedAt: time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// DataMerger merges data from multiple sources
type DataMerger struct {
	*agent.Agent
}

// NewDataMerger creates a new data merger agent
func NewDataMerger(id, name string) *DataMerger {
	return &DataMerger{
		Agent: agent.NewAgent(id, name, "data_merger"),
	}
}

// Execute merges data from multiple sources
func (a *DataMerger) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Handle both []*ValidatedData and []interface{}
	var validatedDataList []*ValidatedData
	switch v := input.(type) {
	case []*ValidatedData:
		validatedDataList = v
	case []interface{}:
		validatedDataList = make([]*ValidatedData, 0, len(v))
		for _, item := range v {
			if vd, ok := item.(*ValidatedData); ok {
				validatedDataList = append(validatedDataList, vd)
			}
		}
	default:
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected []*ValidatedData or []interface{}")
	}

	select {
	case <-time.After(80 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	// Merge data from all sources
	sources := make([]string, 0)
	allRecords := make([]map[string]interface{}, 0)
	totalCount := 0

	for _, vd := range validatedDataList {
		sources = append(sources, vd.Source)
		allRecords = append(allRecords, vd.Data...)
		totalCount += len(vd.Data)
	}

	result := &MergedData{
		Sources:    sources,
		Records:    allRecords,
		TotalCount: totalCount,
		MergedAt:   time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// WarehouseLoader loads data to warehouse
type WarehouseLoader struct {
	*agent.Agent
	WarehouseURL string
	BatchSize    int
}

// NewWarehouseLoader creates a new warehouse loader agent
func NewWarehouseLoader(id, name, warehouseURL string, batchSize int) *WarehouseLoader {
	return &WarehouseLoader{
		Agent:        agent.NewAgent(id, name, "warehouse_loader"),
		WarehouseURL: warehouseURL,
		BatchSize:    batchSize,
	}
}

// Execute loads data to the warehouse
func (a *WarehouseLoader) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	mergedData, ok := input.(*MergedData)
	if !ok {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *MergedData")
	}

	// Simulate batch loading
	batches := (len(mergedData.Records) + a.BatchSize - 1) / a.BatchSize
	for i := 0; i < batches; i++ {
		select {
		case <-time.After(50 * time.Millisecond):
		case <-ctx.Done():
			a.Agent.Status = "failed"
			a.Agent.CompletedAt = time.Now()
			return nil, ctx.Err()
		}
	}

	result := &WarehouseLoadResult{
		Success:  true,
		RowCount: len(mergedData.Records),
		LoadTime: time.Since(a.Agent.StartedAt),
		LoadedAt: time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// ETLPipeline creates an ETL pipeline workflow
func ETLPipeline() (*orchestrator.Orchestrator, error) {
	// Build DAG
	builder := orchestrator.NewWorkflowBuilder()

	// Add nodes - 3 database sources
	builder.AddNode("db1_extract", "Database 1 Extract", "database_extractor")
	builder.AddNode("db2_extract", "Database 2 Extract", "database_extractor")
	builder.AddNode("db3_extract", "Database 3 Extract", "database_extractor")

	// Add validation nodes for each source
	builder.AddNode("db1_validate", "Database 1 Validate", "data_validator")
	builder.AddNode("db2_validate", "Database 2 Validate", "data_validator")
	builder.AddNode("db3_validate", "Database 3 Validate", "data_validator")

	// Add merge node
	builder.AddNode("merge", "Merge Data", "data_merger")

	// Add warehouse load node
	builder.AddNode("warehouse_load", "Load to Warehouse", "warehouse_loader")

	// Add edges - extract -> validate (parallel for each DB)
	builder.AddEdge("db1_extract", "db1_validate", 1.0)
	builder.AddEdge("db2_extract", "db2_validate", 1.0)
	builder.AddEdge("db3_extract", "db3_validate", 1.0)

	// Add edges - validate -> merge (all validators feed into merge)
	builder.AddEdge("db1_validate", "merge", 1.0)
	builder.AddEdge("db2_validate", "merge", 1.0)
	builder.AddEdge("db3_validate", "merge", 1.0)

	// Add edge - merge -> warehouse load
	builder.AddEdge("merge", "warehouse_load", 1.0)

	// Validate
	if err := builder.Validate(); err != nil {
		return nil, err
	}

	dag := builder.Build()

	// Create orchestrator
	orc := orchestrator.NewOrchestrator(dag)

	// Register agents
	orc.RegisterAgent("db1_extract", NewDatabaseExtractor("extract-1", "DB1", "postgres://db1", "SELECT * FROM users"))
	orc.RegisterAgent("db2_extract", NewDatabaseExtractor("extract-2", "DB2", "mysql://db2", "SELECT * FROM orders"))
	orc.RegisterAgent("db3_extract", NewDatabaseExtractor("extract-3", "DB3", "mongodb://db3", "db.products.find()"))
	orc.RegisterAgent("db1_validate", NewDataValidator("validate-1", "Validator 1"))
	orc.RegisterAgent("db2_validate", NewDataValidator("validate-2", "Validator 2"))
	orc.RegisterAgent("db3_validate", NewDataValidator("validate-3", "Validator 3"))
	orc.RegisterAgent("merge", NewDataMerger("merge-1", "Data Merger"))
	orc.RegisterAgent("warehouse_load", NewWarehouseLoader("load-1", "Warehouse Loader", "snowflake://warehouse", 100))

	return orc, nil
}

// PrintETLResults prints ETL pipeline results
func PrintETLResults(orc *orchestrator.Orchestrator) {
	fmt.Println("\n=== ETL Pipeline Results ===")

	// Database extraction results
	fmt.Println("Database Extraction:")
	for _, nodeID := range []string{"db1_extract", "db2_extract", "db3_extract"} {
		if result, err := orc.GetResult(nodeID); err == nil {
			if dbData, ok := result.(*DatabaseData); ok {
				fmt.Printf("  %s: %d records from %s\n", dbData.Source, dbData.Count, dbData.Source)
			}
		}
	}

	// Validation results
	fmt.Println("\nValidation:")
	for _, nodeID := range []string{"db1_validate", "db2_validate", "db3_validate"} {
		if result, err := orc.GetResult(nodeID); err == nil {
			if vData, ok := result.(*ValidatedData); ok {
				fmt.Printf("  %s: %d valid records, %d issues\n", vData.Source, len(vData.Data), len(vData.Issues))
			}
		}
	}

	// Merge results
	if result, err := orc.GetResult("merge"); err == nil {
		if mData, ok := result.(*MergedData); ok {
			fmt.Printf("\nMerge: %d total records from %d sources\n", mData.TotalCount, len(mData.Sources))
			fmt.Printf("  Sources: %v\n", mData.Sources)
		}
	}

	// Warehouse load results
	if result, err := orc.GetResult("warehouse_load"); err == nil {
		if wData, ok := result.(*WarehouseLoadResult); ok {
			fmt.Printf("\nWarehouse Load: %s\n", map[bool]string{true: "SUCCESS", false: "FAILED"}[wData.Success])
			fmt.Printf("  Rows loaded: %d\n", wData.RowCount)
			fmt.Printf("  Load time: %v\n", wData.LoadTime)
		}
	}

	fmt.Println("\n============================")
}
