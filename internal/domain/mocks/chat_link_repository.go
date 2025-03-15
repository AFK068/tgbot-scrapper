// Code generated by mockery v2.52.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/AFK068/bot/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// ChatLinkRepository is an autogenerated mock type for the ChatLinkRepository type
type ChatLinkRepository struct {
	mock.Mock
}

type ChatLinkRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *ChatLinkRepository) EXPECT() *ChatLinkRepository_Expecter {
	return &ChatLinkRepository_Expecter{mock: &_m.Mock}
}

// CheckUserExistence provides a mock function with given fields: ctx, chatID
func (_m *ChatLinkRepository) CheckUserExistence(ctx context.Context, chatID int64) (bool, error) {
	ret := _m.Called(ctx, chatID)

	if len(ret) == 0 {
		panic("no return value specified for CheckUserExistence")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (bool, error)); ok {
		return rf(ctx, chatID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) bool); ok {
		r0 = rf(ctx, chatID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, chatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChatLinkRepository_CheckUserExistence_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckUserExistence'
type ChatLinkRepository_CheckUserExistence_Call struct {
	*mock.Call
}

// CheckUserExistence is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
func (_e *ChatLinkRepository_Expecter) CheckUserExistence(ctx interface{}, chatID interface{}) *ChatLinkRepository_CheckUserExistence_Call {
	return &ChatLinkRepository_CheckUserExistence_Call{Call: _e.mock.On("CheckUserExistence", ctx, chatID)}
}

func (_c *ChatLinkRepository_CheckUserExistence_Call) Run(run func(ctx context.Context, chatID int64)) *ChatLinkRepository_CheckUserExistence_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *ChatLinkRepository_CheckUserExistence_Call) Return(_a0 bool, _a1 error) *ChatLinkRepository_CheckUserExistence_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChatLinkRepository_CheckUserExistence_Call) RunAndReturn(run func(context.Context, int64) (bool, error)) *ChatLinkRepository_CheckUserExistence_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteChat provides a mock function with given fields: ctx, chatID
func (_m *ChatLinkRepository) DeleteChat(ctx context.Context, chatID int64) error {
	ret := _m.Called(ctx, chatID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteChat")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, chatID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ChatLinkRepository_DeleteChat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteChat'
type ChatLinkRepository_DeleteChat_Call struct {
	*mock.Call
}

// DeleteChat is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
func (_e *ChatLinkRepository_Expecter) DeleteChat(ctx interface{}, chatID interface{}) *ChatLinkRepository_DeleteChat_Call {
	return &ChatLinkRepository_DeleteChat_Call{Call: _e.mock.On("DeleteChat", ctx, chatID)}
}

func (_c *ChatLinkRepository_DeleteChat_Call) Run(run func(ctx context.Context, chatID int64)) *ChatLinkRepository_DeleteChat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *ChatLinkRepository_DeleteChat_Call) Return(_a0 error) *ChatLinkRepository_DeleteChat_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ChatLinkRepository_DeleteChat_Call) RunAndReturn(run func(context.Context, int64) error) *ChatLinkRepository_DeleteChat_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteLink provides a mock function with given fields: ctx, uid, link
func (_m *ChatLinkRepository) DeleteLink(ctx context.Context, uid int64, link *domain.Link) error {
	ret := _m.Called(ctx, uid, link)

	if len(ret) == 0 {
		panic("no return value specified for DeleteLink")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *domain.Link) error); ok {
		r0 = rf(ctx, uid, link)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ChatLinkRepository_DeleteLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteLink'
type ChatLinkRepository_DeleteLink_Call struct {
	*mock.Call
}

// DeleteLink is a helper method to define mock.On call
//   - ctx context.Context
//   - uid int64
//   - link *domain.Link
func (_e *ChatLinkRepository_Expecter) DeleteLink(ctx interface{}, uid interface{}, link interface{}) *ChatLinkRepository_DeleteLink_Call {
	return &ChatLinkRepository_DeleteLink_Call{Call: _e.mock.On("DeleteLink", ctx, uid, link)}
}

func (_c *ChatLinkRepository_DeleteLink_Call) Run(run func(ctx context.Context, uid int64, link *domain.Link)) *ChatLinkRepository_DeleteLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(*domain.Link))
	})
	return _c
}

