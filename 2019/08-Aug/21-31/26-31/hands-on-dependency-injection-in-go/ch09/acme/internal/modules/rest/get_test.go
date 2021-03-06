package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"

	"github.com/gorilla/mux"
	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/logging"
	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/data"

	"github.com/stretchr/testify/require"
)

func TestGetHandler_ServeHTTP(t *testing.T) {
	scenarios := []struct {
		desc            string
		inRequest       func() *http.Request
		inModelMock     func() *MockGetModel
		expectedStatus  int
		expectedPayload string
	}{
		{
			desc: "happy path",
			inRequest: func() *http.Request {
				req, err := http.NewRequest("GET", "/person/1/", nil)
				require.NoError(t, err)

				// set values into request (required by the mux)
				return mux.SetURLVars(req, map[string]string{muxVarID: "1"})
			},
			inModelMock: func() *MockGetModel {
				output := &data.Person{
					ID:       1,
					FullName: "John",
					Phone:    "0123456789",
					Currency: "USD",
					Price:    100,
				}

				mockGetModel := &MockGetModel{}
				mockGetModel.On("Do", mock.Anything).Return(output, nil).Once()

				return mockGetModel
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: `{"id":1,"name":"John","phone":"0123456789","currency":"USD","price":100}` + "\n",
		},
		{
			desc: "bad input (ID is invalid)",
			inRequest: func() *http.Request {
				req, err := http.NewRequest("GET", "/person/x/", nil)
				require.NoError(t, err)

				// set values into request (required by the mux)
				return mux.SetURLVars(req, map[string]string{muxVarID: "x"})
			},
			inModelMock: func() *MockGetModel {
				// expect the model not to be called
				mockRegisterModel := &MockGetModel{}
				return mockRegisterModel
			},
			expectedStatus:  http.StatusBadRequest,
			expectedPayload: ``,
		},
		{
			desc: "bad input (ID is missing)",
			inRequest: func() *http.Request {
				req, err := http.NewRequest("GET", "/person//", nil)
				require.NoError(t, err)

				// set values into request (required by the mux)
				return mux.SetURLVars(req, map[string]string{})
			},
			inModelMock: func() *MockGetModel {
				// expect the model not to be called
				mockRegisterModel := &MockGetModel{}
				return mockRegisterModel
			},
			expectedStatus:  http.StatusBadRequest,
			expectedPayload: ``,
		},
		{
			desc: "dependency fail",
			inRequest: func() *http.Request {
				req, err := http.NewRequest("GET", "/person/1/", nil)
				require.NoError(t, err)

				// set values into request (required by the mux)
				return mux.SetURLVars(req, map[string]string{muxVarID: "1"})
			},
			inModelMock: func() *MockGetModel {
				mockRegisterModel := &MockGetModel{}
				mockRegisterModel.On("Do", mock.Anything).Return(nil, errors.New("something failed")).Once()

				return mockRegisterModel
			},
			expectedStatus:  http.StatusNotFound,
			expectedPayload: ``,
		},
		{
			desc: "requested registration does not exist",
			inRequest: func() *http.Request {
				req, err := http.NewRequest("GET", "/person/1/", nil)
				require.NoError(t, err)

				// set values into request (required by the mux)
				return mux.SetURLVars(req, map[string]string{muxVarID: "1"})
			},
			inModelMock: func() *MockGetModel {
				mockRegisterModel := &MockGetModel{}
				mockRegisterModel.On("Do", mock.Anything).Return(nil, errors.New("person not found")).Once()

				return mockRegisterModel
			},
			expectedStatus:  http.StatusNotFound,
			expectedPayload: ``,
		},
	}

	for _, s := range scenarios {
		t.Run(s.desc, func(t *testing.T) {
			// define model layer mock
			modelGetModel := s.inModelMock()

			// build handler
			handler := NewGetHandler(&testConfig{}, modelGetModel)

			// perform request
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, s.inRequest())

			// validate outputs
			require.Equal(t, s.expectedStatus, response.Code, s.desc)

			payload, _ := ioutil.ReadAll(response.Body)
			assert.Equal(t, s.expectedPayload, string(payload), s.desc)
		})
	}
}

type testConfig struct {
}

func (t *testConfig) Logger() logging.Logger {
	return &logging.LoggerStdOut{}
}

func (t *testConfig) BindAddress() string {
	return "0.0.0.0:0"
}
