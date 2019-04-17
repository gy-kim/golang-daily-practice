package main

import (
	"bytes"
	"fmt"
	"unicode"
)

// https://golang.org/pkg/bytes/

func main() {
	// Fields()
	// Fieldsfunc()
	// HasPrefix()
	// HasSuffix()
	// Index()
	// IndexAny()
	// IndexByte()
	// IndexFunc()
	// IndexRune()
	Join()
}

func Fields() {
	fmt.Printf("Fields are: %q\n", bytes.Fields([]byte("  foo  bar baz  ")))

	arr := bytes.Fields([]byte("   foo bar   baz  "))
	for _, a := range arr {
		fmt.Println(string(a))
	}
}

func Fieldsfunc() {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fmt.Printf("Fields are: %q", bytes.FieldsFunc([]byte(" foo1;bar2;baz3..."), f))
}

func HasPrefix() {
	fmt.Println(bytes.HasPrefix([]byte("Gopher"), []byte("Go")))
	fmt.Println(bytes.HasPrefix([]byte("Gopher"), []byte("C")))
	fmt.Println(bytes.HasPrefix([]byte("Gopher"), []byte("")))
}

func HasSuffix() {
	fmt.Println(bytes.HasSuffix([]byte("Amigo"), []byte("go")))
	fmt.Println(bytes.HasSuffix([]byte("Amigo"), []byte("O")))
	fmt.Println(bytes.HasSuffix([]byte("Amigo"), []byte("Ami")))
	fmt.Println(bytes.HasSuffix([]byte("Amigo"), []byte("")))
}

func Index() {
	fmt.Println(bytes.Index([]byte("chicken"), []byte("ken")))
	fmt.Println(bytes.Index([]byte("chicken"), []byte("dmr")))
}

func IndexAny() {
	fmt.Println(bytes.IndexAny([]byte("chicken"), "aeiouy"))
	fmt.Println(bytes.IndexAny([]byte("crwth"), "aeiouy"))
}

func IndexByte() {
	fmt.Println(bytes.IndexByte([]byte("chicken"), byte('k')))
	fmt.Println(bytes.IndexByte([]byte("chicken"), byte('g')))
}

func IndexFunc() {
	f := func(c rune) bool {
		return unicode.Is(unicode.Hangul, c)
	}
	fmt.Println(bytes.IndexFunc([]byte("Hello, 월드"), f))
	fmt.Println(bytes.IndexFunc([]byte("Hello, world"), f))
}

func IndexRune() {
	fmt.Println(bytes.IndexRune([]byte("chicken"), 'k'))
	fmt.Println(bytes.IndexRune([]byte("chicken"), 'd'))
}

func Join() {

}
