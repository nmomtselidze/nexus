# Multi-Agent Graph System: The Ultimate Workflow Orchestration Framework

Multi-Agent Graph System is a groundbreaking Go-based framework designed to orchestrate complex multi-agent workflows using Directed Acyclic Graphs (DAGs). From CI/CD pipelines to machine learning workflows, this system provides unparalleled efficiency and intelligence for developers, data engineers, DevOps engineers, and system architects alike.

## 🌟 Overview

Multi-Agent Graph System is a powerful orchestration framework that models workflows as DAGs, where nodes represent autonomous agents and edges define dependencies. With support for parallel execution, custom agents, and built-in workflow patterns, it's the ultimate tool for automating complex business processes.

## 🚀 Features

- **Graph-Based Orchestration**: Full DAG implementation with cycle detection and topological sorting
- **Parallel Execution**: Execute independent agents concurrently for maximum performance
- **Agent Framework**: Extensible agent types with custom execution logic
- **Workflow Builder**: Intuitive DSL for building complex workflows
- **Built-in Agent Types**: Data processors, transformers, validators, and aggregators
- **Real-time Execution Monitoring**: Track agent status, execution time, and results
- **Result Management**: Comprehensive result tracking and error handling
- **Context-Based Cancellation**: Stop workflows gracefully with context cancellation
- **Cross-Platform Compatibility**: Works seamlessly across Linux, macOS, and Windows
- **Comprehensive Testing**: Unit tests, integration tests, and 49.9% overall coverage
- **Makefile Automation**: Convenient build, test, and run commands
- **CLI Interface**: Easy-to-use command-line interface for all workflows

## 🛠️ Installation

### Prerequisites

- Go 1.27.1 or higher
- Git (for cloning the repository)

### Steps

**Clone the repository:**
```bash
git clone <repository-url>
cd multi-agent-graph
```

**Build the CLI:**
```bash
# Using Makefile (recommended)
make build

# Or build directly
go build -o multi-agent-graph examples/cmd/main.go
```

**Install to GOPATH (optional):**
```bash
make install
```

**Run the application:**
```bash
# Run a workflow
./multi-agent-graph ml --parallel

# Or use Makefile
make run-ml
```

## 🤖 Core Components

### Graph Engine

The `graph/` package provides the foundation for DAG-based workflows:

```go
import "multi-agent-graph/graph"

dag := graph.NewDAG()
node := graph.NewNode("node1", "First Node", "data_processor")
dag.AddNode(node)
dag.AddEdge("node1", "node2", 1.0)
sorted, _ := dag.TopologicalSort()
```

**Key Features:**
- Cycle detection using DFS
- Topological sorting with Kahn's algorithm
- Root and leaf node identification
- Dependency tracking

### Agent Framework

The `agent/` package provides extensible agent types:

```go
import "multi-agent-graph/agent"

// Built-in agents
processor := agent.NewDataProcessorAgent("id", "Processor")
transformer := agent.NewTransformAgent("id", "Transformer")
validator := agent.NewValidationAgent("id", "Validator")
aggregator := agent.NewAggregatorAgent("id", "Aggregator")

// Custom agents
type CustomAgent struct {
    *agent.Agent
    CustomFunc func(interface{}) (interface{}, error)
}

func (a *CustomAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    // Your custom logic
    return output, nil
}
```

**Built-in Agent Types:**
- **DataProcessorAgent**: Processes and transforms data
- **TransformAgent**: Applies transformations to data
- **ValidationAgent**: Validates data integrity
- **AggregatorAgent**: Aggregates results from multiple sources

### Orchestrator

The `orchestrator/` package manages workflow execution:

```go
import "multi-agent-graph/orchestrator"

builder := orchestrator.NewWorkflowBuilder()
builder.AddNode("node1", "Node 1", "data_processor")
builder.AddNode("node2", "Node 2", "transform")
builder.AddEdge("node1", "node2", 1.0)

dag := builder.Build()
orc := orchestrator.NewOrchestrator(dag)

// Register agents
orc.RegisterAgent("node1", agent.NewDataProcessorAgent("id", "Processor"))

// Execute sequentially
orc.Execute(input)

// Or execute in parallel
orc.ExecuteParallel(input)

// Get results
result, _ := orc.GetResult("node1")
results := orc.GetAllResults()
```

**Execution Modes:**
- **Sequential**: Executes nodes one at a time in topological order
- **Parallel**: Executes independent nodes concurrently for better performance

## 📚 Example Workflows

### Linear Workflow

Simple sequential processing pipeline:

```
validate -> process -> transform -> aggregate
```

**Run it:**
```bash
./multi-agent-graph linear --parallel
```

### Parallel Workflow

Parallel branch processing with merge:

```
        -> branch1 ->
input  -> branch2 -> merge -> output
        -> branch3 ->
```

**Run it:**
```bash
./multi-agent-graph parallel --parallel
```

### ETL Pipeline

