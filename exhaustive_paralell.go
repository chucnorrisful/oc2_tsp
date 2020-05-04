package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func testPar() [][]int{
	fmt.Println("testing parallel execution")
	// d := loadTestMatrix("testData2.json")
	d := createDistanceMatrix(11,0)
	//
	// start := time.Now()
	// //betteExhaustiveSearch(0, *problemSize, d)
	in := make(chan []int, 10000)
	out := make(chan int, 5000 )

	s := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)

	for i := 0; i < 10; i++ {
		go process(in, out, &wg, i, d)
	}
	go processOut(out)
	genIn(in, &wg, 0, len(d[0]))
	wg.Wait()
	fmt.Println("ar took ",time.Since(s))
	return d
}
func processOut(out chan int) {
	smallest := math.MaxInt32
	for a := range out {
		if a < smallest {
			fmt.Println(a)
			smallest = a
		}

	}
	fmt.Println("SMALLEST IS ",smallest)
}
func process(in chan []int, out chan int, wg *sync.WaitGroup, k int, dist [][]int) {
	for i := range in {
		d := calcPathDist(i, dist)
		out <- d
	}
}

func genIn(in chan []int, wg *sync.WaitGroup, src, n int) {
	// Indices without/start/end
	left := sliceWithoutSrc(src, n)

	for perm := range permutations(left) {

		// Make a path
		a := make([]int, n+1)
		a[0] = src
		a[n] = src
		k := 1
		for i := 0; i < len(perm); i++ {
			a[k] = perm[i]
			k++
		}
		// Put on input channel
		in <- a
	}
	time.Sleep(time.Millisecond*100)
	wg.Done()
}
