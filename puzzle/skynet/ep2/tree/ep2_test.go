package main

import (
	"strings"
	"testing"
)

var example03 = `12 20 2
0 9
6 1
0 6
0 3
0 2
11 5
10 11
11 9
10 9
8 9
5 9
4 5
0 4
0 1
3 4
8 10
0 5
1 2
6 7
2 3
0
7
8`

var example05 = `37 81 4
2 5
14 13
16 13
19 21
13 7
16 8
35 5
2 35
10 0
8 3
23 16
0 1
31 17
19 22
12 11
1 2
1 4
14 9
17 16
30 29
32 22
28 26
24 23
20 19
15 13
18 17
6 1
29 28
15 14
9 13
32 18
25 26
1 7
34 35
33 34
27 16
27 26
23 25
33 3
16 30
25 24
3 2
5 4
31 32
27 25
19 3
17 8
4 2
32 17
10 11
29 27
30 27
6 4
24 15
9 10
34 2
9 7
11 6
33 2
14 10
12 6
0 6
19 17
20 3
21 20
21 32
15 16
0 9
23 27
11 0
28 27
22 18
3 1
23 15
18 19
7 0
19 8
21 22
7 36
13 36
8 36
0
16
18
26
2`

func BenchmarkSearch(b *testing.B) {
	examples := []struct {
		name  string
		input string
	}{
		{"Example03", example03},
		{"Example05", example05},
	}

	for _, example := range examples {
		r := strings.NewReader(example.input)
		g := readGraph(r, false)
		agt := readAgent(r)
		visited := make([]bool, g.Nodes)

		var sol Sol
		var fail int
		b.Run(example.name, func(b1 *testing.B) {
			for n := 0; n < b1.N; n++ {
				sol, fail = search(g, agt, visited, 0)
			}
		})
	}
}

func TestSearchOnExamples03(t *testing.T) {
	r := strings.NewReader(example03)
	g := readGraph(r, false)
	agt := readAgent(r)
	sol, fail := search(g, agt, make([]bool, g.Nodes), 0)
	if fail != 0 {
		t.Fatalf("expected success, got a failure (%d %d %d)", sol.Exit, sol.Node, fail)
	}
}

// func TestSearchOnExamples05(t *testing.T) {
// 	r := strings.NewReader(example05)
// 	g := readGraph(r, false)
// 	agt := readAgent(r)
// 	sol, fail := search(g, agt, make([]bool, g.Nodes), 0)
// 	if fail != 0 {
// 		t.Fatalf("expected success, got a failure (%d %d %d)", sol.Exit, sol.Node, fail)
// 	}
// }
