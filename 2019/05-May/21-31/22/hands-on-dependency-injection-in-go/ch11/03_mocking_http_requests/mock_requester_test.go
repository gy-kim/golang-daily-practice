// Code generated by mockery v1.0.0. DO NOT EDIT.

// @generated

package mocking_http_requests

import context "context"
import http "net/http"
import mock "github.com/stretchr/testify/mock"

// MockRequester is an autogenerated mock type for the Requester type
type MockRequester struct {
	mock.Mock
}

// doRequest provides a mock function with given fields: ctx, url
func (_m *MockRequester) doRequest(ctx context.Context, url string) (*http.Response, error) {
	ret := _m.Called(ctx, url)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(context.Context, string) *http.Response); ok {
		r0 = rf(ctx, url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}