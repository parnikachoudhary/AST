package ast

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"fmt"
)

var (
	
	classDeclRegex = regexp.MustCompile(`\bclass\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*([:\{]|$)`)
	
	structDeclRegex = regexp.MustCompile(`\bstruct\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*([:\{]|$)`)
	
	funcDeclRegex = regexp.MustCompile(`\b(?:[a-zA-Z_][a-zA-Z0-9_]*::)?([a-zA-Z_][a-zA-Z0-9_]*)\s*\([^;]*\)\s*(?:const\s*)?\{`)
)

		func ResolveScopedName(scope string, symbolName string) string {
			if scope != "" {
				return fmt.Sprintf("%s::%s", scope, symbolName)
			}
			return symbolName
		}

func ExtractSymbols(filePath string, table *SymbolTable) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sourceFile := filepath.Base(filePath)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") {
			continue
		}

		

		// 1. Detect Class Declarations
		if matches := classDeclRegex.FindStringSubmatch(line); len(matches) > 1 {
			className := matches[1]

			if className != "explicit" && className != "public" && className != "private" {
				
				table.AddSymbol(Symbol{
					Name:       className,
					Type:       SymbolClass,
					SourceFile: sourceFile,
					LineNumber: lineNum,
				})
				continue
			}
		}

		// 2. Detect Struct Declarations
		if matches := structDeclRegex.FindStringSubmatch(line); len(matches) > 1 {
			table.AddSymbol(Symbol{
				Name:       matches[1],
				Type:       SymbolStruct,
				SourceFile: sourceFile,
				LineNumber: lineNum,
			})
			continue
		}

		// 3. Detect Real Function Definitions (Ending with '{')
		if matches := funcDeclRegex.FindStringSubmatch(line); len(matches) > 1 {
			funcName := matches[1]
			
	
			ignored := map[string]bool{
				"if": true, "for": true, "while": true, "switch": true,
				"catch": true, "erase": true, "find_if": true, "base": true,
				"length": true, "isspace": true, "explicit": true,
			}

			if !ignored[funcName] {
			
				table.AddSymbol(Symbol{
					
					Name:       funcName,
					Type:       SymbolFunction,
					SourceFile: sourceFile,
					LineNumber: lineNum,
				})
			}
		}
	}

	

	return scanner.Err()
}

