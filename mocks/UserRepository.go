// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	model "pinstack-user-service/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

type UserRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepository) EXPECT() *UserRepository_Expecter {
	return &UserRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, user
func (_m *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) (*model.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) *model.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - user *model.User
func (_e *UserRepository_Expecter) Create(ctx interface{}, user interface{}) *UserRepository_Create_Call {
	return &UserRepository_Create_Call{Call: _e.mock.On("Create", ctx, user)}
}

func (_c *UserRepository_Create_Call) Run(run func(ctx context.Context, user *model.User)) *UserRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.User))
	})
	return _c
}

func (_c *UserRepository_Create_Call) Return(_a0 *model.User, _a1 error) *UserRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Create_Call) RunAndReturn(run func(context.Context, *model.User) (*model.User, error)) *UserRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *UserRepository) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type UserRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *UserRepository_Expecter) Delete(ctx interface{}, id interface{}) *UserRepository_Delete_Call {
	return &UserRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *UserRepository_Delete_Call) Run(run func(ctx context.Context, id int64)) *UserRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *UserRepository_Delete_Call) Return(_a0 error) *UserRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_Delete_Call) RunAndReturn(run func(context.Context, int64) error) *UserRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetByEmail")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type UserRepository_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *UserRepository_Expecter) GetByEmail(ctx interface{}, email interface{}) *UserRepository_GetByEmail_Call {
	return &UserRepository_GetByEmail_Call{Call: _e.mock.On("GetByEmail", ctx, email)}
}

func (_c *UserRepository_GetByEmail_Call) Run(run func(ctx context.Context, email string)) *UserRepository_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_GetByEmail_Call) Return(_a0 *model.User, _a1 error) *UserRepository_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetByEmail_Call) RunAndReturn(run func(context.Context, string) (*model.User, error)) *UserRepository_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*model.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *model.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type UserRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *UserRepository_Expecter) GetByID(ctx interface{}, id interface{}) *UserRepository_GetByID_Call {
	return &UserRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *UserRepository_GetByID_Call) Run(run func(ctx context.Context, id int64)) *UserRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *UserRepository_GetByID_Call) Return(_a0 *model.User, _a1 error) *UserRepository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetByID_Call) RunAndReturn(run func(context.Context, int64) (*model.User, error)) *UserRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUsername provides a mock function with given fields: ctx, username
func (_m *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for GetByUsername")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.User, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUsername'
type UserRepository_GetByUsername_Call struct {
	*mock.Call
}

// GetByUsername is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *UserRepository_Expecter) GetByUsername(ctx interface{}, username interface{}) *UserRepository_GetByUsername_Call {
	return &UserRepository_GetByUsername_Call{Call: _e.mock.On("GetByUsername", ctx, username)}
}

func (_c *UserRepository_GetByUsername_Call) Run(run func(ctx context.Context, username string)) *UserRepository_GetByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_GetByUsername_Call) Return(_a0 *model.User, _a1 error) *UserRepository_GetByUsername_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetByUsername_Call) RunAndReturn(run func(context.Context, string) (*model.User, error)) *UserRepository_GetByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// Search provides a mock function with given fields: ctx, searchQuery, offset, pageSize
func (_m *UserRepository) Search(ctx context.Context, searchQuery string, offset int, pageSize int) ([]*model.User, int, error) {
	ret := _m.Called(ctx, searchQuery, offset, pageSize)

	if len(ret) == 0 {
		panic("no return value specified for Search")
	}

	var r0 []*model.User
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]*model.User, int, error)); ok {
		return rf(ctx, searchQuery, offset, pageSize)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*model.User); ok {
		r0 = rf(ctx, searchQuery, offset, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) int); ok {
		r1 = rf(ctx, searchQuery, offset, pageSize)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int, int) error); ok {
		r2 = rf(ctx, searchQuery, offset, pageSize)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UserRepository_Search_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Search'
type UserRepository_Search_Call struct {
	*mock.Call
}

// Search is a helper method to define mock.On call
//   - ctx context.Context
//   - searchQuery string
//   - offset int
//   - pageSize int
func (_e *UserRepository_Expecter) Search(ctx interface{}, searchQuery interface{}, offset interface{}, pageSize interface{}) *UserRepository_Search_Call {
	return &UserRepository_Search_Call{Call: _e.mock.On("Search", ctx, searchQuery, offset, pageSize)}
}

func (_c *UserRepository_Search_Call) Run(run func(ctx context.Context, searchQuery string, offset int, pageSize int)) *UserRepository_Search_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *UserRepository_Search_Call) Return(_a0 []*model.User, _a1 int, _a2 error) *UserRepository_Search_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *UserRepository_Search_Call) RunAndReturn(run func(context.Context, string, int, int) ([]*model.User, int, error)) *UserRepository_Search_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, user
func (_m *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) (*model.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) *model.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type UserRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - user *model.User
func (_e *UserRepository_Expecter) Update(ctx interface{}, user interface{}) *UserRepository_Update_Call {
	return &UserRepository_Update_Call{Call: _e.mock.On("Update", ctx, user)}
}

func (_c *UserRepository_Update_Call) Run(run func(ctx context.Context, user *model.User)) *UserRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.User))
	})
	return _c
}

func (_c *UserRepository_Update_Call) Return(_a0 *model.User, _a1 error) *UserRepository_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Update_Call) RunAndReturn(run func(context.Context, *model.User) (*model.User, error)) *UserRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAvatar provides a mock function with given fields: ctx, id, avatarURL
func (_m *UserRepository) UpdateAvatar(ctx context.Context, id int64, avatarURL string) error {
	ret := _m.Called(ctx, id, avatarURL)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAvatar")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, avatarURL)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_UpdateAvatar_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAvatar'
type UserRepository_UpdateAvatar_Call struct {
	*mock.Call
}

// UpdateAvatar is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - avatarURL string
func (_e *UserRepository_Expecter) UpdateAvatar(ctx interface{}, id interface{}, avatarURL interface{}) *UserRepository_UpdateAvatar_Call {
	return &UserRepository_UpdateAvatar_Call{Call: _e.mock.On("UpdateAvatar", ctx, id, avatarURL)}
}

func (_c *UserRepository_UpdateAvatar_Call) Run(run func(ctx context.Context, id int64, avatarURL string)) *UserRepository_UpdateAvatar_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *UserRepository_UpdateAvatar_Call) Return(_a0 error) *UserRepository_UpdateAvatar_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_UpdateAvatar_Call) RunAndReturn(run func(context.Context, int64, string) error) *UserRepository_UpdateAvatar_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePassword provides a mock function with given fields: ctx, id, newPassword
func (_m *UserRepository) UpdatePassword(ctx context.Context, id int64, newPassword string) error {
	ret := _m.Called(ctx, id, newPassword)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, newPassword)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_UpdatePassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePassword'
type UserRepository_UpdatePassword_Call struct {
	*mock.Call
}

// UpdatePassword is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - newPassword string
func (_e *UserRepository_Expecter) UpdatePassword(ctx interface{}, id interface{}, newPassword interface{}) *UserRepository_UpdatePassword_Call {
	return &UserRepository_UpdatePassword_Call{Call: _e.mock.On("UpdatePassword", ctx, id, newPassword)}
}

func (_c *UserRepository_UpdatePassword_Call) Run(run func(ctx context.Context, id int64, newPassword string)) *UserRepository_UpdatePassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *UserRepository_UpdatePassword_Call) Return(_a0 error) *UserRepository_UpdatePassword_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_UpdatePassword_Call) RunAndReturn(run func(context.Context, int64, string) error) *UserRepository_UpdatePassword_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
