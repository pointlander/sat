package main

import "fmt"

// GraphNode represensts a node in a directed graph
type GraphNode struct {
	In, Out []int
}

// Graph represents a graph made up of nodes
type Graph []GraphNode

// FindIn finds the in values for the graph nodes
func (g Graph) FindIn() {
	for i := range g {
		for _, out := range g[i].Out {
			g[out].In = append(g[out].In, i)
		}
	}
}

// HamiltonianCircuit creates a hamiltonian circuit problem of a given size
func HamiltonianCircuit(n int) (int, *Node) {
	graph := make(Graph, n)

	out, c := make([]int, n-2), 1
	for i := range out {
		out[i] = c
		c++
	}
	graph[0].Out = out

	for i := 1; i < n-1; i++ {
		c = i + 1
		out = make([]int, n-c)
		for j := range out {
			out[j] = c
			c++
		}
		graph[i].Out = out
	}

	graph[n-1].Out = make([]int, 1)

	graph.FindIn()

	vars, v := make(map[string]int), 0
	getVar := func(a, b int) int {
		key := fmt.Sprintf("%d_%d", a, b)
		value, ok := vars[key]
		if !ok {
			value = v
			v++
			vars[key] = value
		}
		return value
	}

	a := &Node{
		NodeType: And,
	}
	for i := range graph {
		// in clauses
		if length := len(graph[i].In); length > 1 {
			b := &Node{
				NodeType: Or,
				Nodes:    make([]*Node, length),
			}
			for j := range b.Nodes {
				b.Nodes[j] = &Node{
					NodeType: Var,
					A:        getVar(i, graph[i].In[j]),
				}
			}
			a.Nodes = append(a.Nodes, b)
		} else {
			b := &Node{
				NodeType: Var,
				A:        getVar(i, graph[i].In[0]),
			}
			a.Nodes = append(a.Nodes, b)
		}

		// out clauses
		if length := len(graph[i].Out); length > 1 {
			b := &Node{
				NodeType: Or,
				Nodes:    make([]*Node, length),
			}
			for j := range b.Nodes {
				c := &Node{
					NodeType: And,
					Nodes:    make([]*Node, length),
				}
				for k := range c.Nodes {
					if k == j {
						c.Nodes[k] = &Node{
							NodeType: Var,
							A:        getVar(graph[i].Out[k], i),
						}
					} else {
						d := &Node{
							NodeType: Not,
							Nodes:    make([]*Node, 1),
						}
						d.Nodes[0] = &Node{
							NodeType: Var,
							A:        getVar(graph[i].Out[k], i),
						}
						c.Nodes[k] = d
					}
				}
				b.Nodes[j] = c
			}
			a.Nodes = append(a.Nodes, b)
		} else {
			b := &Node{
				NodeType: Var,
				A:        getVar(graph[i].Out[0], i),
			}
			a.Nodes = append(a.Nodes, b)
		}
	}

	return len(vars), a
}
