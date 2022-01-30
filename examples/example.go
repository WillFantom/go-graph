package main

import (
  "fmt"

  graph "github.com/willfantom/go-graph"
)

type MyStruct struct {
  Data string
}

func (ms MyStruct) Print() {
  fmt.Println(ms.Data)
}

func main() {
  myGraph := graph.New[MyStruct](true)

  nodeIDA := myGraph.AddNode(MyStruct{
    Data: "Hello, I am Node A",
  })
  nodeIDB := myGraph.AddNode(MyStruct{
    Data: "Hello, I am Node B",
  })
	myGraph.AddEdge(nodeIDA, nodeIDB)

  if bfsResults, err := myGraph.BFS(nodeIDA); err != nil {
    panic(err)
  } else {
    for _, node := range bfsResults {
			node.Print()
		}
	}
}
