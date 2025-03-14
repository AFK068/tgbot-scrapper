// Code generated by mockery v2.52.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UserChecker is an autogenerated mock type for the UserChecker type
type UserChecker struct {
	mock.Mock
}

type UserChecker_Expecter struct {
	mock *mock.Mock
}

func (_m *UserChecker) EXPECT() *UserChecker_Expecter {
	return &UserChecker_Expecter{mock: &_m.Mock}
}

// CheckUserExistence provides a mock function with given fields: id
func (_m *UserChecker) CheckUserExistence(id int64) bool {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for CheckUserExistence")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// UserChecker_CheckUserExistence_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckUserExistence'
type UserChecker_CheckUserExistence_Call struct {
	*mock.Call
}

// CheckUserExistence is a helper method to define mock.On call
//   - id int64
func (_e *UserChecker_Expecter) CheckUserExistence(id interface{}) *UserChecker_CheckUserExistence_Call {
	return &UserChecker_CheckUserExistence_Call{Call: _e.mock.On("CheckUserExistence", id)}
}

func (_c *UserChecker_CheckUserExistence_Call) Run(run func(id int64)) *UserChecker_CheckUserExistence_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *UserChecker_CheckUserExistence_Call) Return(_a0 bool) *UserChecker_CheckUserExistence_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserChecker_CheckUserExistence_Call) RunAndReturn(run func(int64) bool) *UserChecker_CheckUserExistence_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserChecker creates a new instance of UserChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserChecker(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserChecker {
	mock := &UserChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
