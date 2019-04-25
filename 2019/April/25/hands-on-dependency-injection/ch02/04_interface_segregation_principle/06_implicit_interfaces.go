package isp

import "fmt"

type Talker interface {
	SayHello() string
}

type Dog struct{}

func (d Dog) SayHello() string {
	return "Woof!"
}

func Speak() {
	var talker Talker
	talker = Dog{}

	fmt.Println(talker.SayHello)
}
