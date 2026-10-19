package examples

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"multi-agent-graph/agent"
	"multi-agent-graph/orchestrator"
)

// DataCollectionResult represents collected data
type DataCollectionResult struct {
	Source      string
	RecordCount int
	Features    int
	FilePath    string
	CollectedAt time.Time
}

// FeatureExtractionResult represents extracted features
type FeatureExtractionResult struct {
	FeatureCount   int
	FeatureNames   []string
	ProcessedCount int
	Preprocessing  []string
	ExtractedAt    time.Time
}

// ModelTrainingResult represents trained model
type ModelTrainingResult struct {
	ModelID      string
	ModelType    string
	Hyperparams  map[string]interface{}
	Accuracy     float64
	Precision    float64
	Recall       float64
	F1Score      float64
	TrainingTime time.Duration
	TrainedAt    time.Time
}

// ModelEvaluationResult represents model evaluation
type ModelEvaluationResult struct {
	ModelID         string
	TestAccuracy    float64
	TestLoss        float64
	ConfusionMatrix map[string][]int
	EvaluatedAt     time.Time
}

// ModelComparisonResult represents comparison of multiple models
type ModelComparisonResult struct {
	BestModelID  string
	BestAccuracy float64
	AllModels    []*ModelTrainingResult
	ComparisonAt time.Time
}

// ModelDeploymentResult represents model deployment
type ModelDeploymentResult struct {
	Success     bool
	ModelID     string
	Environment string
	Endpoint    string
	Version     string
	DeployTime  time.Duration
	DeployedAt  time.Time
}

// DataCollector collects training data
type DataCollector struct {
	*agent.Agent
	DataSource  string
	RecordLimit int
}

// NewDataCollector creates a new data collector agent
func NewDataCollector(id, name, dataSource string, recordLimit int) *DataCollector {
	return &DataCollector{
		Agent:       agent.NewAgent(id, name, "data_collector"),
		DataSource:  dataSource,
		RecordLimit: recordLimit,
	}
}

