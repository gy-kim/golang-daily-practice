package main

import (
	"fmt"
	"strings"
)

// ProductType is product type
type ProductType string

func (p ProductType) String() string {
	return string(p)
}

func (p ProductType) Equal(s string) bool {
	return strings.ToUpper(s) == p.String()
}

const (
	PTypeSMicro   ProductType = "MICRO"
	PTypesRegular ProductType = "REGULAR"
)

func main() {
	fmt.Println(PTypeSMicro)
	var micro string = ""

	micro = PTypeSMicro.String()
	fmt.Println(micro)

	fmt.Println(PTypeSMicro.Equal("micro"))
}
