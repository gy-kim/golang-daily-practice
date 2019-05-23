package main

import (
	"fmt"
	"sync"
)

// Multiple functions can read from the same channel until that channel is closed; this is called fan-out.
// This provides a way to distribute work amongst a group of workers to parallelize CPU use and I/O.

func main() {
	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)

	// Consume the merged output from c1 and c2.
	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}

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

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs. output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		fmt.Println("output c:", c)
		for n := range c {
			out <- n
		}
		fmt.Println("wg.Done()")
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are done.
	// This must start after the wg.Add call.
	go func() {
		fmt.Println("wg.Wait()")
		wg.Wait()
		fmt.Println("close(out)")
		close(out)
	}()
	return out
}
