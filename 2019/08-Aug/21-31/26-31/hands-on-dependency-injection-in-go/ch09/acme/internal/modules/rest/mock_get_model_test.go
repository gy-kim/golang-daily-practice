// Code generated by mockery v1.0.0. DO NOT EDIT.

// @generated

package rest

import data "github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/data"
import mock "github.com/stretchr/testify/mock"

// MockGetModel is an autogenerated mock type for the GetModel type
type MockGetModel struct {
	mock.Mock
}

// Do provides a mock function with given fields: ID
func (_m *MockGetModel) Do(ID int) (*data.Person, error) {
	ret := _m.Called(ID)

	var r0 *data.Person
	if rf, ok := ret.Get(0).(func(int) *data.Person); ok {
		r0 = rf(ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.Person)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
