package main

// FalsePeak generates a false peak problem
func FalsePeak(n int) (int, *Node) {
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
	b.Nodes[0] = &Node{
		NodeType: Var,
		A:        0,
	}
	for i := 1; i < len(b.Nodes); i++ {
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
	return n, c
}
