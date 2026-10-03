package graph

// Edge represents a directed connection between two nodes (Functions or Files)
type Edge struct {
	From string
	To   string
}

// Graph holds the adjacency list representation
type Graph struct {
	AdjacencyList map[string][]string
}

// NewGraph initializes an empty Graph structure with allocated memory
func NewGraph() *Graph {
	return &Graph{
		AdjacencyList: make(map[string][]string),
	}
}

// AddEdge inserts a directed relation into the Adjacency List
func (g *Graph) AddEdge(from, to string) {
	g.AdjacencyList[from] = append(g.AdjacencyList[from], to)
}