package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"
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

func (p *Paths) ExitCount(from, to int, g *Graph) int {
	min := math.MaxInt32
	next := -1
	cnt := 0
	for _, e := range g.LinksOf[to] {
		for _, exit := range g.Exits {
			if e == exit {
				cnt++
			}
		}
		if p.Get(e) < min {
			min = p.Get(e)
			next = e
		}
	}
	// fmt.Fprintf(os.Stderr, "%d (%d) -> ", to, cnt)
	if to == from {
		return cnt
	}
	return cnt + p.ExitCount(from, next, g)
}

func dequeue(l []int) (int, []int) {
	n := l[0]
	return n, l[1:]
}

func (g *Graph) ComputePathsFrom(agt int, ignoreExit bool) *Paths {
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
			if ignoreExit {
				isExit := false
				for _, exit := range g.Exits {
					if exit == n2 {
						isExit = true
					}
				}
				if isExit {
					continue
				}
			}
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

func search(g *Graph, agt int) (int, int) {

	// agt is next to an exit cut the link
	paths := g.ComputePathsFrom(agt, false)
	for _, e := range g.Exits {
		if paths.Get(e) == 1 {
			return e, agt
		}
	}

	// sort exits by link count (the more links, the exit will win a tie)
	sort.Slice(g.Exits, func(i, j int) bool {
		return len(g.LinksOf[i]) >= len(g.LinksOf[j])
	})

	// count number of linked exit for other nodes
	exitNeighbourCount := make([]int, g.Nodes)
	for _, e := range g.Exits {
		for _, n := range g.LinksOf[e] {
			exitNeighbourCount[n]++
		}
	}

	// select the neightboor with max linked exit (where exit have most link)
	mostConnecteNeighbors := make([]int, 0, g.Nodes)
	max := 0
	for _, e := range g.Exits {
		for _, n := range g.LinksOf[e] {
			count := exitNeighbourCount[n]
			if count > max {
				max = count
				mostConnecteNeighbors = mostConnecteNeighbors[:0]
				mostConnecteNeighbors = append(mostConnecteNeighbors, n)
			} else if count == max {
				mostConnecteNeighbors = append(mostConnecteNeighbors, n)
			}

		}
	}

	// if multiple, select the node with less turn to cut (the more urgent)
	// on path from agt to a node T: the distance from agt give N turns until agt arrival to T
	// minus the number of exit on the path (you have to cut them -> -1 occasion to cut the node T)
	// ignore exit when computing distance -> agt does not pass by exits or you lose
	paths = g.ComputePathsFrom(agt, true)
	exitOnPath := make([]int, g.Nodes)
	for _, n := range mostConnecteNeighbors {
		exitOnPath[n] = paths.ExitCount(agt, n, g)
		// fmt.Fprintf(os.Stderr, "cost of %d -> %d - %d =  %d\n", n, paths.Get(n), paths.Get(n)-exitOnPath[n], exitOnPath[n])
	}
	// sort neighbor with max links by turn remaining, closest to agt
	sort.Slice(mostConnecteNeighbors, func(i, j int) bool {
		mi, mj := mostConnecteNeighbors[i], mostConnecteNeighbors[j]
		ci, cj := paths.Get(mi)-exitOnPath[mi], paths.Get(mj)-exitOnPath[mj]
		if ci < cj {
			return true
		} else if ci == cj {
			return paths.Get(mi) <= paths.Get(mj)
		}
		return false
	})

	for _, e := range g.LinksOf[mostConnecteNeighbors[0]] {
		for _, exit := range g.Exits {
			if e == exit {
				return e, mostConnecteNeighbors[0]
			}
		}
	}
	panic(fmt.Errorf("%d is not link to an exit", mostConnecteNeighbors[0]))
}

func readGraph(r io.Reader, debug bool) *Graph {
	// N: the total number of nodes in the level, including the gateways
	// L: the number of links
	// E: the number of exit gateways
	var N, L, E int
	fmt.Fscan(r, &N, &L, &E)

	if debug {
		fmt.Fprintln(os.Stderr, N, L, E)
	}

	g := New(N, E)

	for i := 0; i < L; i++ {
		// N1: N1 and N2 defines a link between these nodes
		var N1, N2 int
		fmt.Fscan(r, &N1, &N2)
		if debug {
			fmt.Fprintln(os.Stderr, N1, N2)
		}
		g.AddLink(N1, N2)
	}
	for i := 0; i < E; i++ {
		fmt.Fscan(r, &g.Exits[i])
		if debug {
			fmt.Fprintln(os.Stderr, g.Exits[i])
		}
	}
	return g
}

func readAgent(r io.Reader) int {
	var SI int
	fmt.Fscan(r, &SI)
	return SI
}

func main() {
	g := readGraph(os.Stdin, true)
	for {
		// SI: The index of the node on which the Skynet agent is positioned this turn
		SI := readAgent(os.Stdin)
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
