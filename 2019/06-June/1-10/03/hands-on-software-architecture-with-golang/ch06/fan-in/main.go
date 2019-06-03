package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c := fanIn(emitter("Source1"), emitter("Source2"))

	for i := 0; i < 10; i++ {
		fmt.Println(<-c) // Display the output of the FanIn channel.
	}
}

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string) // The FanIn channel

	// to avoid blocking, listen to the input channels in separate goroutinies
	go func() {
		for {
			c <- <-input1 // Write the message to the FanIn channel, Blocking call.
		}
	}()

	go func() {
		for {
			c <- <-input2 // Write the message to the FanIn channel, Blocking call.
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
