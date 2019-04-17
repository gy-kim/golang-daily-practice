package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(bytes.ContainsAny([]byte("I like seafood."), "해산물."))
	fmt.Println(bytes.ContainsAny([]byte("I like seafood."), ""))
	fmt.Println(bytes.ContainsAny([]byte(""), ""))
}
