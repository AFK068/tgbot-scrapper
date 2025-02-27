// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

type Service_Expecter struct {
	mock *mock.Mock
}

func (_m *Service) EXPECT() *Service_Expecter {
	return &Service_Expecter{mock: &_m.Mock}
}

// Run provides a mock function with given fields:
func (_m *Service) Run() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Run")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Service_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type Service_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
func (_e *Service_Expecter) Run() *Service_Run_Call {
	return &Service_Run_Call{Call: _e.mock.On("Run")}
}

func (_c *Service_Run_Call) Run(run func()) *Service_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Service_Run_Call) Return(_a0 error) *Service_Run_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Service_Run_Call) RunAndReturn(run func() error) *Service_Run_Call {
	_c.Call.Return(run)
	return _c
}

// SendMessage provides a mock function with given fields: chatID, text
func (_m *Service) SendMessage(chatID int64, text string) {
	_m.Called(chatID, text)
}

// Service_SendMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMessage'
type Service_SendMessage_Call struct {
	*mock.Call
}

// SendMessage is a helper method to define mock.On call
//   - chatID int64
//   - text string
func (_e *Service_Expecter) SendMessage(chatID interface{}, text interface{}) *Service_SendMessage_Call {
	return &Service_SendMessage_Call{Call: _e.mock.On("SendMessage", chatID, text)}
}

func (_c *Service_SendMessage_Call) Run(run func(chatID int64, text string)) *Service_SendMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(string))
	})
	return _c
}

func (_c *Service_SendMessage_Call) Return() *Service_SendMessage_Call {
	_c.Call.Return()
	return _c
}

func (_c *Service_SendMessage_Call) RunAndReturn(run func(int64, string)) *Service_SendMessage_Call {
	_c.Call.Return(run)
	return _c
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