Real-world ETL workflow extracting from 3 databases:

```
db1_extract -> db1_validate \
                              -> merge -> warehouse_load
db2_extract -> db2_validate /
db3_extract -> db3_validate /
```

**Features:**
- Extracts data from PostgreSQL, MySQL, MongoDB
- Validates each source independently (parallel)
- Merges all validated data
- Loads to data warehouse (Snowflake) in batches

**Run it:**
```bash
./multi-agent-graph etl --parallel
```

**Results include:**
- Database extraction statistics
- Validation results with issue counts
- Merge statistics
- Warehouse load status and timing

### CI/CD Pipeline

Modern CI/CD workflow with multi-platform testing:

```
        build
       /  |  \
      /   |   \
test_ test_ test_ test_ security lint
linux windows macos integration  scan
      \   |   /      |        |      |
       \  |  /       |        |      |
        deploy_staging
```

**Features:**
- Builds the application
- Runs tests on multiple platforms in parallel (Linux, Windows, macOS)
- Executes integration tests
- Scans for security vulnerabilities
- Lints code for quality
- Deploys to staging (conditional on all tests passing)

**Run it:**
```bash
./multi-agent-graph cicd --parallel
```

**Results include:**
- Build artifact information
- Test results and coverage for all platforms
- Security scan vulnerability report
- Lint results with error/warning counts
- Deployment status and URL

### ML Pipeline

Machine learning pipeline with parallel model training:

```
data_collection -> feature_extraction
                                /-> train_model1 -> eval_model1
                               /-> train_model2 -> eval_model2
                              /-> train_model3 -> eval_model3
                             /-> train_model4 -> eval_model4
                            /
                           compare_models -> deploy_model
```

**Features:**
- Collects training data (10,000 records, 20 features)
- Extracts and preprocesses features
- Trains 4 models in parallel (CNN, RNN, Transformer, Random Forest)
- Evaluates each model independently
- Compares models and selects the best
- Deploys the best model to production

**Run it:**
```bash
./multi-agent-graph ml --parallel
```

**Results include:**
- Data collection statistics
- Feature extraction details
- Model training results (accuracy, F1 score, training time, hyperparameters)
- Model evaluation results (test accuracy, loss)
- Model comparison and ranking
- Deployment status and endpoint

### Flight Booking Workflow

Autonomous flight booking with airport resolution:

```
resolve_airports -> search_flights -> reserve_seat -> process_payment -> issue_ticket -> complete_booking
```

**Features:**
- Resolves city names to airport codes via geocoder API
- Searches available flight schedules
- Reserves seat and creates PNR
- Processes payment
- Issues e-ticket document
- Aggregates complete booking information

**Run it:**
```bash
./multi-agent-graph book-flight --parallel
```

**Results include:**
- Airport codes with multilingual support
- Flight details (number, airline, partner)
- Ticket number and price
- Booking status

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with coverage (no HTML)
make test-coverage-short

# Run unit tests only
make test-unit
```

### Test Coverage

- **agent**: 82.2% coverage
- **orchestrator**: 85.4% coverage
- **graph**: 69.2% coverage
- **examples**: 63.7% coverage
- **Overall**: 49.9% coverage

### Test Files

- `agent/agent_test.go` - Agent creation, configuration, execution tests
- `agent/registry_test.go` - Registry management tests
- `graph/dag_test.go` - Graph operations, cycle detection, topological sort tests
- `orchestrator/orchestrator_test.go` - Workflow execution, parallel execution tests
- `orchestrator/builder_test.go` - Workflow builder and validation tests
- `examples/examples_integration_test.go` - Integration tests for all example workflows

## 🎨 Design Patterns

### DAG-Based Orchestration

Workflows are modeled as Directed Acyclic Graphs where:
- Nodes represent agents/tasks
- Edges represent dependencies
- Topological sort determines execution order
- Cycle detection prevents circular dependencies

### Agent Pattern

Each agent:
- Encapsulates specific processing logic
- Maintains its own state and metadata
- Can be independently tested and reused
- Supports custom execution functions

### Builder Pattern

The WorkflowBuilder provides a fluent interface for:
- Declarative workflow construction
- Validation before execution
- Separation of definition and execution

## 🛠️ Makefile Commands

### Build Commands

```bash
make build          # Build the CLI binary
make install        # Install to GOPATH/bin
make clean          # Remove build artifacts
```

### Test Commands

```bash
make test                  # Run all tests
make test-coverage         # Run tests with HTML coverage report
make test-coverage-short   # Run tests with coverage (no HTML)
make test-unit             # Run unit tests only
```

### Code Quality

```bash
make fmt            # Format Go code
make lint           # Run linter (requires golangci-lint)
make ci             # Run CI checks (fmt, test, lint)
```

### Dependencies

```bash
make deps           # Download dependencies
make deps-tidy      # Tidy and verify dependencies
```

### Workflow Runners

```bash
make run-linear          # Run linear workflow
make run-parallel        # Run parallel workflow
make run-complex         # Run complex workflow
make run-etl             # Run ETL pipeline
make run-cicd            # Run CI/CD pipeline
make run-ml              # Run ML pipeline
make run-book-flight     # Run flight booking workflow
make run-all             # Run all workflows
```

### Help

```bash
make help           # Display all available commands
```

## ⚡ Performance Considerations

- **Sequential Mode**: Simpler, predictable execution order
- **Parallel Mode**: Faster for independent nodes, but higher resource usage
- **Cycle Detection**: O(V + E) complexity using DFS
- **Topological Sort**: O(V + E) complexity using Kahn's algorithm

## 🔧 Extending the System

### Adding New Agent Types

1. Create a new agent struct embedding `*agent.Agent`
2. Implement the `Execute(ctx, input)` method
3. Register the agent with the orchestrator
4. Add type handling in the orchestrator's `executeAgent` method

```go
type MyCustomAgent struct {
    *agent.Agent
    CustomField string
}

