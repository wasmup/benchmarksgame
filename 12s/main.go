// The Computer Language Benchmarks Game
// https://salsa.debian.org/benchmarksgame-team/benchmarksgame/
//
// contributed by wasmup
//

package main

import (
	"fmt"
	"os"
	"sync"
)

type node struct {
	left, right *node
}

func create(depth int) *node {
	if depth == 1 {
		return &node{&node{}, &node{}}
	}
	return &node{create(depth - 1), create(depth - 1)}
}

func (p *node) check() int {
	sum := 1
	c := p.left
	for c != nil {
		sum += 2
		c = c.left
	}
	return sum
}

func main() {
	const MinDepth = 4
	const NoTasks = 4
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

	wg := new(sync.WaitGroup)
	stretchTreeCheck := ""
	wg.Add(1)
	go func() {
		stretchDepth := maxDepth + 1
		stretchTreeCheck = fmt.Sprintf("stretch tree of depth %d\t check: %d",
			stretchDepth, create(stretchDepth).check())
		wg.Done()
	}()

	longLivedText := ""
	wg.Add(1)
	go func() {
		longLivedText = fmt.Sprintf("long lived tree of depth %d\t check: %d",
			maxDepth, create(maxDepth).check())
		wg.Done()
	}()

	results := make([]string, (maxDepth-MinDepth)/2+1)
	for i := range results {
		depth := 2*i + MinDepth
		n := (1 << (maxDepth - depth + MinDepth)) / NoTasks
		tasks := make([]int, NoTasks)
		for t := range tasks {
			wg.Add(1)
			go func(t int) {
				check := 0
				for i := n; i > 0; i-- {
					check += create(depth).check()
				}
				tasks[t] = check
				wg.Done()
			}(t)
		}

		wg.Wait()
		check := 0
		for _, v := range tasks {
			check += v
		}
		results[i] = fmt.Sprintf("%d\t trees of depth %d\t check: %d",
			n*NoTasks, depth, check)
	}

	fmt.Println(stretchTreeCheck)

	for _, s := range results {
		fmt.Println(s)
	}

	fmt.Println(longLivedText)
}
