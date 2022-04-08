// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	shield "github.com/itimky/word-of-wisom/internal/service/shield"
)

// ShieldService is an autogenerated mock type for the shieldService type
type ShieldService struct {
	mock.Mock
}

type ShieldService_Expecter struct {
	mock *mock.Mock
}

func (_m *ShieldService) EXPECT() *ShieldService_Expecter {
	return &ShieldService_Expecter{mock: &_m.Mock}
}

// CheckPuzzle provides a mock function with given fields: clientIP, solution
func (_m *ShieldService) CheckPuzzle(clientIP string, solution *shield.PuzzleSolution) shield.PuzzleCheckResult {
	ret := _m.Called(clientIP, solution)

	var r0 shield.PuzzleCheckResult
	if rf, ok := ret.Get(0).(func(string, *shield.PuzzleSolution) shield.PuzzleCheckResult); ok {
		r0 = rf(clientIP, solution)
	} else {
		r0 = ret.Get(0).(shield.PuzzleCheckResult)
	}

	return r0
}

// ShieldService_CheckPuzzle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckPuzzle'
type ShieldService_CheckPuzzle_Call struct {
	*mock.Call
}

// CheckPuzzle is a helper method to define mock.On call
//  - clientIP string
//  - solution *shield.PuzzleSolution
func (_e *ShieldService_Expecter) CheckPuzzle(clientIP interface{}, solution interface{}) *ShieldService_CheckPuzzle_Call {
	return &ShieldService_CheckPuzzle_Call{Call: _e.mock.On("CheckPuzzle", clientIP, solution)}
}

func (_c *ShieldService_CheckPuzzle_Call) Run(run func(clientIP string, solution *shield.PuzzleSolution)) *ShieldService_CheckPuzzle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*shield.PuzzleSolution))
	})
	return _c
}

func (_c *ShieldService_CheckPuzzle_Call) Return(_a0 shield.PuzzleCheckResult) *ShieldService_CheckPuzzle_Call {
	_c.Call.Return(_a0)
	return _c
}
