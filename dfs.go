package graph

import (
	"fmt"
)

func (g Graph[T]) DFS(startNodeID string) ([]T, error) {
	visitedNodes := make(map[string]bool)
	result := make([]T, 0)
	g.lock.RLock()
	defer g.lock.RUnlock()
	startNode, exists := g.nodes[startNodeID]
	if !exists {
		return result, fmt.Errorf("start node does not exist in graph")
	}
	result = g.dfsRecursive(startNode, visitedNodes)
	return result, nil
}

func (g Graph[T]) dfsRecursive(n *node[T], visited map[string]bool) []T {
	visited[n.id] = true
	result := make([]T, 0)
	result = append(result, n.Value)
	edges := n.edges
	if !g.directed {
		for k, v := range n.reversedEdges {
			edges[k] = v
		}
	}
	for endNodeID, edge := range edges {
		if !visited[endNodeID] {
			visited[endNodeID] = true
			result = append(result, g.dfsRecursive(edge.end, visited)...)
		}
	}
	return result

}