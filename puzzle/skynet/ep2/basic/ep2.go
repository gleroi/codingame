package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

//import "os"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

type Graph struct {
	Nodes   int
	LinksOf [][]int
	Exits   []int
}

func New(N, E int) *Graph {
	g := &Graph{
		Nodes:   N,
		LinksOf: make([][]int, N),
		Exits:   make([]int, E),
	}
	for i := 0; i < g.Nodes; i++ {
		g.LinksOf[i] = make([]int, 0, N)
	}
	return g
}

func (g *Graph) Clone() *Graph {
	c := New(g.Nodes, len(g.Exits))
	copy(c.Exits, g.Exits)

	for n1 := range g.LinksOf {
		c.LinksOf[n1] = append(c.LinksOf[n1], g.LinksOf[n1]...)
	}
	return c
}

func (g *Graph) Print(w io.Writer, paths *Paths) {
	fmt.Fprintln(w, "strict graph {")
	for n1, links := range g.LinksOf {
		for _, n2 := range links {
			fmt.Fprintln(w, n1, "--", n2, "[ label = \"", paths.Get(n2), "\"]")
		}
	}
	fmt.Fprintln(w, "}")
}

func (g *Graph) AddLink(n1, n2 int) {
	g.LinksOf[n1] = append(g.LinksOf[n1], n2)
	g.LinksOf[n2] = append(g.LinksOf[n2], n1)
}

func indexOf(s []int, x int) int {
	for i := range s {
		if s[i] == x {
			return i
		}
	}
	return -1
}

func remove(s []int, i int) []int {
	return append(s[:i], s[i+1:]...)
}

func (g *Graph) RemoveLink(n1, n2 int) {
	i := indexOf(g.LinksOf[n1], n2)
	g.LinksOf[n1] = remove(g.LinksOf[n1], i)

	i = indexOf(g.LinksOf[n2], n1)
	g.LinksOf[n2] = remove(g.LinksOf[n2], i)
}

type Paths struct {
	Costs []int
	N     int
}

func NewPaths(N int) *Paths {
	p := &Paths{
		Costs: make([]int, N),
	}
	for i := range p.Costs {
		p.Costs[i] = math.MaxInt32
	}
	return p
}

func (p *Paths) Set(n1 int, cost int) {
	p.Costs[n1] = cost
}

func (p *Paths) Get(n1 int) int {
	return p.Costs[n1]
}

func dequeue(l []int) (int, []int) {
	n := l[0]
	return n, l[1:]
}

func (g *Graph) ComputePathsFrom(agt int) *Paths {
	visited := make(map[int]bool, g.Nodes)
	toVisit := make([]int, 0, g.Nodes)
	toVisit = append(toVisit, agt)
	paths := NewPaths(g.Nodes)
	paths.Set(agt, 0)

	for len(toVisit) > 0 {
		var n1 int
		n1, toVisit = dequeue(toVisit)
		visited[n1] = true

		for _, n2 := range g.LinksOf[n1] {
			newCost := paths.Get(n1) + 1
			if newCost < paths.Get(n2) {
				paths.Set(n2, newCost)
				toVisit = append(toVisit, n2)
			} else if !visited[n2] {
				toVisit = append(toVisit, n2)
			}
		}
	}
	return paths
}

// LinkToCut returns the node to cut on the shortest to exit from agent and the cost
// of this path
func (g *Graph) LinkToCut(paths *Paths, exit int) (int, int) {
	cost := paths.Get(exit)
	n := exit
	min := math.MaxInt32
	for _, n1 := range g.LinksOf[exit] {
		if paths.Get(n1) < min {
			min = paths.Get(n1)
			n = n1
		}
	}
	return n, cost
}

func search(g *Graph, agt int) (int, int) {

	// if all neighbour have a cost <= to the number of exit
	// cut link to exit with most nodes
	//TODO: rework this part of the strategy, probably need a tree of all decision :(
	paths := g.ComputePathsFrom(agt)
	exitsNeighbors := make([]int, g.Nodes)
	for _, e := range g.Exits {
		for _, n2 := range g.LinksOf[e] {
			exitsNeighbors[n2] += 1
		}
	}
	good := true
	for n, exitCount := range exitsNeighbors {
		if exitCount > paths.Get(n) {
			fmt.Fprintf(os.Stderr, "invalid by node %d: cost: %d, exitCount: %d\n", n, paths.Get(n), exitCount)
			good = false
		}
	}
	max := 0
	if good {
		bestN := 0
		for n, exitCount := range exitsNeighbors {
			if exitCount > max {
				max = exitCount
				bestN = n
			}
		}
		for _, e := range g.Exits {
			if indexOf(g.LinksOf[e], bestN) != -1 {
				return e, bestN
			}
		}
	}

	// this part works, but only when one gate is at risk. If no gate is at risk no
	// solution is found.
	for _, n1 := range g.Exits {
		for _, n2 := range g.LinksOf[n1] {
			g1 := g.Clone()
			g1.RemoveLink(n1, n2)
			paths := g1.ComputePathsFrom(agt)

			exitsNeighbors := make([]int, g.Nodes)
			for _, e := range g1.Exits {
				for _, n2 := range g1.LinksOf[e] {
					exitsNeighbors[n2] += 1
				}
			}

			// good if all neightbors should have a cost <= to the number of exit
			for _, exitCount := range exitsNeighbors {
				if exitCount == 0 {
					continue
				}
			}
			good := true
			for n, exitCount := range exitsNeighbors {
				if exitCount > paths.Get(n) {
					good = false
				}
			}
			if good {
				return n1, n2
			}
		}
	}
	panic("no solution!")
}

func main() {
	// N: the total number of nodes in the level, including the gateways
	// L: the number of links
	// E: the number of exit gateways
	var N, L, E int
	fmt.Scan(&N, &L, &E)
	g := New(N, E)

	for i := 0; i < L; i++ {
		// N1: N1 and N2 defines a link between these nodes
		var N1, N2 int
		fmt.Scan(&N1, &N2)
		g.AddLink(N1, N2)
	}
	for i := 0; i < E; i++ {
		fmt.Scan(&g.Exits[i])
	}

	for {
		// SI: The index of the node on which the Skynet agent is positioned this turn
		var SI int
		fmt.Scan(&SI)

		/*
			For each exits find the shortest path to the agent
			Amongst theses paths, select the shortest and cut a link
		*/
		fmt.Fprintf(os.Stderr, "agent is in %d\n", SI)
		fmt.Fprintf(os.Stderr, "exits are %d\n", g.Exits)

		n1, n2 := search(g, SI)

		fmt.Printf("%d %d\n", n1, n2)

		g.RemoveLink(n1, n2)
	}
}
