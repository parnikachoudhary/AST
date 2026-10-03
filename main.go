package main

import(
	"ast-engine/pkg/parser"
	"fmt"
)
 

func main() {
	
	targetDir := "./test_project"

	fmt.Printf("[1] Scanning C/C++ repository at: %s\n", targetDir)
	rawGraph, err := parser.ScanDirectory(targetDir)
	if err != nil {
		fmt.Printf("[ERROR] Directory scan failed: %v\n", err)
		return
	}

	fmt.Println("\nRaw Include Dependencies Found:")
	for node, neighbors := range rawGraph.AdjacencyList {
		fmt.Printf("  %s includes -> %v\n", node, neighbors)
	}

	// Apply Transitive Reduction
	fmt.Println("\n[2] Applying Transitive Reduction...")
	reducedGraph := rawGraph.TransitiveReduction()

	// Export Graph
	outputFile := "file_dependencies_go.dot"
	err = reducedGraph.ExportToDOT(outputFile)
	if err != nil {
		fmt.Printf("[ERROR] Export failed: %v\n", err)
		return
	}

	fmt.Printf("\n[SUCCESS] Engine successfully exported architecture to: %s\n", outputFile)

}
