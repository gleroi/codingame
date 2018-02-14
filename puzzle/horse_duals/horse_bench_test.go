package main

import (
	"strings"
	"testing"
)

var result int

func BenchmarkSortSliceInt(b *testing.B) {
	var d int
	r := strings.NewReader(examples[3].input)
	for n := 0; n < b.N; n++ {
		r.Reset(examples[3].input)
		d = compute2(r, examples[3].length)
	}
	result = d
}

func BenchmarkInsertionLinkedList(b *testing.B) {
	var d int
	r := strings.NewReader(examples[3].input)
	for n := 0; n < b.N; n++ {
		r.Reset(examples[3].input)
		d = compute(r, examples[3].length)
	}
	result = d
}
