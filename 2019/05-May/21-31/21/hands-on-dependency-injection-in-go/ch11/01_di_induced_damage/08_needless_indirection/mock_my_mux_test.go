// Code generated by mockery v1.0.0. DO NOT EDIT.

// @generated

package needless_indirection

import http "net/http"
import mock "github.com/stretchr/testify/mock"

// MockMyMux is an autogenerated mock type for the MyMux type
type MockMyMux struct {
	mock.Mock
}

// Handle provides a mock function with given fields: pattern, handler
func (_m *MockMyMux) Handle(pattern string, handler http.Handler) {
	_m.Called(pattern, handler)
}

// Handler provides a mock function with given fields: req
func (_m *MockMyMux) Handler(req *http.Request) (http.Handler, string) {
	ret := _m.Called(req)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(*http.Request) http.Handler); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(*http.Request) string); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// ServeHTTP provides a mock function with given fields: resp, req
func (_m *MockMyMux) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_m.Called(resp, req)
}
