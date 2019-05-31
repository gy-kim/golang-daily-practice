package main

import "fmt"

type MasterAlgorithm struct {
	template Template
}

func (c *MasterAlgorithm) TemplateMethod() {
	c.template.Step1()
	c.template.Step2()
}

type Template interface {
	Step1()
	Step2()
}

type VariantA struct{}

func (c *VariantA) Step1() {
	fmt.Println("VariantA step 1")
}

func (c *VariantA) Step2() {
	fmt.Println("VariantA step 2")
}

func main() {
	masterAlgorithm := MasterAlgorithm{new(VariantA)}
	masterAlgorithm.TemplateMethod()
}
