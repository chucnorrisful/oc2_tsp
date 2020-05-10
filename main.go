package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

const (
	MAX_DIST = 1000
	MIN_DIST = 1
)

func main() {
	var (
		testFlag    = flag.String("data", "", "Path to predefined distance matrix.")
		problemSize = flag.Int("n", 0, "Number of cities used.")
		seed        = flag.Int64("seed", 0, "Custom seed for random number generation.")
		cpuprofile  = flag.String("cpuprofile", "", "write cpu profile to file")
	)

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *testFlag != "" {
		a := loadTestMatrix(*testFlag)
		betteExhaustiveSearch(0, len(a[0]), a)

	} else {
		if *problemSize == 0 {
			panic("please provide the number of cities (-n 5)")
		}

		if *seed == 0 {
			*seed = time.Now().UnixNano()
		}

		d := createDistanceMatrix(*problemSize, *seed)
		betteExhaustiveSearch(0, *problemSize, d)
	}
}

func createDistanceMatrix(n int, source int64) [][]int {
	fmt.Println("crafting", n, "x", n, " distance matrix")
	r := rand.New(rand.NewSource(source))
	a := make([][]int, n)
	for i := range a {
		a[i] = make([]int, n)
	}

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if i == j {
				// Diagonal
				a[i][j] = 0
			} else if j > i {
				a[i][j] = r.Intn(MAX_DIST-MIN_DIST) + MIN_DIST
			} else {
				// Symmetric
				a[i][j] = a[j][i]
			}
		}
	}

	return a
}
func calcPathDist(path []int, dist [][]int) int {
	d := 0
	for i := 0; i < len(path)-1; i++ {
		d += dist[path[i]][path[i+1]]
	}
	return d
}

func loadTestMatrix(path string) [][]int {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var a [][]int

	err = json.Unmarshal(d, &a)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(a), "x", len(a), "test matrix loaded ")

	return a
}
