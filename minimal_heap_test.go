package main

import (
	"math/rand"
	"testing"
)

type IntEva int

func (v IntEva) Value() int {
	return int(v)
}

func TestMinimalHeap(t *testing.T) {
	var heap IHeap = &Heap{cap: 4}
	heap.Init()
	numbers := []IntEva{5, 3, 6, 2, 1, 9, 2, 7}
	sorted := []IntEva{1, 2, 2, 3, 5, 6, 7, 9}
	half := len(numbers) / 2
	for _, n := range numbers {
		heap.Add(n)
	}
	for _, v := range sorted[:half] {
		if r := heap.Peek(); r != nil {
			if heap.Poll() != v {
				t.Errorf("Expecting %q got %q", v, r)
			}
		}
	}
}

func buildWithRand(n int) (heap IHeap) {
	heap = &Heap{cap: n, useStaticArray: true}
	heap.Init()
	for i := 0; i < n; i++ { 
		v := IntEva(rand.Intn(n))
		heap.Add(v)
	}
	return
}

func peekAndPollAll(heap IHeap) {
	for !heap.IsEmpty() {
		if peek := heap.Peek(); peek != nil {
			poll := heap.Poll()
			if peek != poll {
				panic("peek and poll are diferente!")
			}
		}
	}
}

func BenchmarkMinimalHeap(b *testing.B) {
	heap := buildWithRand(b.N)
	peekAndPollAll(heap)
}
