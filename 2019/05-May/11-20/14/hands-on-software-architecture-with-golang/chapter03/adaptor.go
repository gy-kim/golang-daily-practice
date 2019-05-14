package main

import "fmt"

// Adaptee is the existing structure - something we need to use
type Adaptee struct{}

func (a *Adaptee) ExistingMethod() {
	fmt.Println("using existing method")
}

// Adapter is the structure we use to glue things together
type Adapter struct {
	adaptee *Adaptee
}

func NewAdapter() *Adapter {
	return &Adapter{new(Adaptee)}
}

func (a *Adapter) ExpectedMethod() {
	fmt.Println("doing some work")
	a.adaptee.ExistingMethod()
}

// func main() {
// 	adapter := NewAdapter()
// 	adapter.ExpectedMethod()
// }
