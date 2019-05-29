package main

import "fmt"

// Adaptee is the existing structure - something we need to use
type Adaptee struct{}

func (a *Adaptee) ExistingMethod() {
	fmt.Println("using existing method.")
}

// Adapter is the structure we use to glue things together
type Adapter struct {
	adaptee *Adaptee
}

func NewAdapter() *Adapter {
	return &Adapter{new(Adaptee)}
}

// ExpectedMethod is the method clients in current code are using. This honors the exptected interface and fulfils it
// using the Adaptee's method
func (a *Adapter) ExpectedMethod() {
	fmt.Println("doing some work")
	a.adaptee.ExistingMethod()
}

func main() {
	adaptor := NewAdapter()
	adaptor.ExpectedMethod()
}
