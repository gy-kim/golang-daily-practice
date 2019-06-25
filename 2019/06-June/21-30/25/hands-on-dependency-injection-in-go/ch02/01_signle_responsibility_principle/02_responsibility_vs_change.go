package srp

import (
	"fmt"
	"io"
)

type CalculatorV2 struct {
	data map[string]float64
}

func (c CalculatorV2) Calculate(path string) error {
	return nil
}

func (c CalculatorV2) Output(writer io.Writer) {
	for path, result := range c.data {
		fmt.Fprintf(writer, "%s -> %.1f\n", path, result)
	}
}

func (c CalculatorV2) OutputCSV(writer io.Writer) {
	for path, result := range c.data {
		fmt.Fprintf(writer, "%s -> %.1f\n", path, result)
	}
}
