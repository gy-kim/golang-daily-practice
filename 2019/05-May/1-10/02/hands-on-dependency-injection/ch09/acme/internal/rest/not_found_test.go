package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotFound_ServeHTTP(t *testing.T) {
	// build inputs
	response := httptest.NewRecorder()
	request := &http.Request{}

	// call handler
	notFoundHandler(response, request)

	// validate output
	require.Equal(t, http.StatusNotFound, response.Code)
}
