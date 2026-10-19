package examples

import (
	"context"
	"fmt"
	"time"

	"multi-agent-graph/agent"
	"multi-agent-graph/orchestrator"
)

// BuildResult represents build output
type BuildResult struct {
	Success    bool
	ArtifactID string
	BinaryPath string
	BuildTime  time.Duration
	BuiltAt    time.Time
}

// TestResult represents test output
type TestResult struct {
	Platform   string
	TotalTests int
	Passed     int
	Failed     int
	Skipped    int
	Coverage   float64
	Success    bool
	TestTime   time.Duration
	TestedAt   time.Time
}

// LintResult represents lint output
type LintResult struct {
	Success     bool
	IssuesFound int
	Warnings    int
	Errors      int
	LintTime    time.Duration
	LintedAt    time.Time
}

// SecurityScanResult represents security scan output
type SecurityScanResult struct {
	Success         bool
	Vulnerabilities int
	Critical        int
	High            int
	Medium          int
	Low             int
	ScanTime        time.Duration
	ScannedAt       time.Time
}

// DeployResult represents deployment output
type DeployResult struct {
	Success     bool
	Environment string
	Version     string
	URL         string
	DeployTime  time.Duration
	DeployedAt  time.Time
}

// BuildAgent builds the application
type BuildAgent struct {
	*agent.Agent
	ProjectPath string
	BuildConfig map[string]string
}

// NewBuildAgent creates a new build agent
func NewBuildAgent(id, name, projectPath string) *BuildAgent {
	return &BuildAgent{
		Agent:       agent.NewAgent(id, name, "build"),
		ProjectPath: projectPath,
		BuildConfig: make(map[string]string),
	}
}

