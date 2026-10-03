package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"ast-engine/pkg/graph"
)

// Regex to capture local includes: #include "filename.h"
var includeRegex = regexp.MustCompile(`^\s*#include\s*["<]([^">]+)[">]`)

// ScanDirectory walks through the target directory :
func ScanDirectory(rootDir string) (*graph.Graph, error) {
	g := graph.NewGraph()

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".cpp" || ext == ".c" || ext == ".h" || ext == ".hpp" {
			parseFileIncludes(path, g)

			
		}

		return nil
	})

	return g, err
}

func parseFileIncludes(filePath string, g *graph.Graph) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	sourceFile := filepath.Base(filePath)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		matches := includeRegex.FindStringSubmatch(line)

		if len(matches) > 1 {
			targetFile := matches[1]

			if !strings.HasPrefix(line, "#include <") {
				g.AddEdge(sourceFile, filepath.Base(targetFile))
			}
		}
	}

	// error check
	if err := scanner.Err(); err != nil {
		return
	}
}