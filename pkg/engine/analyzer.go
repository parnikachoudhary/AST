package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	

	"ast-engine/pkg/ast"
	"ast-engine/pkg/graph"
	"ast-engine/pkg/parser"
)

type AnalysisResult struct {
	FileGraph   *graph.Graph
	SymbolGraph *graph.Graph
	SymbolTable *ast.SymbolTable
	FileCount int
	
}


func isSourceFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".cpp" || ext == ".hpp" || ext == ".cc" || ext == ".c" || ext == ".h"
}

func AnalyzeRepo(targetDir string) (*AnalysisResult, error) {
	fmt.Printf("[ENGINE] Scanning directory tree recursively at: %s\n", targetDir)

	// 1. Concurrent File Dependency Scan
	fileG, err := parser.ScanDirectoryConcurrent(targetDir)
	if err != nil {
		return nil, fmt.Errorf("file scan failed: %w", err)
	}

	symbolTable := ast.NewSymbolTable()
	symbolG := graph.NewGraph()

	// 2. Discover all C/C++ source files dynamically
	var sourceFiles []string
	err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isSourceFile(path) {
			sourceFiles = append(sourceFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed walking directory: %w", err)
	}

	fmt.Printf("[ENGINE] Found %d C/C++ source files for AST indexing.\n", len(sourceFiles))

	// 3. Extract Symbols across all discovered files in parallel (Goroutines)

	var wg sync.WaitGroup
	for _, file := range sourceFiles {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			_ = ast.ExtractSymbols(f, symbolTable)
		}(file)

		wg.Wait()
	}

	// 4. Build Global Symbol-Level Call Graph
	for _, file := range sourceFiles {
		_ = ast.BuildSymbolCallGraph(file, symbolTable, symbolG)
	}

	return &AnalysisResult{
		FileGraph:   fileG,
		SymbolGraph: symbolG,
		SymbolTable: symbolTable,
		FileCount: len(sourceFiles),
	}, nil
}