package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// https://golang.org/pkg/bytes/#Buffer

func Buffer() {
	var b bytes.Buffer // A Buffer needs to ninitialization
	b.Write([]byte("Hello "))
	fmt.Fprintf(&b, "world!")
	b.WriteTo(os.Stdout)
}

func BufferReader() {
	// A Buffer can turn a string or a []byte into an io.Reader.
	buf := bytes.NewBufferString("R29waGVycyBydWxlIQ==")
	dec := base64.NewDecoder(base64.StdEncoding, buf)
	io.Copy(os.Stdout, dec)
}

func Grow() {
	var b bytes.Buffer
	b.Grow(64)
	bb := b.Bytes()
	b.Write([]byte("64 bytes or fewer"))
	fmt.Printf("%q", bb[:b.Len()])
}

func Len() {
	var b bytes.Buffer
	b.Grow(64)
	b.Write([]byte("abcd"))
	fmt.Printf("%d", b.Len())
}

func ReaderLen() {
	fmt.Println(bytes.NewReader([]byte("Hi!")).Len())
	fmt.Println(bytes.NewReader([]byte("안녕하세요!")).Len())
}

func main() {
	// Buffer()
	// BufferReader()
	// Grow()
	// Len()
	ReaderLen()
}
