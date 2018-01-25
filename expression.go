package main

import "fmt"

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

// String generates a string for the expression
func (n *Node) String() string {
	var str func(n *Node) string
	str = func(n *Node) string {
		switch n.NodeType {
		case NoOp:
			panic("noop")
		case And:
			s := "("
			and := ""
			for _, node := range n.Nodes {
				s += and + str(node)
				and = " & "
			}
			return s + ")"
		case Or:
			s := "("
			or := ""
			for _, node := range n.Nodes {
				s += or + str(node)
				or = " | "
			}
			return s + ")"
		case Not:
			return "!" + str(n.Nodes[0])
		case Var:
			return fmt.Sprintf("%d", n.A)
		default:
			panic("invalid node class")
		}
	}

	return str(n)
}