// Execute builds the application
func (a *BuildAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Simulate build process
	select {
	case <-time.After(200 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	result := &BuildResult{
		Success:    true,
		ArtifactID: fmt.Sprintf("artifact-%d", time.Now().Unix()),
		BinaryPath: fmt.Sprintf("/builds/%s/app-binary", a.Agent.Name),
		BuildTime:  time.Since(a.Agent.StartedAt),
		BuiltAt:    time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// TestAgent runs tests on a specific platform
type TestAgent struct {
	*agent.Agent
	Platform string
	TestType string // "unit", "integration", "e2e"
}

// NewTestAgent creates a new test agent
func NewTestAgent(id, name, platform, testType string) *TestAgent {
	return &TestAgent{
		Agent:    agent.NewAgent(id, name, "test"),
		Platform: platform,
		TestType: testType,
	}
}

// Execute runs tests
func (a *TestAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Get build result
	_, ok := input.(*BuildResult)
	if !ok {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *BuildResult")
	}

	// Simulate test execution based on platform
	var testDuration time.Duration
	switch a.Platform {
	case "linux":
		testDuration = 150 * time.Millisecond
	case "windows":
		testDuration = 180 * time.Millisecond
	case "macos":
		testDuration = 160 * time.Millisecond
	default:
		testDuration = 140 * time.Millisecond
	}

	select {
	case <-time.After(testDuration):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	// Simulate test results
	totalTests := 100
	passed := totalTests // All tests pass
	failed := 0
	skipped := 0
	coverage := 85.5

	result := &TestResult{
		Platform:   a.Platform,
		TotalTests: totalTests,
		Passed:     passed,
		Failed:     failed,
		Skipped:    skipped,
		Coverage:   coverage,
		Success:    failed == 0,
		TestTime:   time.Since(a.Agent.StartedAt),
		TestedAt:   time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// LintAgent runs linting
type LintAgent struct {
	*agent.Agent
	Linter string
	Rules  []string
}

// NewLintAgent creates a new lint agent
func NewLintAgent(id, name, linter string) *LintAgent {
	return &LintAgent{
		Agent:  agent.NewAgent(id, name, "lint"),
		Linter: linter,
		Rules:  []string{"golint", "staticcheck", "vet"},
	}
}

// Execute runs linting
func (a *LintAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Simulate linting
	select {
	case <-time.After(80 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	result := &LintResult{
		Success:     true,
		IssuesFound: 5,
		Warnings:    3,
		Errors:      0,
		LintTime:    time.Since(a.Agent.StartedAt),
		LintedAt:    time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// SecurityScanAgent runs security scanning
type SecurityScanAgent struct {
	*agent.Agent
	Scanner string
}

// NewSecurityScanAgent creates a new security scan agent
func NewSecurityScanAgent(id, name, scanner string) *SecurityScanAgent {
	return &SecurityScanAgent{
		Agent:   agent.NewAgent(id, name, "security_scan"),
		Scanner: scanner,
	}
}

// Execute runs security scan
func (a *SecurityScanAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Get build result
	_, ok := input.(*BuildResult)
	if !ok {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *BuildResult")
	}

	// Simulate security scan
	select {
	case <-time.After(120 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	result := &SecurityScanResult{
		Success:         true,
		Vulnerabilities: 2,
		Critical:        0,
		High:            1,
		Medium:          1,
		Low:             0,
		ScanTime:        time.Since(a.Agent.StartedAt),
		ScannedAt:       time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// DeployAgent deploys the application
type DeployAgent struct {
	*agent.Agent
	Environment  string
	DeployConfig map[string]string
}

// NewDeployAgent creates a new deploy agent
func NewDeployAgent(id, name, environment string) *DeployAgent {
	return &DeployAgent{
		Agent:        agent.NewAgent(id, name, "deploy"),
		Environment:  environment,
		DeployConfig: make(map[string]string),
	}
}

// Execute deploys the application
func (a *DeployAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	a.Agent.Status = "running"
	a.Agent.StartedAt = time.Now()

	// Handle both single test result and multiple test results
	var allTestsPassed bool
	var testResults []*TestResult

	switch v := input.(type) {
	case *TestResult:
		testResults = []*TestResult{v}
		allTestsPassed = v.Success
	case []*TestResult:
		testResults = v
		allTestsPassed = true
		for _, tr := range testResults {
			if !tr.Success {
				allTestsPassed = false
				break
			}
		}
	case []interface{}:
		testResults = make([]*TestResult, 0)
		allTestsPassed = true
		for _, item := range v {
			if tr, ok := item.(*TestResult); ok {
				testResults = append(testResults, tr)
				if !tr.Success {
					allTestsPassed = false
				}
			}
		}
	default:
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("invalid input type, expected *TestResult or []*TestResult")
	}

	// Conditional deployment - only deploy if all tests pass
	if !allTestsPassed {
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, fmt.Errorf("deployment skipped: tests failed")
	}

	// Simulate deployment
	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		a.Agent.Status = "failed"
		a.Agent.CompletedAt = time.Now()
		return nil, ctx.Err()
	}

	result := &DeployResult{
		Success:     true,
		Environment: a.Environment,
		Version:     fmt.Sprintf("v1.0.%d", time.Now().Unix()),
		URL:         fmt.Sprintf("https://%s.example.com", a.Environment),
		DeployTime:  time.Since(a.Agent.StartedAt),
		DeployedAt:  time.Now(),
	}

	a.Agent.Status = "completed"
	a.Agent.CompletedAt = time.Now()
	return result, nil
}

// CICDPipeline creates a CI/CD pipeline workflow
func CICDPipeline() (*orchestrator.Orchestrator, error) {
	// Build DAG
	builder := orchestrator.NewWorkflowBuilder()

	// Add nodes
	builder.AddNode("build", "Build Application", "build")

	// Test nodes - run in parallel on different platforms
	builder.AddNode("test_linux", "Unit Tests (Linux)", "test")
	builder.AddNode("test_windows", "Unit Tests (Windows)", "test")
	builder.AddNode("test_macos", "Unit Tests (macOS)", "test")
	builder.AddNode("test_integration", "Integration Tests", "test")
	builder.AddNode("security_scan", "Security Scan", "security_scan")

	// Lint node
	builder.AddNode("lint", "Code Linting", "lint")

	// Deploy node
	builder.AddNode("deploy_staging", "Deploy to Staging", "deploy")

	// Add edges
	// Build -> all tests and security scan (parallel)
	builder.AddEdge("build", "test_linux", 1.0)
	builder.AddEdge("build", "test_windows", 1.0)
	builder.AddEdge("build", "test_macos", 1.0)
	builder.AddEdge("build", "test_integration", 1.0)
	builder.AddEdge("build", "security_scan", 1.0)

	// Build -> lint
	builder.AddEdge("build", "lint", 1.0)

	// All tests -> deploy (conditional)
	builder.AddEdge("test_linux", "deploy_staging", 1.0)
	builder.AddEdge("test_windows", "deploy_staging", 1.0)
	builder.AddEdge("test_macos", "deploy_staging", 1.0)
	builder.AddEdge("test_integration", "deploy_staging", 1.0)

	// Security scan -> deploy
	builder.AddEdge("security_scan", "deploy_staging", 1.0)

	// Lint -> deploy
	builder.AddEdge("lint", "deploy_staging", 1.0)

	// Validate
	if err := builder.Validate(); err != nil {
		return nil, err
	}

	dag := builder.Build()

	// Create orchestrator
	orc := orchestrator.NewOrchestrator(dag)

	// Register agents
	orc.RegisterAgent("build", NewBuildAgent("build-1", "Build Agent", "/app"))
	orc.RegisterAgent("test_linux", NewTestAgent("test-linux", "Linux Tests", "linux", "unit"))
	orc.RegisterAgent("test_windows", NewTestAgent("test-windows", "Windows Tests", "windows", "unit"))
	orc.RegisterAgent("test_macos", NewTestAgent("test-macos", "macOS Tests", "macos", "unit"))
	orc.RegisterAgent("test_integration", NewTestAgent("test-integration", "Integration Tests", "linux", "integration"))
	orc.RegisterAgent("security_scan", NewSecurityScanAgent("security-1", "Security Scanner", "snyk"))
	orc.RegisterAgent("lint", NewLintAgent("lint-1", "Linter", "golangci-lint"))
	orc.RegisterAgent("deploy_staging", NewDeployAgent("deploy-1", "Staging Deployer", "staging"))

	return orc, nil
}

// PrintCICDResults prints CI/CD pipeline results
func PrintCICDResults(orc *orchestrator.Orchestrator) {
	fmt.Println("\n=== CI/CD Pipeline Results ===")

	// Build results
	if result, err := orc.GetResult("build"); err == nil {
		if build, ok := result.(*BuildResult); ok {
			fmt.Println("Build:")
			fmt.Printf("  Status: %s\n", map[bool]string{true: "SUCCESS", false: "FAILED"}[build.Success])
			fmt.Printf("  Artifact ID: %s\n", build.ArtifactID)
			fmt.Printf("  Binary Path: %s\n", build.BinaryPath)
			fmt.Printf("  Build Time: %v\n", build.BuildTime)
		}
	}

	// Test results
	fmt.Println("\nTest Results:")
	testNodes := []string{"test_linux", "test_windows", "test_macos", "test_integration"}
	allTestsPassed := true
	for _, nodeID := range testNodes {
		if result, err := orc.GetResult(nodeID); err == nil {
			if test, ok := result.(*TestResult); ok {
				status := "PASSED"
				if !test.Success {
					status = "FAILED"
					allTestsPassed = false
				}
				fmt.Printf("  %s: %s (%d/%d passed, %.1f%% coverage)\n",
					test.Platform, status, test.Passed, test.TotalTests, test.Coverage)
			}
		}
	}
	fmt.Printf("  Overall: %s\n", map[bool]string{true: "ALL TESTS PASSED", false: "SOME TESTS FAILED"}[allTestsPassed])

	// Security scan results
	if result, err := orc.GetResult("security_scan"); err == nil {
		if scan, ok := result.(*SecurityScanResult); ok {
			fmt.Println("\nSecurity Scan:")
			fmt.Printf("  Status: %s\n", map[bool]string{true: "SUCCESS", false: "FAILED"}[scan.Success])
			fmt.Printf("  Vulnerabilities: %d (Critical: %d, High: %d, Medium: %d, Low: %d)\n",
				scan.Vulnerabilities, scan.Critical, scan.High, scan.Medium, scan.Low)
			fmt.Printf("  Scan Time: %v\n", scan.ScanTime)
		}
	}

	// Lint results
	if result, err := orc.GetResult("lint"); err == nil {
		if lint, ok := result.(*LintResult); ok {
			fmt.Println("\nLint:")
			fmt.Printf("  Status: %s\n", map[bool]string{true: "SUCCESS", false: "FAILED"}[lint.Success])
			fmt.Printf("  Issues Found: %d (Warnings: %d, Errors: %d)\n",
				lint.IssuesFound, lint.Warnings, lint.Errors)
			fmt.Printf("  Lint Time: %v\n", lint.LintTime)
		}
	}

	// Deploy results
	if result, err := orc.GetResult("deploy_staging"); err == nil {
		if deploy, ok := result.(*DeployResult); ok {
			fmt.Println("\nDeployment:")
			fmt.Printf("  Status: %s\n", map[bool]string{true: "SUCCESS", false: "FAILED"}[deploy.Success])
			fmt.Printf("  Environment: %s\n", deploy.Environment)
			fmt.Printf("  Version: %s\n", deploy.Version)
			fmt.Printf("  URL: %s\n", deploy.URL)
			fmt.Printf("  Deploy Time: %v\n", deploy.DeployTime)
		}
	} else {
		fmt.Println("\nDeployment:")
		fmt.Printf("  Status: SKIPPED (tests or checks failed)\n")
	}

	fmt.Println("\n==============================")
}
