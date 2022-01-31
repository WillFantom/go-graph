package graph

import (
	"sync"
	"fmt"
)

type Graph[T any] struct {
	nodes    map[NodeID]*node[T]

	directed bool
	idTicker int
	lock     sync.RWMutex
}

type NodeID int

type node[T any] struct {
	Value T

	id NodeID
	edges map[NodeID]*edge[T]
	reversedEdges map[NodeID]*edge[T]
}

type edge[T any] struct {
	end *node[T]
	weight int
}

func New[T any](directed bool) *Graph[T] {
	return &Graph[T]{
		nodes:    make(map[NodeID]*node[T]),

		directed: directed,
		idTicker: 0,
		lock:     sync.RWMutex{},
	}
}


func (g *Graph[T]) Node(id NodeID) (T, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	var value T
	if node, exists := g.nodes[id]; !exists {
		return value, fmt.Errorf("node does not exist in graph: %s", id)
	} else {
		return node.Value, nil
	}
}

func (g *Graph[T]) NodeIDs() ([]NodeID) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	ids := make([]NodeID, 0)
	for _, n := range g.nodes {
		ids = append(ids, n.id)
	}
	return ids
}

func (g *Graph[T]) AddNode(nodeData T) NodeID {
	g.lock.Lock()
	id := NodeID(g.idTicker)
	g.idTicker++
	g.nodes[id] = &node[T]{
		Value: nodeData,
		id: id,
		edges: make(map[NodeID]*edge[T]),
		reversedEdges: make(map[NodeID]*edge[T]),
	}
	g.lock.Unlock()
	return id
}

func (g *Graph[T]) RemoveNode(id NodeID) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	_, exists := g.nodes[id]
	if !exists {
		fmt.Errorf("node does not exist in graph: %v", id)
	}
	for _, endNode := range g.nodes {
		delete(endNode.edges, id)
		delete(endNode.reversedEdges, id)
	}
	delete(g.nodes, id)
	return nil
}

func (g *Graph[T]) AddEdge(startNodeID, endNodeID NodeID) error {
	return g.AddWeightedEdge(startNodeID, endNodeID, 1)
}

func (g *Graph[T]) AddWeightedEdge(startNodeID, endNodeID NodeID, weight int) error {
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

func (g *Graph[T]) RemoveEdge(startNodeID, endNodeID NodeID) error {
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