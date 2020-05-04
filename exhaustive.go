package main

import (
	"fmt"
	"math"

	"github.com/gosuri/uiprogress"
)

// 1: [0,1,2,0]
// 2: 0210

func exhaustiveSearch(src, n int, dist [][]int) {
	fmt.Println("exhaustive search, there are ", fac(n-1), " paths to check")

	perms := genPermutations(src, n)

	var path []int
	var shortest int = math.MaxInt32

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
	fmt.Println("better exhaustive search, there are ", fac(n-1), " paths to check")

	var path []int
	var shortest int = math.MaxInt32

	// Fancy
	uiprogress.Start()                   // start rendering
	bar := uiprogress.AddBar(fac(n - 1)) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	bar.AppendCompleted()
	// bar.PrependElapsed()
	bar.AppendElapsed()

	// Indices without/start/end
	left := sliceWithoutSrc(src, n)

	for perm := range permutations(left) {
		bar.Incr()

		// Make a path
		a := make([]int, n+1)
		a[0] = src
		a[n] = src
		k := 1
		for i := 0; i < len(perm); i++ {
			a[k] = perm[i]
			k++
		}

		// Check it
		d := calcPathDist(a, dist)
		if d < shortest {
			shortest = d
			path = a
		}

	}

	uiprogress.Stop()

	fmt.Println(" B shortest path is ", path, "len", shortest)
}

//uses parallel permutator, otherwise identical to betteExhaustiveSearch
func betteExhaustiveSearch2(src, n int, dist [][]int) {
	fmt.Println("better exhaustive search, there are ", fac(n-1), " paths to check")

	var path []int
	var shortest int = math.MaxInt32

	// Fancy
	uiprogress.Start()                   // start rendering
	bar := uiprogress.AddBar(fac(n - 1)) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	bar.AppendCompleted()
	// bar.PrependElapsed()
	bar.AppendElapsed()

	// Indices without/start/end
	left := sliceWithoutSrc(src, n)

	for perm := range parallelPermutations(left, 20) {
		bar.Incr()

		// Make a path
		a := make([]int, n+1)
		a[0] = src
		a[n] = src
		k := 1
		for i := 0; i < len(perm); i++ {
			a[k] = perm[i]
			k++
		}

		// Check it
		d := calcPathDist(a, dist)
		if d < shortest {
			shortest = d
			path = a
		}

	}

	fmt.Println("shortest path is ", path, "len", shortest)
}

//runs fastest if workerCnt==1 :(
func parallelExhaustiveSearch(workerCnt, src, n int, dist [][]int) {
	fmt.Println("parallel exhaustive search, there are ", fac(n-1), " paths to check")
	fmt.Println(workerCnt, " workers started")

	var path []int
	var shortest = math.MaxInt32

	// Fancy
	//uiprogress.Start()                   // start rendering
	//bar := uiprogress.AddBar(fac(n - 1)) // Add a new bar
	// optionally, append and prepend completion and elapsed time
	//bar.AppendCompleted()
	// bar.PrependElapsed()
	//bar.AppendElapsed()

	// Indices without/start/end
	left := sliceWithoutSrc(src, n)

	//startWorkers
	in := parallelPermutations(left, workerCnt*2)

	outs := make([]chan []int, workerCnt)
	for i := 0; i < workerCnt; i++ {
		outs[i] = exhaustiveWorker(i, src, n, dist, in)
	}

	for i := 0; i < workerCnt; i++ {
		short := <-outs[i] //blocking, unbuffered out chan
		d := calcPathDist(short, dist)
		if d < shortest {
			shortest = d
			path = short
		}
	}

	fmt.Println("shortest path is ", path, "len", shortest)
}
func exhaustiveWorker(worker, src, n int, dist [][]int, in <-chan []int) (out chan []int) {
	out = make(chan []int, 0)

	go func() {

		shortest := math.MaxInt32
		var path []int

		cnt := 0
		for perm := range in {
			cnt++
			if cnt%1000000 == 0 && worker == 0 {
				//bar.Incr()
				fmt.Println(len(in))
			}

			// Make a path
			a := make([]int, n+1)
			a[0] = src
			a[n] = src

			copy(a[1:], perm)

			// Check it
			d := calcPathDist(a, dist)
			if d < shortest {
				shortest = d
				path = a
			}
		}

		out <- path
		fmt.Println("shortest path worker ", worker, " is ", path, "len", shortest)
	}()

	return out
}

func parallelExhaustiveSearch2(workerCnt, src, n int, dist [][]int) {
	fmt.Println("parallel exhaustive search, there are ", fac(n-1), " paths to check")
	fmt.Println(workerCnt, " workers started")

	//c chan []int, inputs []int
	inputs := sliceWithoutSrc(src, n)

	c := make([]chan []int, len(inputs))

	for i := 0; i < len(inputs); i++ {
		c[i] = make(chan []int)

		temp := make([]int, len(inputs)-1)
		copy(temp, inputs[:i])
		copy(temp[i:], inputs[i+1:])

		go func(t []int, ou chan []int, in int) {
			permutate3(ou, t, in, dist)
		}(temp, c[i], inputs[i])
	}

	var path []int
	var shortest = math.MaxInt32

	for i := 0; i < len(inputs); i++ {
		aChan := c[i]
		aPath := <-aChan
		d := calcPathDist(aPath, dist)
		if d < shortest {
			shortest = d
			path = aPath
		}
	}

	fmt.Println("A shortest path is ", path, "len", shortest)

	return
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
