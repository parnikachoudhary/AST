package graph

// HasPathWithoutEdge runs a DFS to check if 'target' is reachable from 'start'

func (g *Graph) HasPathWithoutEdge(start, target, disabledFrom, disabledTo string) bool {
	// 1. visited keeps track of traversed nodes to prevent infinite loops (like set() in Python)
	visited := make(map[string]bool)

	// 2. stack acts as our LIFO queue for DFS traversal
	stack := []string{start}

	for len(stack) > 0 {
		// Pop the last element from stack
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Target reached via alternative path!
		if curr == target && curr != start {
			return true
		}

		if !visited[curr] {
			visited[curr] = true

			// Traverse neighbors of current node
			for _, neighbor := range g.AdjacencyList[curr] {
				// Block the disabled direct edge
				if curr == disabledFrom && neighbor == disabledTo {
					continue
				}
				if !visited[neighbor] {
					stack = append(stack, neighbor)
				}
			}
		}
	}

	return false
}

// TransitiveReduction creates a new reduced Graph removing redundant edges.
func (g *Graph) TransitiveReduction() *Graph {
	reducedGraph := NewGraph()

	for u, neighbors := range g.AdjacencyList {
		for _, v := range neighbors {
			// If no alternative path exists, this edge must be added.
			if !g.HasPathWithoutEdge(u, v, u, v) {
				reducedGraph.AddEdge(u, v)
			}
		}
	}

	return reducedGraph
}