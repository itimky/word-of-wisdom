// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// QuoteService is an autogenerated mock type for the quoteService type
type QuoteService struct {
	mock.Mock
}

type QuoteService_Expecter struct {
	mock *mock.Mock
}

func (_m *QuoteService) EXPECT() *QuoteService_Expecter {
	return &QuoteService_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields:
func (_m *QuoteService) Get() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// QuoteService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type QuoteService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
func (_e *QuoteService_Expecter) Get() *QuoteService_Get_Call {
	return &QuoteService_Get_Call{Call: _e.mock.On("Get")}
}

func (_c *QuoteService_Get_Call) Run(run func()) *QuoteService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *QuoteService_Get_Call) Return(_a0 string) *QuoteService_Get_Call {
	_c.Call.Return(_a0)
	return _c
}