func (a *MyCustomAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    a.Agent.Status = "running"
    a.Agent.StartedAt = time.Now()
    
    // Your custom logic
    output, err := a.process(input)
    
    if err != nil {
        a.Agent.Status = "failed"
        a.Agent.CompletedAt = time.Now()
        return nil, err
    }
    
    a.Agent.Status = "completed"
    a.Agent.CompletedAt = time.Now()
    return output, nil
}
```

### Adding Graph Features

The graph engine can be extended with:
- Edge weights for priority-based execution
- Conditional edges for dynamic workflows
- Sub-graphs for modular workflows
- Graph visualization utilities

## 🛡️ Security

Multi-Agent Graph System employs robust security mechanisms:

- **Context-Based Cancellation**: Prevent runaway workflows
- **Agent Isolation**: Each agent operates independently
- **Error Handling**: Comprehensive error tracking and reporting
- **Input Validation**: Built-in validation agents for data integrity

## 📊 Analytics

The orchestrator provides detailed execution analytics:

- **Agent Status Tracking**: Real-time status for all agents
- **Execution Duration**: Time tracking for each agent
- **Result Storage**: Comprehensive result management
- **Error Logging**: Detailed error tracking and reporting

## 🌐 Integration

Multi-Agent Graph System integrates with:
- **Go Standard Library**: Native Go integration
- **Custom APIs**: Easy integration with external services
- **Cloud Platforms**: Can be extended for AWS, GCP, Azure workflows
- **Databases**: ETL examples show PostgreSQL, MySQL, MongoDB integration

## 🤝 Contributing

We welcome contributions! To get started:

1. Fork the repository
2. Create a new branch (`git checkout -b feature-name`)
3. Commit your changes (`git commit -m 'Add feature'`)
4. Push to the branch (`git push origin feature-name`)
5. Open a pull request

### Development Workflow

```bash
# Create feature branch
git checkout -b feature/my-feature

# Make your changes
# ...

# Run tests
make test

# Run linting
make lint

# Format code
make fmt

# Commit and push
git add .
git commit -m "Add my feature"
git push origin feature/my-feature
```

## 🐛 Troubleshooting

### Common Issues

**Build fails:**
```bash
# Ensure Go is installed
go version

# Update dependencies
make deps-tidy
```

**Tests fail:**
```bash
# Run tests with verbose output
make test

# Check specific package
go test ./agent -v
```

**Workflow execution errors:**
```bash
# Check DAG structure
./multi-agent-graph <workflow-type>

# Verify agents are registered
# Check agent implementation
```

## 💡 Inspiration

Multi-Agent Graph System is inspired by a vision to create the most advanced, efficient workflow orchestration framework capable of handling any complex task with precision and intelligence, leveraging the power of Directed Acyclic Graphs and autonomous agents.

## ✨ Future Goals

- **Web Dashboard**: Real-time workflow visualization and monitoring
- **GraphQL API**: RESTful API for workflow management
- **Persistent Storage**: Database integration for workflow history
- **Advanced Scheduling**: Cron-based workflow execution
- **Retry Mechanisms**: Automatic retry on failure with exponential backoff
- **Workflow Templates**: Pre-built templates for common patterns
- **Distributed Execution**: Support for multi-machine execution
- **Event-Driven Triggers**: Start workflows based on events
- **Kubernetes Integration**: Native Kubernetes operator
- **Enhanced Monitoring**: Prometheus metrics and Grafana dashboards

## 📝 License

[Add your license here]

## 📞 Support

For issues, questions, or contributions, please open an issue on the repository.

## 🎓 Conclusion

Multi-Agent Graph System is more than just a workflow framework—it's your ultimate tool for orchestrating complex multi-agent processes with unparalleled efficiency, intelligence, and flexibility. Whether you're building CI/CD pipelines, ML workflows, or data processing pipelines, this system provides the foundation for automation at scale.

---

Built with ❤️ using Go and Directed Acyclic Graphs
