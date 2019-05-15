package main

import (
	"fmt"
	"math/rand"
	"time"
)

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string) // The FanIn channel

	go func() {
		for {
			c <- <-input1 // Write the message to the FanIn channel, Blocking Call
		}
	}()

	go func() {
		for {
			c <- <-input2 // Write the
		}
	}()

	return c
}

func emitter(name string) <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("[%s] says %d", name, i)
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
	}()
	return c
}

func main() {
	c := fanIn(emitter("Source1"), emitter("Source2"))

	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
}
