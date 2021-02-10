// The Computer Language Benchmarks Game
// http://benchmarksgame.alioth.debian.org/
//
// Go adaptation of binary-trees Rust #4 program
// Use semaphores to match the number of workers with the CPU count
//
// contributed by Marcel Ibes
// modified by Isaac Gouy
// modified by Adam Shaver(use the struct constructor for bottomUpTree)
//
//

package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"runtime"
	"sort"
	"strconv"

	"golang.org/x/sync/semaphore"
)

type Tree struct {
	Left  *Tree
	Right *Tree
}

type Message struct {
	Pos  uint32
	Text string
}

var unity = uint32(1)

type ByPos []Message

func (m ByPos) Len() int           { return len(m) }
func (m ByPos) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByPos) Less(i, j int) bool { return m[i].Pos < m[j].Pos }

func itemCheck(tree *Tree) uint32 {
	if tree.Left != nil && tree.Right != nil {
		return uint32(1) + itemCheck(tree.Right) + itemCheck(tree.Left)
	}

	return 1
}

func bottomUpTree(depth uint32) *Tree {
	if depth > uint32(0) {
		return &Tree{Left: bottomUpTree(depth - 1), Right: bottomUpTree(depth - 1)}
	} else {
		return &Tree{}
	}
}

func inner(depth, iterations uint32) string {
	chk := uint32(0)
	for i := uint32(0); i < iterations; i++ {
		a := bottomUpTree(depth)
		chk += itemCheck(a)
	}
	return fmt.Sprintf("%d\t trees of depth %d\t check: %d",
		iterations, depth, chk)
}

const minDepth = uint32(4)

func main() {
	n := 0
	flag.Parse()
	if flag.NArg() > 0 {
		n, _ = strconv.Atoi(flag.Arg(0))
	}

	run(uint32(n))
}

func run(n uint32) {
	cpuCount := runtime.NumCPU()
	sem := semaphore.NewWeighted(int64(cpuCount))

	maxDepth := n
	if minDepth+2 > n {
		maxDepth = minDepth + 2
	}

	depth := maxDepth + 1

	messages := make(chan Message, cpuCount)
	expected := uint32(2) // initialize with the 2 summary messages

	go func() {
		// do stretch tree and longLivedTree
		if err := sem.Acquire(context.TODO(), 1); err == nil {
			go func() {
				defer sem.Release(1)
				tree := bottomUpTree(depth)
				messages <- Message{0,
					fmt.Sprintf("stretch tree of depth %d\t check: %d",
						depth, itemCheck(tree))}
			}()
		} else {
			panic(err)
		}

		if err := sem.Acquire(context.TODO(), 1); err == nil {
			go func() {
				defer sem.Release(1)
				longLivedTree := bottomUpTree(maxDepth)
				messages <- Message{math.MaxUint32,
					fmt.Sprintf("long lived tree of depth %d\t check: %d",
						maxDepth, itemCheck(longLivedTree))}
			}()
		} else {
			panic(err)
		}

		for halfDepth := minDepth / 2; halfDepth < maxDepth/2+1; halfDepth++ {
			depth := halfDepth * 2
			iterations := uint32(1 << (maxDepth - depth + minDepth))
			expected++

			func(d, i uint32) {
				if err := sem.Acquire(context.TODO(), 1); err == nil {
					go func() {
						defer sem.Release(1)
						messages <- Message{d, inner(d, i)}
					}()
				} else {
					panic(err)
				}
			}(depth, iterations)
		}
	}()

	var sortedMsg []Message
	for m := range messages {
		sortedMsg = append(sortedMsg, m)
		expected--
		if expected == 0 {
			close(messages)
		}
	}

	sort.Sort(ByPos(sortedMsg))
	for _, m := range sortedMsg {
		fmt.Println(m.Text)
	}
}

/*
NOTES:
64-bit Ubuntu quad core
go version go1.15.2 linux/amd64


Sat, 24 Oct 2020 02:29:05 GMT

MAKE:
/opt/src/go1.15.2/go/bin/go build -o binarytrees.go-8.go_run

3.07s to complete and log all make actions

COMMAND LINE:
./binarytrees.go-8.go_run 21

PROGRAM OUTPUT:
stretch tree of depth 22	 check: 8388607
2097152	 trees of depth 4	 check: 65011712
524288	 trees of depth 6	 check: 66584576
131072	 trees of depth 8	 check: 66977792
32768	 trees of depth 10	 check: 67076096
8192	 trees of depth 12	 check: 67100672
2048	 trees of depth 14	 check: 67106816
512	 trees of depth 16	 check: 67108352
128	 trees of depth 18	 check: 67108736
32	 trees of depth 20	 check: 67108832
long lived tree of depth 21	 check: 4194303
*/
