package graph

import (
	"fmt"
)

type fifo[T any] struct {
	queue []*node[T]
}

func newFIFO[T any]() *fifo[T]{
	return &fifo[T]{
		queue: make([]*node[T], 0),
	}
} 

func (q *fifo[T]) pushToBase(n *node[T]) {
	q.queue = append(q.queue, n)
}

func (q *fifo[T]) popFromTop() (*node[T], error) {
	var n *node[T]
	if len(q.queue) == 0 {
		return n, fmt.Errorf("fifo queue is empty")
	}
	n = q.queue[0]
	q.queue[0] = nil
	q.queue = q.queue[1:]
	return n, nil
}

func (q *fifo[T]) empty() bool {
	if len(q.queue) == 0 {
		return true
	}
	return false
} 

func (g Graph[T]) BFS(startNodeID string) ([]T, error) {
	queue := newFIFO[T]()
	visitedNodes := make(map[string]bool)
	result := make([]T, 0)
	g.lock.RLock()
	defer g.lock.RUnlock()
	startNode, exists := g.nodes[startNodeID]
	if !exists {
		return result, fmt.Errorf("start node does not exist in graph")
	}
	queue.pushToBase(startNode)
	for !queue.empty() {
		currNode, err := queue.popFromTop()
		if err != nil {
			return result, fmt.Errorf("attempted to pop from empty fifo: should not reach here")
		}
		visitedNodes[currNode.id] = true
		result = append(result, currNode.Value)
		edges := currNode.edges
		if !g.directed {
			for k, v := range currNode.reversedEdges {
				edges[k] = v
			}
		}
		for endNodeID, edge := range edges {
			if !visitedNodes[endNodeID] {
				queue.pushToBase(edge.end)
				visitedNodes[endNodeID] = true
			}
		}
	}
	return result, nil
}

func (g Graph[T]) Order() ([]T, error) { 
	result := make([]T, 0)
	if !g.directed {
		return result, fmt.Errorf("undirected graph can not be ordered")
	}
	indegrees := make(map[string]int)
	queue := newFIFO[T]()
	for id, n := range g.nodes {
		indegrees[id] = len(n.reversedEdges)
		if indegrees[id] == 0 {
			queue.pushToBase(n)
		}
	}
	for !queue.empty() {
		currNode, err := queue.popFromTop()
		if err != nil {
			return result, fmt.Errorf("attempted to pop from empty fifo: should not reach here")
		}
		result = append(result, currNode.Value)
		delete(indegrees, currNode.id)
		for id, n := range g.nodes {
			if _, exists := indegrees[id]; exists {
				indegrees[id]--
				if indegrees[id] == 0 {
					queue.pushToBase(n)
				} 
			}
		}
	}
	return result, nil
}