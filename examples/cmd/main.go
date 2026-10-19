package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"multi-agent-graph/examples"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	workflowType := os.Args[1]
	parallel := len(os.Args) > 2 && os.Args[2] == "--parallel"

	switch workflowType {
	case "linear":
		runLinearWorkflow(parallel)
	case "parallel":
		runParallelWorkflow(parallel)
	case "complex":
		runComplexWorkflow(parallel)
	case "etl":
		runETLWorkflow(parallel)
	case "cicd":
		runCICDWorkflow(parallel)
	case "ml":
		runMLWorkflow(parallel)
	case "book-flight", "book_flight":
		runFlightBookingWorkflow(parallel)
	default:
		fmt.Printf("Unknown workflow type: %s\n", workflowType)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Multi-Agent Graph Workflow System")
	fmt.Println("\nUsage:")
	fmt.Println("  go run examples/cmd/main.go <workflow-type> [--parallel]")
	fmt.Println("\nWorkflow Types:")
	fmt.Println("  linear      - Linear workflow")
	fmt.Println("  parallel    - Parallel branches workflow")
	fmt.Println("  complex     - Complex DAG workflow")
	fmt.Println("  etl         - ETL pipeline")
	fmt.Println("  cicd        - CI/CD pipeline")
	fmt.Println("  ml          - ML pipeline")
	fmt.Println("  book-flight - Autonomous flight booking (resolves airport, schedule, and ticket)")
}

func runFlightBookingWorkflow(parallel bool) {
	fmt.Println("=== Autonomous Flight Booking Workflow ===")

	orc, err := examples.FlightBookingWorkflow()
	if err != nil {
		log.Fatalf("Failed to create flight booking workflow: %v", err)
	}

	examples.PrintDAG(orc.DAG)

	// User input contains ONLY what is known:
	input := map[string]interface{}{
		"origin_city": "Tbilisi",
		"dest_city":   "Madrid",
		"date":        "2026-11-15",
	}

	fmt.Println("User Input:")
	fmt.Printf("  %+v\n\n", input)

	// Silence internal orchestrator log output
	log.SetOutput(io.Discard)

	start := time.Now()
	var execErr error
	if parallel {
		execErr = orc.ExecuteParallel(input)
	} else {
		execErr = orc.Execute(input)
	}
	duration := time.Since(start)

	// Restore logger output
	log.SetOutput(os.Stderr)

	if execErr != nil {
		log.Fatalf("Workflow execution failed: %v", execErr)
	}

	fmt.Printf("Execution completed in %v\n\n", duration)

	results := orc.GetAllResults()
	booking, ok := results["complete_booking"].(map[string]interface{})
	if !ok {
		booking, _ = results["issue_ticket"].(map[string]interface{})
	}

	if booking == nil {
		booking = make(map[string]interface{})
	}

	fmt.Println("==================================================")
	fmt.Println("            FLIGHT TICKET / АВИАБИЛЕТ             ")
	fmt.Println("==================================================")
	fmt.Printf("City (Город):        %s -> %s\n", booking["origin_city"], booking["dest_city"])
	fmt.Printf("Airport (Аэропорт):  %s -> %s\n", booking["origin_airport"], booking["dest_airport"])
	fmt.Printf("Date (Дата):         %s\n", booking["date"])
	fmt.Printf("Flight (Рейс):       %s (%s)\n", booking["flight"], booking["airline"])
	fmt.Printf("Price (Цена):        %s\n", booking["price"])
	fmt.Printf("Ticket (Номер):      %s\n", booking["ticket_number"])
	fmt.Printf("Status (Статус):     %s\n", booking["status"])
	fmt.Println("==================================================")

	fmt.Println("\n=== Node Status ===")
	for nodeID, status := range orc.GetStatus() {
		fmt.Printf("%s: %s\n", nodeID, status)
	}
}
func runLinearWorkflow(parallel bool) {
	orc, err := examples.LinearWorkflow()
	if err != nil {
		log.Fatalf("Failed to create workflow: %v", err)
	}
	examples.PrintDAG(orc.DAG)
	input := map[string]interface{}{"data": "sample data", "timestamp": time.Now().Format(time.RFC3339)}
	if parallel {
		orc.ExecuteParallel(input)
	} else {
		orc.Execute(input)
	}
}

func runParallelWorkflow(parallel bool) {
	orc, err := examples.ParallelWorkflow()
	if err != nil {
		log.Fatalf("Failed to create workflow: %v", err)
	}
	examples.PrintDAG(orc.DAG)
	input := map[string]interface{}{"data": "sample data", "timestamp": time.Now().Format(time.RFC3339)}
	if parallel {
		orc.ExecuteParallel(input)
	} else {
		orc.Execute(input)
	}
}

func runComplexWorkflow(parallel bool) {
	orc, err := examples.ComplexWorkflow()
	if err != nil {
		log.Fatalf("Failed to create workflow: %v", err)
	}
	examples.PrintDAG(orc.DAG)
	input := map[string]interface{}{"data": "sample data", "timestamp": time.Now().Format(time.RFC3339)}
	if parallel {
		orc.ExecuteParallel(input)
	} else {
		orc.Execute(input)
	}
}

func runETLWorkflow(parallel bool) {
	orc, err := examples.ETLPipeline()
	if err != nil {
		log.Fatalf("Failed to create workflow: %v", err)
	}
	examples.PrintDAG(orc.DAG)
	if parallel {
		orc.ExecuteParallel(nil)
	} else {
		orc.Execute(nil)
	}
	examples.PrintETLResults(orc)
}

func runCICDWorkflow(parallel bool) {
	orc, err := examples.CICDPipeline()
	if err != nil {
		log.Fatalf("Failed to create workflow: %v", err)
	}
	examples.PrintDAG(orc.DAG)
	if parallel {
		orc.ExecuteParallel(nil)
	} else {
		orc.Execute(nil)
	}
	examples.PrintCICDResults(orc)
}

func runMLWorkflow(parallel bool) {
	orc, err := examples.MLPipeline()
	if err != nil {
		log.Fatalf("Failed to create workflow: %v", err)
	}
	examples.PrintDAG(orc.DAG)
	if parallel {
		orc.ExecuteParallel(nil)
	} else {
		orc.Execute(nil)
	}
	examples.PrintMLResults(orc)
}
