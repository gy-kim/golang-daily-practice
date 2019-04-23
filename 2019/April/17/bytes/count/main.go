package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(bytes.Count([]byte("chees"), []byte("e")))
	fmt.Println(bytes.Count([]byte("five"), []byte(""))) // before & afte each rune
}
