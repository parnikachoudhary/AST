package graph

import (
	"fmt"
	"os"
)

// ExportToDOT serializes the graph into Graphviz DOT syntax and writes to disk.
func (g *Graph) ExportToDOT(filename string) error {
	// Create or overwrite the target file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	// defer ensures the file closes when this function finishes executing
	defer file.Close()

	file.WriteString("digraph CallGraph {\n")
	file.WriteString("  node [shape=box, fontname=\"Helvetica\", style=filled, fillcolor=\"#f4f4f4\"];\n")

	for u, neighbors := range g.AdjacencyList {
		for _, v := range neighbors {
			line := fmt.Sprintf("  \"%s\" -> \"%s\";\n", u, v)
			file.WriteString(line)
		}
	}

	file.WriteString("}\n")
	return nil
}