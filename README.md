# Graph

> 🚧 Work-In-Progress

A (very) simple graph package that utilizes the Generics features in Go 1.18

---

## Usage

```go
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
```

> 🎉 No need to cast an `interface{}` to the desired node type anymore!

---

## Generics

Typically, a graph package will ensure that graph nodes have some form of open type in their structs, for example:
```go
type Node struct {
  Value interface{}
}
```

However, this requires some type casting during runtime to get back the more specific struct type you want as nodes... (this can easily cause errors)... Using the generics features added in Go 1.18, this package can catch those errors at build time, and thus it has a node structure like:
```go
type Node[T any] struct {
  Value T
  
  ...
}
```