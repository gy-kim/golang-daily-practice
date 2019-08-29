// Code generated by mockery v1.0.0. DO NOT EDIT.

// @generated

package register

import context "context"
import data "github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/data"
import mock "github.com/stretchr/testify/mock"

// mockMySaver is an autogenerated mock type for the mySaver type
type mockMySaver struct {
	mock.Mock
}

// Save provides a mock function with given fields: ctx, in
func (_m *mockMySaver) Save(ctx context.Context, in *data.Person) (int, error) {
	ret := _m.Called(ctx, in)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, *data.Person) int); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *data.Person) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
