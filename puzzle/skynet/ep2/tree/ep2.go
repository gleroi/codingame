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

func (p *Paths) FindClosest(from int, tos []int) int {
	closest, best := math.MaxInt32, -1
	for _, to := range tos {
		if p.Get(to) < closest {
			closest = p.Get(to)
			best = to
		}
	}
	return best
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

func (g *Graph) Move(from int, to int, paths *Paths) int {
	closest := paths.FindClosest(from, g.LinksOf[to])
	if closest == from {
		return to
	}
	return g.Move(from, closest, paths)
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

type Sol struct {
	Exit, Node int
}

func search(g *Graph, agt int, visited []bool, indent int) (Sol, int) {

	// fmt.Fprintf(os.Stderr, "agt: %d, visited: %v\n", agt, prevVisited)
	/*
			  recursive search of the solution:
			  - agent on exit -> Empty solution, 1 fail
			  - no links on any exit -> Empty solution, 0 fail (success!)
			  - exit with links
			  	- remove a link
				- for all possible agent move
		      	  - search and collect sol, fail count
			    - removing this link is a good solution if fail = 0
	*/
	reached := indexOf(g.Exits, agt)
	if reached != -1 {
		return Sol{0, 0}, 1
	}

	success := true
	for _, exit := range g.Exits {
		if len(g.LinksOf[exit]) != 0 {
			success = false
			break
		}
	}
	if success {
		return Sol{0, 0}, 0
	}

	visited[agt] = true
	defer func() { visited[agt] = false }()

	for _, exit := range g.Exits {
		for _, n1 := range g.LinksOf[exit] {
			// fmt.Fprintf(os.Stderr, "%sagt: %d %d -> %d\n", strings.Repeat(" ", indent), agt, exit, n1)

			g1 := g.Clone()
			g1.RemoveLink(exit, n1)

			sol := Sol{Exit: exit, Node: n1}
			fail := 0
			for _, newAgt := range g1.LinksOf[agt] {
				if visited[newAgt] {
					// fmt.Fprintf(os.Stderr, "%sagt: %d %d already seen\n", strings.Repeat(" ", indent), agt, newAgt)
					continue
				}
				_, subFail := search(g1, newAgt, visited, indent+2)
				fail += subFail
			}

			g1.AddLink(exit, n1)
			if fail == 0 {
				// fmt.Fprintf(os.Stderr, "%sagt: %d %d -> %d no fail\n", strings.Repeat(" ", indent), agt, exit, n1)
				return sol, 0
			}
		}
	}

	return Sol{0, 0}, 1
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

		sol, fail := search(g, SI, make([]bool, g.Nodes), 0)

		fmt.Printf("%d %d (%d)\n", sol.Exit, sol.Node, fail)

		g.RemoveLink(sol.Exit, sol.Node)
	}
}
