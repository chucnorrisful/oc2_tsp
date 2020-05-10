package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gosuri/uiprogress"
)

// First try, overflows memory on large problems
// since all permutations are created and stored in memory.
func exhaustiveSearch(src, n int, dist [][]int) {
	fmt.Println("exhaustive search, there are ", fac(n-1), " paths to check")

	perms := genPermutations(src, n)

	var (
		path     []int
		shortest int = math.MaxInt32
	)

	for _, perm := range perms {
		// Check if path ist shorter than path
		d := calcPathDist(perm, dist)
		if d < shortest {
			shortest = d
			path = perm
		}
	}
	fmt.Println("shortest path is ", path, "len", shortest)
}

// Do not preallocate all possible paths, be merciful to your memory.
func betteExhaustiveSearch(src, n int, dist [][]int) {
	fmt.Println("better exhaustive search, there are ", fac(n-1), " paths to check \nsearching ...\n")

	var (
		path         []int
		shortest     int = math.MaxInt32
		shortestTemp int
		a            = make([]int, n+1)
		rate         float64
		tLeft        float64
		perm         = sliceWithoutSrc(src, n)
		n_len        = len(perm) - 1
		i, j, k      int
	)

	// Uhhh yes a fancy progressbar

	uiprogress.Start()
	bar := uiprogress.AddBar(fac(n - 1))
	bar.Width = 42
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		tLeft = float64(bar.Total-bar.Current()) / rate
		return "\tpaths/s\t\t" + strconv.FormatFloat(rate, 'f', 0, 64) + "\n" +
		    "\ttime left\t" + time.Duration(tLeft*1000000000).String() + "\n" +
		    "\tshortest\t" + strconv.Itoa(shortest) + "\n" +
		    " " + strconv.FormatFloat(time.Since(b.TimeStarted).Seconds(), 'f', 0, 64) + "s"
	})
	bar.AppendCompleted()

	// Start/End are the same for every path
	a[0] = src
	a[n] = src

	for c := 1; c < fac(n-1); c++ {
		bar.Incr()

		// For the sake of performance permutations are generated inline
		i = n_len - 1
		j = n_len
		for perm[i] > perm[i+1] {
			i--
		}
		for perm[j] < perm[i] {
			j--
		}
		perm[i], perm[j] = perm[j], perm[i]
		j = n_len
		i += 1
		for i < j {
			perm[i], perm[j] = perm[j], perm[i]
			i++
			j--
		}

		// Squeeze a permutation of the free indices in
		// start=end=0: 0 - x - x - x 0
		for k = 0; k < len(perm); k++ {
			a[k+1] = perm[k]
		}

		// Check the path

		shortestTemp = calcPathDist(a, dist)
		if shortestTemp < shortest {
			shortest = shortestTemp
			path = a
		}

		// Compute the current computation rate every now and then,
		// this does slow down everything a bit...

		if bar.Current()%100000 == 0 {
			rate = float64(bar.Current()) / time.Since(bar.TimeStarted).Seconds()
		}

	}

	exportD(dist)
	exportPath(path)
	draw()

	fmt.Println("shortest path is ", path, "len", shortest, " took ", time.Since(bar.TimeStarted))
}

// Faculty n!
func fac(n int) (result int) {
	if n > 0 {
		result = n * fac(n-1)
		return result
	}
	return 1
}

func sliceWithoutSrc(v, n int) []int {
	var left []int
	for i := 0; i < n; i++ {
		if i == v {
			continue
		}
		left = append(left, i)
	}
	return left

}

// This does work, but since all possible paths are generated beforehand
// we may run into a little RAM overflow if the problem is too big...
// In my case (32gb) if n>13
func genPermutations(start int, n int) [][]int {
	a := make([][]int, fac(n-1))

	// Indices without/start/end
	var left []int
	for i := 0; i < n; i++ {
		if i == start {
			continue
		}
		left = append(left, i)
	}

	tmp := 0
	for perm := range permutations(left) {
		a[tmp] = make([]int, n+1)
		a[tmp][0] = start
		a[tmp][n] = start
		k := 1
		for i := 0; i < len(perm); i++ {
			a[tmp][k] = perm[i]
			k++
		}
		tmp++
	}
	return a
}

// Based on on the QuickPerm algorithm,
// unfortunately the overhead of the go runtime mitigate
// any performance improvements compared to a simple inline heaps algorithm.
func permutations(data []int) <-chan []int {
	c := make(chan []int)
	go func(c chan []int) {
		defer close(c)
		permutate(c, data)
	}(c)
	return c
}
func permutate(c chan []int, inputs []int) {
	output := make([]int, len(inputs))
	copy(output, inputs)
	c <- output

	size := len(inputs)
	p := make([]int, size+1)
	for i := 0; i < size+1; i++ {
		p[i] = i
	}
	for i := 1; i < size; {
		p[i]--
		j := 0
		if i%2 == 1 {
			j = p[i]
		}
		tmp := inputs[j]
		inputs[j] = inputs[i]
		inputs[i] = tmp
		output := make([]int, len(inputs))
		copy(output, inputs)
		c <- output
		for i = 1; p[i] == 0; i++ {
			p[i] = i
		}
	}
}


func permutateParallel(c chan []int, inputs []int) {

	finisher := make(chan bool)
	for i := 0; i < len(inputs); i++ {

		temp := make([]int, len(inputs)-1)
		copy(temp, inputs[:i])
		copy(temp[i:], inputs[i+1:])

		//fmt.Println(temp, i)

		go permutate2(c, temp, inputs[i], finisher)
	}

	for i := 0; i < len(inputs); i++ {
		<-finisher
	}

	return
}
func permutate2(c chan []int, inputs []int, first int, finisher chan bool) {
	output := make([]int, len(inputs))
	copy(output, inputs)
	c <- output

	size := len(inputs)
	p := make([]int, size+1)
	for i := 0; i < size+1; i++ {
		p[i] = i
	}
	for i := 1; i < size; {
		p[i]--
		j := 0
		if i%2 == 1 {
			j = p[i]
		}
		inputs[j], inputs[i] = inputs[i], inputs[j]
		output := make([]int, len(inputs)+1)
		copy(output[1:], inputs)
		output[0] = first
		c <- output
		for i = 1; p[i] == 0; i++ {
			p[i] = i
		}
	}
	finisher <- true
}
func permutate3(c chan []int, inputs []int, first int, dist [][]int) {
	var path []int
	var shortest = math.MaxInt32


	d := calcPathDist(inputs, dist)
	if d < shortest {
		shortest = d
		output := make([]int, len(inputs)+1)
		copy(output[1:], inputs)
		output[0] = first
		path = output
	}

	size := len(inputs)
	p := make([]int, size+1)
	for i := 0; i < size+1; i++ {
		p[i] = i
	}
	for i := 1; i < size; {
		p[i]--
		j := 0
		if i%2 == 1 {
			j = p[i]
		}
		inputs[j], inputs[i] = inputs[i], inputs[j]

		d := calcPathDist(inputs, dist)
		if d < shortest {
			shortest = d
			output := make([]int, len(inputs)+1)
			copy(output[1:], inputs)
			output[0] = first
			path = output
		}

		for i = 1; p[i] == 0; i++ {
			p[i] = i
		}
	}

	c <- path
}

func parallelPermutations(data []int, buffer int) <-chan []int {
	c := make(chan []int, buffer)
	go func(c chan []int) {
		defer close(c)
		permutateParallel(c, data)
	}(c)
	return c
}
