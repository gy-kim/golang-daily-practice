package main

import (
	"bytes"
	"fmt"
	"os"
	"unicode"
)

func TrimLeftFunc() {
	fmt.Println(string(bytes.TrimLeftFunc([]byte("go-gopher"), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimLeftFunc([]byte("go-gopher"), unicode.IsPunct)))
	fmt.Println(string(bytes.TrimLeftFunc([]byte("1234go-gopher!567"), unicode.IsNumber)))
}

func TrimRight() {
	fmt.Print(string(bytes.TrimRight([]byte("453gopher8257"), "123456789")))
}

func TrimRightFunc() {
	fmt.Println(string(bytes.TrimRightFunc([]byte("go-gopher"), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimRightFunc([]byte("go-gopher!"), unicode.IsPunct)))
	fmt.Println(string(bytes.TrimRightFunc([]byte("1234go-gopher!567"), unicode.IsNumber)))
}

func TrimSpace() {
	fmt.Printf("%s", bytes.TrimSpace([]byte(" \t\n a long gopher \n\t\r\n")))
}

func TrimSuffix() {
	var b = []byte("Hello, goodbye, etc!")
	b = bytes.TrimSuffix(b, []byte("goodbye, etc!"))
	b = bytes.TrimSuffix(b, []byte("gopher"))
	b = append(b, bytes.TrimSuffix([]byte("world!"), []byte("x!"))...)
	os.Stdout.Write(b)
}

func Buffer() {
	var b bytes.Buffer // A Buffer needs no initialization
	b.Write([]byte("Hello "))
	fmt.Fprintf(&b, "world!")
	b.WriteTo(os.Stdout)
}

func main() {
	// TrimLeftFunc()
	// TrimRight()
	// TrimRightFunc()
	// TrimSpace()
	// TrimSuffix()
	Buffer()
}
