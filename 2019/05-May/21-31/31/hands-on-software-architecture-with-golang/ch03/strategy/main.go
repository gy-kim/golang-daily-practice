package main

import (
	"fmt"
	"math"
	"sort"
)

type Stategy interface {
	FindBreadth([]int) int
}

type NativeAlgo struct{}

func (n *NativeAlgo) FindBreadth(set []int) int {
	sort.Ints(set)
	return set[len(set)-1] - set[0]
}

type FastAlgo struct{}

func (n *FastAlgo) FindBreadth(set []int) int {
	min := math.MaxInt32
	max := math.MinInt64

	for _, x := range set {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}

	return max - min
}

func client(s Stategy) int {
	a := []int{-1, 10, 3, 1}
	return s.FindBreadth(a)
}

func main() {
	fmt.Println(client(&NativeAlgo{}))
	fmt.Println(client(&FastAlgo{}))
}
