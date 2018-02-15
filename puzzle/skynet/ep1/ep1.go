package main

import (
	"fmt"
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

func (g *Graph) AddLink(n1, n2 int) {
	g.LinksOf[n1] = append(g.LinksOf[n1], n2)
	g.LinksOf[n2] = append(g.LinksOf[n2], n1)
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
		paths := g.ComputePathsFrom(SI)
		bestExit, bestNode, min := 0, 0, math.MaxInt32
		for _, exit := range g.Exits {
			n, cost := g.LinkToCut(paths, exit)
			fmt.Fprintf(os.Stderr, "link %d %d cost %d\n", exit, n, cost)
			if cost < min {
				min = cost
				bestExit = exit
				bestNode = n
			}
		}

		// Example: 0 1 are the indices of the nodes you wish to sever the link between
		fmt.Printf("%d %d\n", bestNode, bestExit)
	}
}
