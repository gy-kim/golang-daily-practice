package get

import (
	"context"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/1-10/02/hands-on-dependency-injection/ch09/acme/internal/modules/data"
	"github.com/stretchr/testify/mock"
)

type mockMyLoader struct {
	mock.Mock
}

func (_m *mockMyLoader) Load(ctx context.Context, ID int) (*data.Person, error) {
	ret := _m.Called(ctx, ID)

	var r0 *data.Person
	if rf, ok := ret.Get(0).(func(context.Context, int) *data.Person); ok {
		r0 = rf(ctx, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.Person)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
