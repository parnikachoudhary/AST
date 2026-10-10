package ast

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"ast-engine/pkg/graph"
)


func BuildSymbolCallGraph(filePath string, table *SymbolTable, symbolGraph *graph.Graph) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sourceFile := filepath.Base(filePath)
	scanner := bufio.NewScanner(f)

	var currentCaller string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		
		for _, sym := range table.Symbols {
			if sym.SourceFile == sourceFile && sym.Type == SymbolFunction {
			
				funcDefPattern := regexp.MustCompile(`\b` + sym.Name + `\s*\(`)
				if funcDefPattern.MatchString(line) && strings.Contains(line, "{") {
					currentCaller = sym.Name
					break
				}
			}
		}

		
		if currentCaller != "" {
			for _, sym := range table.Symbols {
			
				if sym.Name != currentCaller {
					callPattern := regexp.MustCompile(`\b` + sym.Name + `\s*\(`)
					if callPattern.MatchString(line) {
						symbolGraph.AddEdge(currentCaller, sym.Name)
					}
				}
			}

			
			if line == "}" || line == "};" {
				currentCaller = ""
			}
		}
	}

	return scanner.Err()
}