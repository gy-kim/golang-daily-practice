package rest

// import (
// 	"net/http"
// 	"testing"

// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"

// 	"github.com/gy-kim/golang-daily-practice/2019/05-May/1-10/02/hands-on-dependency-injection/ch09/acme/internal/logging"
// 	"github.com/gy-kim/golang-daily-practice/2019/05-May/1-10/02/hands-on-dependency-injection/ch09/acme/internal/modules/data"
// )

// func TestGetHandler_ServeHTTP(t *testing.T) {
// 	scenarios := []struct {
// 		desc            string
// 		inRequest       func() *http.Request
// 		inModelMock     func() *MockGetModel
// 		expectedStatus  int
// 		expectedPayload string
// 	}{
// 		{
// 			desc: "happy path",
// 			inRequest: func() *http.Request {
// 				req, err := http.NewRequest("GET", "/person/1/", nil)
// 				require.NoError(t, err)

// 				// set values into request (required by the mux)
// 				return mux.SetURLVars(req, map[string]string{muxVarID: "1"})
// 			},
// 			inModelMock: func() *MockGetModel {
// 				output := &data.Person{
// 					ID:       1,
// 					FullName: "John",
// 					Phone:    "0123456789",
// 					Price:    100,
// 				}

// 				mockGetModel := &MockGetModel{}
// 				mockGetModel.On("Do", mock.Anything).Return(output, nil).Once()

// 				return mockGetModel
// 			},
// 			expectedStatus:  http.StatusOK,
// 			expectedPayload: `{"id":1,"name":"John","phone":"0123456789","currency":"USD","price":100}` + "\n",
// 		}, {
// 			desc: "bad input (ID is invalid)",
// 			inRequest: func() *http.Request {
// 				req, err := http.NewRequest("GET", "/person/x/", nil)
// 				require.NoError(t, err)

// 				// set values into request (required by the mux)
// 				return mux.SetURLVars(req, map[string]string{muxVarID: "x"})
// 			},
// 			expectedStatus:  http.StatusBadRequest,
// 			expectedPayload: ``,
// 		},

// 	}
// }

// type testConfig struct{}

// func (t *testConfig) Logger() logging.Logger {
// 	return &logging.LoggerStdOut{}
// }

// func (*testConfig) BindAddress() string {
// 	return "0.0.0.0:0"
// }
