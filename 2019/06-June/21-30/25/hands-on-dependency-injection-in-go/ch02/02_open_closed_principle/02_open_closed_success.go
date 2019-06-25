package ocp

import (
	"io"
	"net/http"
)

func BuildOutputOCPSuccess(response http.ResponseWriter, formatter PersonFormatter, person Person) {
	err := formatter.Format(response, person)
	if err != nil {
		// output a server error and quit
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

type PersonFormatter interface {
	Format(writer io.Writer, person Person) error
}

type CSVPersonFormatter struct{}

func (c *CSVPersonFormatter) Format(writer io.Writer, person Person) error {
	return nil
}

// output the person as JSON
type JSONPersonFormatter struct{}

func (j *JSONPersonFormatter) Format(writer io.Writer, person Person) error {
	return nil
}
