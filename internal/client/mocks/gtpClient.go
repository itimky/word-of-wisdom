// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	msgp "github.com/tinylib/msgp/msgp"

	server "github.com/itimky/word-of-wisom/api/server"
)

// GtpClient is an autogenerated mock type for the gtpClient type
type GtpClient struct {
	mock.Mock
}

type GtpClient_Expecter struct {
	mock *mock.Mock
}

func (_m *GtpClient) EXPECT() *GtpClient_Expecter {
	return &GtpClient_Expecter{mock: &_m.Mock}
}

// MakeRequest provides a mock function with given fields: reqType
func (_m *GtpClient) MakeRequest(reqType server.RequestType) (msgp.Raw, error) {
	ret := _m.Called(reqType)

	var r0 msgp.Raw
	if rf, ok := ret.Get(0).(func(server.RequestType) msgp.Raw); ok {
		r0 = rf(reqType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(msgp.Raw)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(server.RequestType) error); ok {
		r1 = rf(reqType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GtpClient_MakeRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakeRequest'
type GtpClient_MakeRequest_Call struct {
	*mock.Call
}

// MakeRequest is a helper method to define mock.On call
//  - reqType server.RequestType
func (_e *GtpClient_Expecter) MakeRequest(reqType interface{}) *GtpClient_MakeRequest_Call {
	return &GtpClient_MakeRequest_Call{Call: _e.mock.On("MakeRequest", reqType)}
}

func (_c *GtpClient_MakeRequest_Call) Run(run func(reqType server.RequestType)) *GtpClient_MakeRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(server.RequestType))
	})
	return _c
}

func (_c *GtpClient_MakeRequest_Call) Return(_a0 msgp.Raw, _a1 error) *GtpClient_MakeRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}
