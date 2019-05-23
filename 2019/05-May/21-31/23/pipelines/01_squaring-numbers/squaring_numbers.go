package main

import (
	"fmt"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	// Set up the pipeline.
	c := gen(2, 3)
	out := sq(c)

	// Consume the output.
	// fmt.Println(<-out)
	// fmt.Println(<-out)
	for n := range out {
		fmt.Println(n)
	}

	//////////////////////
	// Set up the pipeline and consume the output.
	// for n := range sq(sq(gen(2, 3))) {
	// 	fmt.Println(n)
	// }

}
