// Code generated by mockery v2.43.0. DO NOT EDIT.

package wrapper

import (
	mock "github.com/stretchr/testify/mock"
	ssh "gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
)

// mockDialerFactory is an autogenerated mock type for the dialerFactory type
type mockDialerFactory struct {
	mock.Mock
}

type mockDialerFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDialerFactory) EXPECT() *mockDialerFactory_Expecter {
	return &mockDialerFactory_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: flags, def
func (_m *mockDialerFactory) Create(flags ssh.Flags, def ssh.TargetDef) (ssh.Dialer, error) {
	ret := _m.Called(flags, def)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 ssh.Dialer
	var r1 error
	if rf, ok := ret.Get(0).(func(ssh.Flags, ssh.TargetDef) (ssh.Dialer, error)); ok {
		return rf(flags, def)
	}
	if rf, ok := ret.Get(0).(func(ssh.Flags, ssh.TargetDef) ssh.Dialer); ok {
		r0 = rf(flags, def)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ssh.Dialer)
		}
	}

	if rf, ok := ret.Get(1).(func(ssh.Flags, ssh.TargetDef) error); ok {
		r1 = rf(flags, def)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDialerFactory_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockDialerFactory_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - flags ssh.Flags
//   - def ssh.TargetDef
func (_e *mockDialerFactory_Expecter) Create(flags interface{}, def interface{}) *mockDialerFactory_Create_Call {
	return &mockDialerFactory_Create_Call{Call: _e.mock.On("Create", flags, def)}
}

func (_c *mockDialerFactory_Create_Call) Run(run func(flags ssh.Flags, def ssh.TargetDef)) *mockDialerFactory_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ssh.Flags), args[1].(ssh.TargetDef))
	})
	return _c
}

func (_c *mockDialerFactory_Create_Call) Return(_a0 ssh.Dialer, _a1 error) *mockDialerFactory_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDialerFactory_Create_Call) RunAndReturn(run func(ssh.Flags, ssh.TargetDef) (ssh.Dialer, error)) *mockDialerFactory_Create_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDialerFactory creates a new instance of mockDialerFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDialerFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDialerFactory {
	mock := &mockDialerFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