func (_c *ChatLinkRepository_DeleteLink_Call) Return(_a0 error) *ChatLinkRepository_DeleteLink_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ChatLinkRepository_DeleteLink_Call) RunAndReturn(run func(context.Context, int64, *domain.Link) error) *ChatLinkRepository_DeleteLink_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllLinks provides a mock function with given fields: ctx
func (_m *ChatLinkRepository) GetAllLinks(ctx context.Context) ([]*domain.Link, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllLinks")
	}

	var r0 []*domain.Link
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*domain.Link, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*domain.Link); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Link)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChatLinkRepository_GetAllLinks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllLinks'
type ChatLinkRepository_GetAllLinks_Call struct {
	*mock.Call
}

// GetAllLinks is a helper method to define mock.On call
//   - ctx context.Context
func (_e *ChatLinkRepository_Expecter) GetAllLinks(ctx interface{}) *ChatLinkRepository_GetAllLinks_Call {
	return &ChatLinkRepository_GetAllLinks_Call{Call: _e.mock.On("GetAllLinks", ctx)}
}

func (_c *ChatLinkRepository_GetAllLinks_Call) Run(run func(ctx context.Context)) *ChatLinkRepository_GetAllLinks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *ChatLinkRepository_GetAllLinks_Call) Return(_a0 []*domain.Link, _a1 error) *ChatLinkRepository_GetAllLinks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChatLinkRepository_GetAllLinks_Call) RunAndReturn(run func(context.Context) ([]*domain.Link, error)) *ChatLinkRepository_GetAllLinks_Call {
	_c.Call.Return(run)
	return _c
}

// GetChatIDsByLink provides a mock function with given fields: ctx, link
func (_m *ChatLinkRepository) GetChatIDsByLink(ctx context.Context, link *domain.Link) ([]int64, error) {
	ret := _m.Called(ctx, link)

	if len(ret) == 0 {
		panic("no return value specified for GetChatIDsByLink")
	}

	var r0 []int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Link) ([]int64, error)); ok {
		return rf(ctx, link)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Link) []int64); ok {
		r0 = rf(ctx, link)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Link) error); ok {
		r1 = rf(ctx, link)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChatLinkRepository_GetChatIDsByLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChatIDsByLink'
type ChatLinkRepository_GetChatIDsByLink_Call struct {
	*mock.Call
}

// GetChatIDsByLink is a helper method to define mock.On call
//   - ctx context.Context
//   - link *domain.Link
func (_e *ChatLinkRepository_Expecter) GetChatIDsByLink(ctx interface{}, link interface{}) *ChatLinkRepository_GetChatIDsByLink_Call {
	return &ChatLinkRepository_GetChatIDsByLink_Call{Call: _e.mock.On("GetChatIDsByLink", ctx, link)}
}

func (_c *ChatLinkRepository_GetChatIDsByLink_Call) Run(run func(ctx context.Context, link *domain.Link)) *ChatLinkRepository_GetChatIDsByLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Link))
	})
	return _c
}

func (_c *ChatLinkRepository_GetChatIDsByLink_Call) Return(_a0 []int64, _a1 error) *ChatLinkRepository_GetChatIDsByLink_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChatLinkRepository_GetChatIDsByLink_Call) RunAndReturn(run func(context.Context, *domain.Link) ([]int64, error)) *ChatLinkRepository_GetChatIDsByLink_Call {
	_c.Call.Return(run)
	return _c
}

