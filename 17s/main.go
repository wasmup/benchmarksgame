// The Computer Language Benchmarks Game
// https://salsa.debian.org/benchmarksgame-team/benchmarksgame/
//
// based on Jarkko Miettinen's Java program
// contributed by Tristan Dupont
// ported from Java version by wasmup

package main

import (
	"fmt"
	"os"
	"sync"
)

type node struct {
	left, right *node
}

func (p *node) check() int {
	if p.left == nil {
		return 1
	}
	return 1 + p.left.check() + p.right.check()
}

func create(depth int) *node {
	if depth > 1 {
		return &node{create(depth - 1), create(depth - 1)}
	}
	return &node{&node{}, &node{}}
}

func main() {
	const MinDepth = 4
	maxDepth := 10
	if len(os.Args) > 1 {
		_, err := fmt.Sscan(os.Args[1], &maxDepth)
		if err != nil {
			panic(err)
		}
		if MinDepth+2 > maxDepth {
			maxDepth = MinDepth + 2
		}
	}

	stretchDepth := maxDepth + 1
	fmt.Printf("stretch tree of depth %d\t check: %d\n",
		stretchDepth, create(stretchDepth).check())

	longLivedTree := create(maxDepth)

	wg := new(sync.WaitGroup)
	results := make([]string, (maxDepth-MinDepth)/2+1)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			depth := 2*i + MinDepth
			n := 1 << (maxDepth - depth + MinDepth)
			check := 0
			for j := n; j > 0; j-- {
				check += create(depth).check()
			}
			results[i] = fmt.Sprintf("%d\t trees of depth %d\t check: %d",
				n, depth, check)
			wg.Done()
		}(i)
	}
	wg.Wait()
	for _, s := range results {
		fmt.Println(s)
	}

	fmt.Printf("long lived tree of depth %d\t check: %d\n",
		maxDepth, longLivedTree.check())
}
