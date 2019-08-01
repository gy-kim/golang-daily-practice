package srp

import (
	"fmt"
	"io"
)

type CalculatorV3 struct {
	data map[string]float64
}

func (c *CalculatorV3) Calculate(path string) error {
	return nil
}

func (c *CalculatorV3) getData() map[string]float64 {
	return nil
}

type Printer interface {
	Output(data map[string]float64)
}

type DefaultPrinter struct {
	Writer io.Writer
}

func (d *DefaultPrinter) Output(data map[string]float64) {
	for path, result := range data {
		fmt.Fprintf(d.Writer, "%s -> %.1f\n", path, result)
	}
}

type CSVPrinter struct {
	Writer io.Writer
}

func (d *CSVPrinter) Output(data map[string]float64) {
	for path, result := range data {
		fmt.Fprintf(d.Writer, "%s, %.1f\n", path, result)
	}
}
