package main

import (
	"fmt"

	"github.com/pointlander/go-galib"
)

// NodeType is the type of node
type NodeType int

const (
	// NoOp is a non-operation
	NoOp NodeType = iota
	// And is an and function
	And
	// Or is an or function
	Or
	// Not is a not function
	Not
	// Var is a variable
	Var
)

const (
	// Decimals is the number of fixed point decimals
	Decimals = 30
	// One is 1 in fixed point
	One = 1 << Decimals
	// ProblemSize is the size of the problems
	ProblemSize = 1000
)

// Node is represensts part of a program
type Node struct {
	NodeType
	A     int
	Nodes []*Node
}

// Eval evaluates an expression
func (n *Node) Eval(vars []int32) int32 {
	var eval func(n *Node) int32
	eval = func(n *Node) int32 {
		switch n.NodeType {
		case NoOp:
			panic("noop")
		case And:
			var sum int64
			for _, node := range n.Nodes {
				sum += int64(eval(node))
			}
			return int32(sum / int64(len(n.Nodes)))
		case Or:
			var max int32
			for _, node := range n.Nodes {
				if value := eval(node); value > max {
					max = value
				}
			}
			return max
		case Not:
			return One - eval(n.Nodes[0])
		case Var:
			return vars[n.A]
		default:
			panic("invalid node class")
		}
	}

	return eval(n)
}

// TwoPeak generates a two peak problem
func TwoPeak(n int) *Node {
	a := &Node{
		NodeType: And,
		Nodes:    make([]*Node, n),
	}
	for i := range a.Nodes {
		a.Nodes[i] = &Node{
			NodeType: Var,
			A:        i,
		}
	}

	b := &Node{
		NodeType: And,
		Nodes:    make([]*Node, n),
	}
	for i := range b.Nodes {
		c := &Node{
			NodeType: Not,
			Nodes:    make([]*Node, 1),
		}
		c.Nodes[0] = &Node{
			NodeType: Var,
			A:        i,
		}
		b.Nodes[i] = c
	}

	c := &Node{
		NodeType: Or,
		Nodes:    make([]*Node, 2),
	}
	c.Nodes[0] = a
	c.Nodes[1] = b
	return c
}

var problem = TwoPeak(ProblemSize)

func score(g *ga.GAFixedBitstringGenome) float64 {
	vars := make([]int32, len(g.Gene))
	for i, bit := range g.Gene {
		if bit {
			vars[i] = One
		}
	}
	return float64(One-problem.Eval(vars)) / One
}

func main() {
	m := ga.NewMultiMutator()
	m.Add(new(ga.GAShiftMutator))
	m.Add(new(ga.GASwitchMutator))
	m.Add(new(ga.GAMutatorRandom))
	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		PMutate:     0.2,
		PBreed:      0.2,
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     m,
	}

	gao := ga.NewGA(param)
	genome := ga.NewFixedBitstringGenome(make([]bool, ProblemSize), score)
	gao.Init(128, genome)
	gao.OptimizeUntil(func(best ga.GAGenome) bool {
		score := best.Score()
		fmt.Printf("best = %v\n", score)
		return score < 1e-3
	})
}
