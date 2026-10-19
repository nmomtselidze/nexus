package graph

// Edge represents a directed edge in the graph
type Edge struct {
	From   string
	To     string
	Weight float64
}

// NewEdge creates a new edge
func NewEdge(from, to string, weight float64) *Edge {
	return &Edge{
		From:   from,
		To:     to,
		Weight: weight,
	}
}
