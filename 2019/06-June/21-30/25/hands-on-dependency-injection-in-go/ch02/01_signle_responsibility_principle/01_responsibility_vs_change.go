package srp

import (
	"fmt"
	"io"
)

// CalculatorV1 calaulates the test coverage for a directory and it's sub-directories
type CalculatorV1 struct {
	data map[string]float64
}

func (c *CalculatorV1) Calculate(path string) error {
	return nil
}

func (c *CalculatorV1) Output(writer io.Writer) {
	for path, result := range c.data {
		fmt.Fprintf(writer, "%s -> %.f\n", path, result)
	}
}
