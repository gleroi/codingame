package main

import (
	"container/list"
	"fmt"
	"io"
	"os"
	"sort"
)

//import "os"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func debug(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func main() {
	var N int
	fmt.Scan(&N)

	fmt.Println(compute2(os.Stdin, N))
}

/*
	N insert +
	n*log(n) sort +
	N comparison/subtraction +
*/

func compute2(r io.Reader, N int) int {
	horses := make([]int, N)
	for n := 0; n < N; n++ {
		fmt.Fscan(r, &horses[n])
	}

	sort.Ints(horses)

	closest := 10000001
	for i := 1; i < N; i++ {
		closest = min(closest, horses[i]-horses[i-1])
	}
	return closest
}

/*
	N insertion * N comparison + 2comparison
*/
func compute(r io.Reader, N int) int {
	horses := list.New()
	closest := 10000001
	for n := 0; n < N; n++ {
		var h int
		fmt.Fscan(r, &h)
		var elt *list.Element
		for e := horses.Front(); e != nil; e = e.Next() {
			if h < e.Value.(int) {
				elt = e
				break
			}
		}
		if elt == nil {
			elt = horses.PushBack(h)
		} else {
			elt = horses.InsertBefore(h, elt)
		}
		closest = findClosest(closest, elt)
	}
	return closest
}

func findClosest(closest int, elt *list.Element) int {
	if prev := elt.Prev(); prev != nil {
		closest = min(closest, elt.Value.(int)-prev.Value.(int))
	}
	if next := elt.Next(); next != nil {
		closest = min(closest, next.Value.(int)-elt.Value.(int))
	}
	return closest
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