// GetListLinks provides a mock function with given fields: ctx, chatID
func (_m *ChatLinkRepository) GetListLinks(ctx context.Context, chatID int64) ([]*domain.Link, error) {
	ret := _m.Called(ctx, chatID)

	if len(ret) == 0 {
		panic("no return value specified for GetListLinks")
	}

	var r0 []*domain.Link
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]*domain.Link, error)); ok {
		return rf(ctx, chatID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*domain.Link); ok {
		r0 = rf(ctx, chatID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Link)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, chatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChatLinkRepository_GetListLinks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetListLinks'
type ChatLinkRepository_GetListLinks_Call struct {
	*mock.Call
}

// GetListLinks is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
func (_e *ChatLinkRepository_Expecter) GetListLinks(ctx interface{}, chatID interface{}) *ChatLinkRepository_GetListLinks_Call {
	return &ChatLinkRepository_GetListLinks_Call{Call: _e.mock.On("GetListLinks", ctx, chatID)}
}

func (_c *ChatLinkRepository_GetListLinks_Call) Run(run func(ctx context.Context, chatID int64)) *ChatLinkRepository_GetListLinks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *ChatLinkRepository_GetListLinks_Call) Return(_a0 []*domain.Link, _a1 error) *ChatLinkRepository_GetListLinks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChatLinkRepository_GetListLinks_Call) RunAndReturn(run func(context.Context, int64) ([]*domain.Link, error)) *ChatLinkRepository_GetListLinks_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterChat provides a mock function with given fields: ctx, chatID
func (_m *ChatLinkRepository) RegisterChat(ctx context.Context, chatID int64) error {
	ret := _m.Called(ctx, chatID)

	if len(ret) == 0 {
		panic("no return value specified for RegisterChat")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, chatID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ChatLinkRepository_RegisterChat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterChat'
type ChatLinkRepository_RegisterChat_Call struct {
	*mock.Call
}

// RegisterChat is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
func (_e *ChatLinkRepository_Expecter) RegisterChat(ctx interface{}, chatID interface{}) *ChatLinkRepository_RegisterChat_Call {
	return &ChatLinkRepository_RegisterChat_Call{Call: _e.mock.On("RegisterChat", ctx, chatID)}
}

func (_c *ChatLinkRepository_RegisterChat_Call) Run(run func(ctx context.Context, chatID int64)) *ChatLinkRepository_RegisterChat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *ChatLinkRepository_RegisterChat_Call) Return(_a0 error) *ChatLinkRepository_RegisterChat_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ChatLinkRepository_RegisterChat_Call) RunAndReturn(run func(context.Context, int64) error) *ChatLinkRepository_RegisterChat_Call {
	_c.Call.Return(run)
	return _c
}

// SaveLink provides a mock function with given fields: ctx, uid, link
func (_m *ChatLinkRepository) SaveLink(ctx context.Context, uid int64, link *domain.Link) error {
	ret := _m.Called(ctx, uid, link)

	if len(ret) == 0 {
		panic("no return value specified for SaveLink")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *domain.Link) error); ok {
		r0 = rf(ctx, uid, link)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ChatLinkRepository_SaveLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveLink'
type ChatLinkRepository_SaveLink_Call struct {
	*mock.Call
}

// SaveLink is a helper method to define mock.On call
//   - ctx context.Context
//   - uid int64
//   - link *domain.Link
func (_e *ChatLinkRepository_Expecter) SaveLink(ctx interface{}, uid interface{}, link interface{}) *ChatLinkRepository_SaveLink_Call {
	return &ChatLinkRepository_SaveLink_Call{Call: _e.mock.On("SaveLink", ctx, uid, link)}
}

func (_c *ChatLinkRepository_SaveLink_Call) Run(run func(ctx context.Context, uid int64, link *domain.Link)) *ChatLinkRepository_SaveLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(*domain.Link))
	})
	return _c
}

func (_c *ChatLinkRepository_SaveLink_Call) Return(_a0 error) *ChatLinkRepository_SaveLink_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ChatLinkRepository_SaveLink_Call) RunAndReturn(run func(context.Context, int64, *domain.Link) error) *ChatLinkRepository_SaveLink_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateLastCheck provides a mock function with given fields: ctx, link
func (_m *ChatLinkRepository) UpdateLastCheck(ctx context.Context, link *domain.Link) error {
	ret := _m.Called(ctx, link)

	if len(ret) == 0 {
		panic("no return value specified for UpdateLastCheck")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Link) error); ok {
		r0 = rf(ctx, link)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ChatLinkRepository_UpdateLastCheck_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateLastCheck'
type ChatLinkRepository_UpdateLastCheck_Call struct {
	*mock.Call
}

// UpdateLastCheck is a helper method to define mock.On call
//   - ctx context.Context
//   - link *domain.Link
func (_e *ChatLinkRepository_Expecter) UpdateLastCheck(ctx interface{}, link interface{}) *ChatLinkRepository_UpdateLastCheck_Call {
	return &ChatLinkRepository_UpdateLastCheck_Call{Call: _e.mock.On("UpdateLastCheck", ctx, link)}
}

func (_c *ChatLinkRepository_UpdateLastCheck_Call) Run(run func(ctx context.Context, link *domain.Link)) *ChatLinkRepository_UpdateLastCheck_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Link))
	})
	return _c
}

func (_c *ChatLinkRepository_UpdateLastCheck_Call) Return(_a0 error) *ChatLinkRepository_UpdateLastCheck_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ChatLinkRepository_UpdateLastCheck_Call) RunAndReturn(run func(context.Context, *domain.Link) error) *ChatLinkRepository_UpdateLastCheck_Call {
	_c.Call.Return(run)
	return _c
}

// NewChatLinkRepository creates a new instance of ChatLinkRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChatLinkRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ChatLinkRepository {
	mock := &ChatLinkRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
