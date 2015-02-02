package sudoku

import (
	"errors"
	"fmt"
)

type Grid struct {
	Nodes      [9][9]Node
	FixedNodes int
}

type Node struct {
	Vals  map[int]bool
	Fixed bool
}

func NewNode() Node {
	n := Node{
		Vals:  make(map[int]bool),
		Fixed: false,
	}
	for i := 1; i <= 9; i++ {
		n.Vals[i] = true
	}
	return n
}

func (n *Node) Len() int {
	l := 0
	for _, v := range n.Vals {
		if v {
			l++
		}
	}
	return l
}

func (n *Node) Val() int {
	if n.Len() > 1 {
		panic("sudoku: node has multiple values")
	}
	for k, v := range n.Vals {
		if v {
			return k
		}
	}
	panic("sudoku: node has no values")
}

func (n *Node) Fix(val int) {
	for k, _ := range n.Vals {
		n.Vals[k] = false
	}
	n.Vals[val] = true
	n.Fixed = true
}

func (n *Node) Remove(val int) bool {
	if n.Fixed {
		return false
	}
	if n.Vals[val] {
		n.Vals[val] = false
		return true
	}
	return false
}

func (n *Node) Copy() Node {
	nn := NewNode()
	nn.Fixed = n.Fixed
	for k, v := range n.Vals {
		nn.Vals[k] = v
	}
	return nn
}

func NewGrid() Grid {
	g := Grid{}
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			g.Nodes[x][y] = NewNode()
		}
	}
	return g
}

func (g *Grid) Set(s string) (bool, error) {
	if len(s) != 9*9 {
		return false, errors.New("sudoku: invalid length")
	}
	for i := 0; i < 9*9; i++ {
		c := s[i]
		switch {
		case c == '0' || c == ' ' || c == '_':
		case c >= '1' && c <= '9':
			if !g.Fix(i%9, i/9, int(byte(c)-'1')+1) {
				return false, errors.New("sudoku: invalid grid")
			}
		default:
			return false, fmt.Errorf("sudoku: invalid char %c", c)
		}

	}
	return true, nil
}

func (g *Grid) String() string {
	var s string
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			n := g.NodeAt(x, y)
			if n.Fixed {
				s += fmt.Sprintf("%d", n.Val())
			} else {
				s += "_"
			}
		}
	}
	return s
}

func (g *Grid) NodeAt(x, y int) *Node {
	return &g.Nodes[x][y]
}

func (g *Grid) Copy() Grid {
	ng := NewGrid()
	ng.FixedNodes = g.FixedNodes
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			ng.Nodes[x][y] = g.Nodes[x][y].Copy()
		}
	}
	return ng
}

func (g *Grid) Complete() bool {
	return g.FixedNodes == 9*9
}

func (g *Grid) Print() {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if x > 0 && x%3 == 0 {
				fmt.Printf(" ")
			}
			n := g.NodeAt(x, y)
			if n.Fixed {
				fmt.Printf("%d", n.Val())
			} else {
				fmt.Printf("_")
			}
			if x == 8 {
				fmt.Printf("\n")
			}
		}
		if y != 8 && (y+1)%3 == 0 {
			fmt.Println()
		}
	}
}

func (g *Grid) Fix(x, y, val int) bool {
	return g.FixNode(g.NodeAt(x, y), val)
}

func (g *Grid) FixNode(n *Node, val int) bool {
	if !n.Vals[val] {
		return false
	}
	n.Fix(val)
	for {
		red, inval := g.Reduce()
		if inval {
			return false
		}
		if !red {
			break
		}
	}
	g.FixedNodes++
	return true
}

func (g *Grid) FixNext() (fixed, inval bool) {
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			n := g.NodeAt(x, y)
			if n.Fixed || n.Len() != 1 {
				continue
			}
			return true, !g.FixNode(n, n.Val())
		}
	}
	return false, false
}

func (g *Grid) Reduce() (reduced, invalid bool) {
	reduced = false

	for x := 0; x < 9; x++ {
		red, inval := g.ReduceCol(x)
		if inval {
			return true, true
		}
		if red {
			reduced = true
		}
	}

	for y := 0; y < 9; y++ {
		red, inval := g.ReduceRow(y)
		if inval {
			return true, true
		}
		if red {
			reduced = true
		}
	}

	for x := 0; x < 9; x += 3 {
		for y := 0; y < 9; y += 3 {
			red, inval := g.ReduceBox(x, y)
			if inval {
				return true, true
			}
			if red {
				reduced = true
			}
		}
	}

	return reduced, false
}

func (g *Grid) ReduceCol(x int) (reduced, invalid bool) {
	return g.reduceNodes(g.NodesForCol(x))
}

func (g *Grid) NodesForCol(x int) []*Node {
	var nodes []*Node
	for y := 0; y < 9; y++ {
		n := g.NodeAt(x, y)
		nodes = append(nodes, n)
	}
	return nodes
}

func (g *Grid) ReduceRow(y int) (reduced, invalid bool) {
	return g.reduceNodes(g.NodesForRow(y))
}

func (g *Grid) NodesForRow(y int) []*Node {
	var nodes []*Node
	for x := 0; x < 9; x++ {
		n := g.NodeAt(x, y)
		nodes = append(nodes, n)
	}
	return nodes
}

func (g *Grid) reduceNodes(ns []*Node) (reduced, invalid bool) {
	reduced = false
	for i := 0; i < 9; i++ {
		n := ns[i]
		if n.Len() != 1 {
			continue
		}
		v := n.Val()
		for j := 0; j < 9; j++ {
			if j == i {
				continue
			}
			nn := ns[j]
			if nn.Remove(v) {
				if nn.Len() == 0 {
					return true, true
				}
				reduced = true
			}
		}
	}
	return reduced, false
}

func (g *Grid) ReduceBox(x, y int) (reduced, invalid bool) {
	return g.reduceNodes(g.NodesForBox(x, y))
}

func (g *Grid) NodesForBox(x, y int) []*Node {
	x = 3 * (x / 3)
	y = 3 * (y / 3)
	var nodes []*Node
	for yy := y; yy < y+3; yy++ {
		for xx := x; xx < x+3; xx++ {
			n := g.NodeAt(xx, yy)
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (g *Grid) Solve() bool {
	var x, y int
	var n *Node

search:
	for x = 0; x < 9; x++ {
		for y = 0; y < 9; y++ {
			n = g.NodeAt(x, y)
			if n.Fixed {
				continue
			}
			break search
		}
	}
	if n == nil {
		return false
	}

vals:
	for i := 1; i <= 9; i++ {
		if !n.Vals[i] {
			continue
		}
		ng := g.Copy()
		if !ng.Fix(x, y, i) {
			continue
		}
		for {
			fixed, inval := ng.FixNext()
			if inval {
				continue vals
			}
			if !fixed {
				break
			}
		}
		if ng.Complete() || ng.Solve() {
			g.Nodes = ng.Nodes
			return true
		}
	}
	return false
}
