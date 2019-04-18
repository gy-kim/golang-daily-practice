package main

import (
	"bytes"
	"fmt"
	"unicode"
)

func Join() {
	s := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
	fmt.Printf("%s", bytes.Join(s, []byte(", ")))
}

func LastIndex() {
	fmt.Println(bytes.Index([]byte("go gopher"), []byte("go")))
	fmt.Println(bytes.LastIndex([]byte("go gopher"), []byte("go")))
	fmt.Println(bytes.LastIndex([]byte("go gopher"), []byte("rodent")))
}

func LastIndexAny() {
	fmt.Println(bytes.LastIndexAny([]byte("go gopher"), "MuQp"))
	fmt.Println(bytes.LastIndexAny([]byte("go gopher"), "z,!."))
}

func LastIndexByte() {
	fmt.Println(bytes.LastIndexByte([]byte("go gopher"), byte('g')))
	fmt.Println(bytes.LastIndexByte([]byte("go gopher"), byte('r')))
	fmt.Println(bytes.LastIndexByte([]byte("go gopher"), byte('z')))
}

func LastIndexFunc() {
	fmt.Println(bytes.LastIndexFunc([]byte("go gopher!"), unicode.IsLetter))
	fmt.Println(bytes.LastIndexFunc([]byte("go gopher!"), unicode.IsPunct))
	fmt.Println(bytes.LastIndexFunc([]byte("go gopher!"), unicode.IsNumber))
}

func Map() {
	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26
		}
		return r
	}
	fmt.Printf("%s", bytes.Map(rot13, []byte("'Twas brillig and the slithy gopher...")))
}

func Repeat() {
	fmt.Printf("ba%s", bytes.Repeat([]byte("na"), 2))
}

func Replace() {
	fmt.Printf("%s\n", bytes.Replace([]byte("oink oink oink"), []byte("k"), []byte("ky"), 2))
	fmt.Printf("%s\n", bytes.Replace([]byte("oink oink oink"), []byte("oink"), []byte("moo"), -1))
}

func ReplaceAll() {
	// bytes.ReplaceAll was added in Go1.12
	// fmt.Printf("%s\n", bytes.ReplaceAll([]byte("oink oink oink"), []byte("oink"), []byte("moo")))
}

func Runes() {
	rs := bytes.Runes([]byte("go gopher"))
	for _, r := range rs {
		fmt.Printf("%#U\n", r)
	}
}

func Split() {
	fmt.Printf("%q\n", bytes.Split([]byte("a,b,c"), []byte(",")))
	fmt.Printf("%q\n", bytes.Split([]byte("a man a plan a cancel panama"), []byte("a ")))
	fmt.Printf("%q\n", bytes.Split([]byte("  xyz   "), []byte("")))
	fmt.Printf("%q\n", bytes.Split([]byte(""), []byte("Bernardo O'Higgins")))
}

func SplitAfter() {
	fmt.Printf("%q\n", bytes.SplitAfter([]byte("a,b,c"), []byte(",")))
}

func SplitAfterN() {
	fmt.Printf("%q\n", bytes.SplitAfterN([]byte("a,b,c,d,e,f"), []byte(","), 3))
}

func SplitN() {
	fmt.Printf("%q\n", bytes.SplitN([]byte("a,b,c"), []byte(","), 2))
	z := bytes.SplitN([]byte("a,b,c"), []byte(","), 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
}

func Title() {
	fmt.Printf("%s", bytes.Title([]byte("her royal highness")))
}

func ToLower() {
	fmt.Printf("%s", bytes.ToLower([]byte("Gopher")))
}

func ToTitle() {
	fmt.Printf("%s\n", bytes.ToTitle([]byte("loud noises")))
	fmt.Printf("%s\n", bytes.ToTitle([]byte("хлеб")))
}

func Trim() {
	fmt.Printf("[%q]", bytes.Trim([]byte("  !!! Achtung! Achtung! !!! "), "! "))
}

func TrimFunc() {
	fmt.Println(string(bytes.TrimFunc([]byte("go-gopher!"), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimFunc([]byte("\"go-gopher!\""), unicode.IsLetter)))
	fmt.Println(string(bytes.TrimFunc([]byte("go-gopher!"), unicode.IsPunct)))
	fmt.Println(string(bytes.TrimFunc([]byte("1234go-gopher!567"), unicode.IsNumber)))
}

func TrimLeft() {
	fmt.Print(string(bytes.TrimLeft([]byte("453gopher8257"), "0123456789")))
}

func main() {
	// Join()
	// LastIndex()
	// LastIndexAny()
	// LastIndexByte()
	// LastIndexFunc()
	// Map()
	// Repeat()
	// Replace()
	// ReplaceAll()
	// Runes()
	// Split()
	// SplitAfter()
	// SplitAfterN()
	// SplitN()
	// Title()
	// ToLower()
	// ToTitle()
	// Trim()
	// TrimFunc()
	// TrimLeft()
}
