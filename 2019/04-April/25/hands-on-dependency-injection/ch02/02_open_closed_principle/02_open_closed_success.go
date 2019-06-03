package ocp

import "io"

type PersonFormatter interface {
	Format(writer io.Writer, person Person) error
}

// output the person as CSV
type CSVPersonFormatter struct{}

// Format implements the PersonFormatter interface
func (c *CSVPersonFormatter) Format(writer io.Writer, person Person) error {
	// TODO : implement
	return nil
}

// output the person as JSON
type JSONPersonFormatter struct{}

func (j *JSONPersonFormatter) Format(write io.Writer, person Person) error {
	// TODO: implement
	return nil
}

type Person struct {
	Name  string
	Email string
}