// Execute collects data
func (a *DataCollector) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Simulate data collection
	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	result := &DataCollectionResult{
		Source:      a.DataSource,
		RecordCount: a.RecordLimit,
		Features:    20,
		FilePath:    fmt.Sprintf("/data/%s.csv", a.Agent.Name),
		CollectedAt: time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// FeatureExtractor extracts and preprocesses features
type FeatureExtractor struct {
	*agent.Agent
	PreprocessingSteps []string
}

// NewFeatureExtractor creates a new feature extractor agent
func NewFeatureExtractor(id, name string, steps []string) *FeatureExtractor {
	return &FeatureExtractor{
		Agent:              agent.NewAgent(id, name, "feature_extractor"),
		PreprocessingSteps: steps,
	}
}

// Execute extracts features
func (a *FeatureExtractor) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Get data collection result
	dataResult, ok := input.(*DataCollectionResult)
	if !ok {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *DataCollectionResult")
	}

	// Simulate feature extraction
	select {
	case <-time.After(80 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	featureNames := make([]string, 0)
	for i := 0; i < dataResult.Features; i++ {
		featureNames = append(featureNames, fmt.Sprintf("feature_%d", i))
	}

	result := &FeatureExtractionResult{
		FeatureCount:   dataResult.Features,
		FeatureNames:   featureNames,
		ProcessedCount: dataResult.RecordCount,
		Preprocessing:  a.PreprocessingSteps,
		ExtractedAt:    time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// ModelTrainer trains a model with specific hyperparameters
type ModelTrainer struct {
	*agent.Agent
	ModelType   string
	Hyperparams map[string]interface{}
}

// NewModelTrainer creates a new model trainer agent
func NewModelTrainer(id, name, modelType string, hyperparams map[string]interface{}) *ModelTrainer {
	return &ModelTrainer{
		Agent:       agent.NewAgent(id, name, "model_trainer"),
		ModelType:   modelType,
		Hyperparams: hyperparams,
	}
}

// Execute trains the model
func (a *ModelTrainer) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Get feature extraction result
	_, ok := input.(*FeatureExtractionResult)
	if !ok {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *FeatureExtractionResult")
	}

	// Simulate training based on hyperparameters
	var trainingDuration time.Duration
	epochs, _ := a.Hyperparams["epochs"].(int)
	batchSize, _ := a.Hyperparams["batch_size"].(int)
	learningRate, _ := a.Hyperparams["learning_rate"].(float64)

	// Training duration depends on hyperparameters
	durationMultiplier := float64(epochs) / float64(batchSize) * (1.0 / learningRate)
	trainingDuration = time.Duration(durationMultiplier * 50 * float64(time.Millisecond))

	select {
	case <-time.After(trainingDuration):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	// Simulate different accuracies based on hyperparameters
	// Better hyperparameters = better accuracy
	baseAccuracy := 0.85
	accuracyBonus := (float64(epochs) / 100.0) * 0.05
	accuracyBonus += (1.0 / learningRate) * 0.02
	accuracyBonus += (1.0 / float64(batchSize)) * 0.01
	accuracy := math.Min(baseAccuracy+accuracyBonus, 0.98)

	precision := accuracy - 0.02
	recall := accuracy - 0.03
	f1Score := 2 * (precision * recall) / (precision + recall)

	result := &ModelTrainingResult{
		ModelID:      fmt.Sprintf("%s-%d", a.ModelType, time.Now().UnixNano()),
		ModelType:    a.ModelType,
		Hyperparams:  a.Hyperparams,
		Accuracy:     accuracy,
		Precision:    precision,
		Recall:       recall,
		F1Score:      f1Score,
		TrainingTime: time.Since(a.Agent.StartedAt),
		TrainedAt:    time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// ModelEvaluator evaluates a trained model
type ModelEvaluator struct {
	*agent.Agent
	TestSetSize int
}

// NewModelEvaluator creates a new model evaluator agent
func NewModelEvaluator(id, name string, testSetSize int) *ModelEvaluator {
	return &ModelEvaluator{
		Agent:       agent.NewAgent(id, name, "model_evaluator"),
		TestSetSize: testSetSize,
	}
}

// Execute evaluates the model
func (a *ModelEvaluator) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Get model training result
	modelResult, ok := input.(*ModelTrainingResult)
	if !ok {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *ModelTrainingResult")
	}

	// Simulate evaluation
	select {
	case <-time.After(60 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	// Test accuracy is slightly lower than training accuracy
	testAccuracy := modelResult.Accuracy - 0.02
	testLoss := 1.0 - testAccuracy

	confusionMatrix := map[string][]int{
		"actual_positive": {int(testAccuracy * 100), int((1 - testAccuracy) * 100)},
		"actual_negative": {int((1 - testAccuracy) * 100), int(testAccuracy * 100)},
	}

	result := &ModelEvaluationResult{
		ModelID:         modelResult.ModelID,
		TestAccuracy:    testAccuracy,
		TestLoss:        testLoss,
		ConfusionMatrix: confusionMatrix,
		EvaluatedAt:     time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// ModelComparer compares multiple models and selects the best
type ModelComparer struct {
	*agent.Agent
	Metric string // "accuracy", "f1_score", etc.
}

// NewModelComparer creates a new model comparer agent
func NewModelComparer(id, name, metric string) *ModelComparer {
	return &ModelComparer{
		Agent:  agent.NewAgent(id, name, "model_comparer"),
		Metric: metric,
	}
}

// Execute compares models and selects the best
func (a *ModelComparer) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Handle both single and multiple model results
	var modelResults []*ModelTrainingResult

	switch v := input.(type) {
	case *ModelTrainingResult:
		modelResults = []*ModelTrainingResult{v}
	case []*ModelTrainingResult:
		modelResults = v
	case []interface{}:
		modelResults = make([]*ModelTrainingResult, 0)
		for _, item := range v {
			if mr, ok := item.(*ModelTrainingResult); ok {
				modelResults = append(modelResults, mr)
			}
		}
	default:
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *ModelTrainingResult or []*ModelTrainingResult")
	}

	// Simulate comparison
	select {
	case <-time.After(40 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	// Find best model based on metric
	if len(modelResults) == 0 {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("no models to compare")
	}

	bestModel := modelResults[0]
	bestValue := a.getMetricValue(bestModel)

	for _, model := range modelResults[1:] {
		value := a.getMetricValue(model)
		if value > bestValue {
			bestModel = model
			bestValue = value
		}
	}

	result := &ModelComparisonResult{
		BestModelID:  bestModel.ModelID,
		BestAccuracy: bestModel.Accuracy,
		AllModels:    modelResults,
		ComparisonAt: time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// getMetricValue returns the value of the specified metric
func (a *ModelComparer) getMetricValue(model *ModelTrainingResult) float64 {
	switch a.Metric {
	case "accuracy":
		return model.Accuracy
	case "precision":
		return model.Precision
	case "recall":
		return model.Recall
	case "f1_score":
		return model.F1Score
	default:
		return model.Accuracy
	}
}

// ModelDeployer deploys the best model
type ModelDeployer struct {
	*agent.Agent
	Environment  string
	DeployConfig map[string]string
}

// NewModelDeployer creates a new model deployer agent
func NewModelDeployer(id, name, environment string) *ModelDeployer {
	return &ModelDeployer{
		Agent:        agent.NewAgent(id, name, "model_deployer"),
		Environment:  environment,
		DeployConfig: make(map[string]string),
	}
}

// Execute deploys the model
func (a *ModelDeployer) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Get model comparison result
	comparisonResult, ok := input.(*ModelComparisonResult)
	if !ok {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *ModelComparisonResult")
	}

	// Simulate deployment
	select {
	case <-time.After(120 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	result := &ModelDeploymentResult{
		Success:     true,
		ModelID:     comparisonResult.BestModelID,
		Environment: a.Environment,
		Endpoint:    fmt.Sprintf("https://api.%s.example.com/v1/predict", a.Environment),
		Version:     fmt.Sprintf("v1.0.%d", time.Now().Unix()),
		DeployTime:  time.Since(a.Agent.StartedAt),
		DeployedAt:  time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// MLPipeline creates a Machine Learning pipeline workflow
func MLPipeline() (*orchestrator.Orchestrator, error) {
	// Build DAG
	builder := orchestrator.NewWorkflowBuilder()

	// Add nodes
	builder.AddNode("data_collection", "Data Collection", "data_collector")
	builder.AddNode("feature_extraction", "Feature Extraction", "feature_extractor")

	// Model training nodes - train multiple models in parallel with different hyperparameters
	builder.AddNode("train_model1", "Train Model 1 (CNN)", "model_trainer")
	builder.AddNode("train_model2", "Train Model 2 (RNN)", "model_trainer")
	builder.AddNode("train_model3", "Train Model 3 (Transformer)", "model_trainer")
	builder.AddNode("train_model4", "Train Model 4 (Random Forest)", "model_trainer")

	// Model comparison
	builder.AddNode("compare_models", "Compare Models", "model_comparer")

	// Model evaluation nodes
	builder.AddNode("eval_model1", "Evaluate Model 1", "model_evaluator")
	builder.AddNode("eval_model2", "Evaluate Model 2", "model_evaluator")
	builder.AddNode("eval_model3", "Evaluate Model 3", "model_evaluator")
	builder.AddNode("eval_model4", "Evaluate Model 4", "model_evaluator")

	// Model deployment
	builder.AddNode("deploy_model", "Deploy Best Model", "model_deployer")

	// Add edges
	// Data collection -> feature extraction
	builder.AddEdge("data_collection", "feature_extraction", 1.0)

	// Feature extraction -> all model trainings (parallel)
	builder.AddEdge("feature_extraction", "train_model1", 1.0)
	builder.AddEdge("feature_extraction", "train_model2", 1.0)
	builder.AddEdge("feature_extraction", "train_model3", 1.0)
	builder.AddEdge("feature_extraction", "train_model4", 1.0)

	// All model trainings -> comparison (for model selection)
	builder.AddEdge("train_model1", "compare_models", 1.0)
	builder.AddEdge("train_model2", "compare_models", 1.0)
	builder.AddEdge("train_model3", "compare_models", 1.0)
	builder.AddEdge("train_model4", "compare_models", 1.0)

	// Each model training -> corresponding evaluation
	builder.AddEdge("train_model1", "eval_model1", 1.0)
	builder.AddEdge("train_model2", "eval_model2", 1.0)
	builder.AddEdge("train_model3", "eval_model3", 1.0)
	builder.AddEdge("train_model4", "eval_model4", 1.0)

	// Comparison -> deployment
	builder.AddEdge("compare_models", "deploy_model", 1.0)

	// Validate
	if err := builder.Validate(); err != nil {
		return nil, err
	}

	dag := builder.Build()

	// Create orchestrator
	orc := orchestrator.NewOrchestrator(dag)

	// Register agents
	orc.RegisterAgent("data_collection", NewDataCollector("collector-1", "Data Collector", "s3://bucket/data", 10000))
	orc.RegisterAgent("feature_extraction", NewFeatureExtractor("extractor-1", "Feature Extractor", []string{"normalization", "scaling", "encoding"}))

	// Register model trainers with different hyperparameters
	orc.RegisterAgent("train_model1", NewModelTrainer("train-1", "CNN Trainer", "CNN", map[string]interface{}{
		"epochs":        2,
		"batch_size":    32,
		"learning_rate": 0.001,
	}))
	orc.RegisterAgent("train_model2", NewModelTrainer("train-2", "RNN Trainer", "RNN", map[string]interface{}{
		"epochs":        2,
		"batch_size":    64,
		"learning_rate": 0.0005,
	}))
	orc.RegisterAgent("train_model3", NewModelTrainer("train-3", "Transformer Trainer", "Transformer", map[string]interface{}{
		"epochs":        2,
		"batch_size":    128,
		"learning_rate": 0.0001,
	}))
	orc.RegisterAgent("train_model4", NewModelTrainer("train-4", "Random Forest Trainer", "RandomForest", map[string]interface{}{
		"epochs":        1,
		"batch_size":    16,
		"learning_rate": 0.01,
	}))

	// Register model evaluators
	orc.RegisterAgent("eval_model1", NewModelEvaluator("eval-1", "Model 1 Evaluator", 2000))
	orc.RegisterAgent("eval_model2", NewModelEvaluator("eval-2", "Model 2 Evaluator", 2000))
	orc.RegisterAgent("eval_model3", NewModelEvaluator("eval-3", "Model 3 Evaluator", 2000))
	orc.RegisterAgent("eval_model4", NewModelEvaluator("eval-4", "Model 4 Evaluator", 2000))

	// Register model comparer
	orc.RegisterAgent("compare_models", NewModelComparer("compare-1", "Model Comparer", "accuracy"))

	// Register model deployer
	orc.RegisterAgent("deploy_model", NewModelDeployer("deploy-1", "Model Deployer", "production"))

	return orc, nil
}

// PrintMLResults prints ML pipeline results
func PrintMLResults(orc *orchestrator.Orchestrator) {
	fmt.Println("\n=== Machine Learning Pipeline Results ===")

	// Data collection results
	if result, err := orc.GetResult("data_collection"); err == nil {
		if data, ok := result.(*DataCollectionResult); ok {
			fmt.Println("Data Collection:")
			fmt.Printf("  Source: %s\n", data.Source)
			fmt.Printf("  Records: %d\n", data.RecordCount)
			fmt.Printf("  Features: %d\n", data.Features)
			fmt.Printf("  File: %s\n", data.FilePath)
		}
	}

	// Feature extraction results
	if result, err := orc.GetResult("feature_extraction"); err == nil {
		if features, ok := result.(*FeatureExtractionResult); ok {
			fmt.Println("\nFeature Extraction:")
			fmt.Printf("  Features: %d\n", features.FeatureCount)
			fmt.Printf("  Processed: %d records\n", features.ProcessedCount)
			fmt.Printf("  Preprocessing: %v\n", features.Preprocessing)
		}
	}

	// Model training results
	fmt.Println("\nModel Training:")
	modelNodes := []string{"train_model1", "train_model2", "train_model3", "train_model4"}
	for _, nodeID := range modelNodes {
		if result, err := orc.GetResult(nodeID); err == nil {
			if model, ok := result.(*ModelTrainingResult); ok {
				fmt.Printf("  %s:\n", model.ModelType)
				fmt.Printf("    Accuracy: %.2f%%\n", model.Accuracy*100)
				fmt.Printf("    F1 Score: %.2f%%\n", model.F1Score*100)
				fmt.Printf("    Training Time: %v\n", model.TrainingTime)
				fmt.Printf("    Hyperparams: epochs=%d, batch_size=%d, lr=%.4f\n",
					model.Hyperparams["epochs"], model.Hyperparams["batch_size"], model.Hyperparams["learning_rate"])
			}
		}
	}

	// Model evaluation results
	fmt.Println("\nModel Evaluation:")
	evalNodes := []string{"eval_model1", "eval_model2", "eval_model3", "eval_model4"}
	for _, nodeID := range evalNodes {
		if result, err := orc.GetResult(nodeID); err == nil {
			if eval, ok := result.(*ModelEvaluationResult); ok {
				fmt.Printf("  %s: Test Accuracy: %.2f%%, Loss: %.4f\n", eval.ModelID, eval.TestAccuracy*100, eval.TestLoss)
			}
		}
	}

	// Model comparison results
	if result, err := orc.GetResult("compare_models"); err == nil {
		if comparison, ok := result.(*ModelComparisonResult); ok {
			fmt.Println("\nModel Comparison:")
			fmt.Printf("  Best Model: %s\n", comparison.BestModelID)
			fmt.Printf("  Best Accuracy: %.2f%%\n", comparison.BestAccuracy*100)
			fmt.Printf("  Total Models: %d\n", len(comparison.AllModels))

			// Sort models by accuracy
			sortedModels := make([]*ModelTrainingResult, len(comparison.AllModels))
			copy(sortedModels, comparison.AllModels)
			sort.Slice(sortedModels, func(i, j int) bool {
				return sortedModels[i].Accuracy > sortedModels[j].Accuracy
			})

			fmt.Println("\n  Model Ranking:")
			for i, model := range sortedModels {
				fmt.Printf("    %d. %s: %.2f%% accuracy\n", i+1, model.ModelType, model.Accuracy*100)
			}
		}
	}

	// Model deployment results
	if result, err := orc.GetResult("deploy_model"); err == nil {
		if deploy, ok := result.(*ModelDeploymentResult); ok {
			fmt.Println("\nModel Deployment:")
			fmt.Printf("  Status: %s\n", map[bool]string{true: "SUCCESS", false: "FAILED"}[deploy.Success])
			fmt.Printf("  Model ID: %s\n", deploy.ModelID)
			fmt.Printf("  Environment: %s\n", deploy.Environment)
			fmt.Printf("  Endpoint: %s\n", deploy.Endpoint)
			fmt.Printf("  Version: %s\n", deploy.Version)
			fmt.Printf("  Deploy Time: %v\n", deploy.DeployTime)
		}
	}

	fmt.Println("\n====================================")
}
