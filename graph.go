package graph

import (
	"sync"
	"fmt"
)

type Graph[T any] struct {
	nodes    map[string]*node[T]

	directed bool
	idTicker int
	lock     sync.RWMutex
}

type node[T any] struct {
	Value T

	id string
	edges map[string]*edge[T]
	reversedEdges map[string]*edge[T]
}

type edge[T any] struct {
	end *node[T]
	weight int
}

func New[T any](directed bool) *Graph[T] {
	return &Graph[T]{
		nodes:    make(map[string]*node[T]),

		directed: directed,
		idTicker: 0,
		lock:     sync.RWMutex{},
	}
}


func (g *Graph[T]) Node(id string) (T, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	var value T
	if node, exists := g.nodes[id]; !exists {
		return value, fmt.Errorf("node does not exist in graph: %s", id)
	} else {
		return node.Value, nil
	}
}

func (g *Graph[T]) NodeIDs() ([]string) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	ids := make([]string, 0)
	for _, n := range g.nodes {
		ids = append(ids, n.id)
	}
	return ids
}

func (g *Graph[T]) AddNode(nodeData T) string {
	g.lock.Lock()
	id := string(g.idTicker)
	g.idTicker++
	g.nodes[id] = &node[T]{
		Value: nodeData,
		id: id,
		edges: make(map[string]*edge[T]),
		reversedEdges: make(map[string]*edge[T]),
	}
	g.lock.Unlock()
	return id
}

func (g *Graph[T]) RemoveNode(id string) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	_, exists := g.nodes[id]
	if !exists {
		fmt.Errorf("node does not exist in graph: %s", id)
	}
	for _, endNode := range g.nodes {
		delete(endNode.edges, id)
		delete(endNode.reversedEdges, id)
	}
	delete(g.nodes, id)
	return nil
}

func (g *Graph[T]) AddEdge(startNodeID, endNodeID string) error {
	return g.AddWeightedEdge(startNodeID, endNodeID, 1)
}

func (g *Graph[T]) AddWeightedEdge(startNodeID, endNodeID string, weight int) error {
	if startNodeID == endNodeID {
		return fmt.Errorf("edge can not have same start and end")
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	startNode, exists := g.nodes[startNodeID]
	if !exists {
		return fmt.Errorf("start node does not exist")
	}
	endNode, exists := g.nodes[endNodeID]
	if !exists {
		return fmt.Errorf("start node does not exist")
	}
	startNode.edges[endNodeID] = &edge[T]{
		end: endNode,
		weight: weight,
	}
	endNode.reversedEdges[startNodeID] = &edge[T]{
		end: startNode,
		weight: weight,
	}
	return nil
}

func (g *Graph[T]) RemoveEdge(startNodeID, endNodeID string) error {
	if startNodeID == endNodeID {
		return fmt.Errorf("edge can not have same start and end")
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	startNode, exists := g.nodes[startNodeID]
	if !exists {
		return fmt.Errorf("start node does not exist")
	}
	endNode, exists := g.nodes[endNodeID]
	if !exists {
		return fmt.Errorf("start node does not exist")
	}
	delete(startNode.edges, endNodeID)
	delete(endNode.edges, startNodeID)
	return nil
}