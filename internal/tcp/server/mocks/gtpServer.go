// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	gtp "github.com/itimky/word-of-wisom/pkg/gtp"
	mock "github.com/stretchr/testify/mock"

	server "github.com/itimky/word-of-wisom/internal/gtp/server"
)

// GtpServer is an autogenerated mock type for the gtpServer type
type GtpServer struct {
	mock.Mock
}

type GtpServer_Expecter struct {
	mock *mock.Mock
}

func (_m *GtpServer) EXPECT() *GtpServer_Expecter {
	return &GtpServer_Expecter{mock: &_m.Mock}
}

// CheckPuzzle provides a mock function with given fields: clientIP, solution
func (_m *GtpServer) CheckPuzzle(clientIP string, solution *gtp.PuzzleSolution) server.PuzzleCheckResult {
	ret := _m.Called(clientIP, solution)

	var r0 server.PuzzleCheckResult
	if rf, ok := ret.Get(0).(func(string, *gtp.PuzzleSolution) server.PuzzleCheckResult); ok {
		r0 = rf(clientIP, solution)
	} else {
		r0 = ret.Get(0).(server.PuzzleCheckResult)
	}

	return r0
}

// GtpServer_CheckPuzzle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckPuzzle'
type GtpServer_CheckPuzzle_Call struct {
	*mock.Call
}

// CheckPuzzle is a helper method to define mock.On call
//  - clientIP string
//  - solution *gtp.PuzzleSolution
func (_e *GtpServer_Expecter) CheckPuzzle(clientIP interface{}, solution interface{}) *GtpServer_CheckPuzzle_Call {
	return &GtpServer_CheckPuzzle_Call{Call: _e.mock.On("CheckPuzzle", clientIP, solution)}
}

func (_c *GtpServer_CheckPuzzle_Call) Run(run func(clientIP string, solution *gtp.PuzzleSolution)) *GtpServer_CheckPuzzle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*gtp.PuzzleSolution))
	})
	return _c
}

func (_c *GtpServer_CheckPuzzle_Call) Return(_a0 server.PuzzleCheckResult) *GtpServer_CheckPuzzle_Call {
	_c.Call.Return(_a0)
	return _c
}
