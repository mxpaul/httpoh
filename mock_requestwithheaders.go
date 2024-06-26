// Code generated by mockery. DO NOT EDIT.

package httpoh

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockRequestWithHeaders is an autogenerated mock type for the RequestWithHeaders type
type MockRequestWithHeaders struct {
	mock.Mock
}

type MockRequestWithHeaders_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRequestWithHeaders) EXPECT() *MockRequestWithHeaders_Expecter {
	return &MockRequestWithHeaders_Expecter{mock: &_m.Mock}
}

// Headers provides a mock function with given fields:
func (_m *MockRequestWithHeaders) Headers() http.Header {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Headers")
	}

	var r0 http.Header
	if rf, ok := ret.Get(0).(func() http.Header); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Header)
		}
	}

	return r0
}

// MockRequestWithHeaders_Headers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Headers'
type MockRequestWithHeaders_Headers_Call struct {
	*mock.Call
}

// Headers is a helper method to define mock.On call
func (_e *MockRequestWithHeaders_Expecter) Headers() *MockRequestWithHeaders_Headers_Call {
	return &MockRequestWithHeaders_Headers_Call{Call: _e.mock.On("Headers")}
}

func (_c *MockRequestWithHeaders_Headers_Call) Run(run func()) *MockRequestWithHeaders_Headers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRequestWithHeaders_Headers_Call) Return(_a0 http.Header) *MockRequestWithHeaders_Headers_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRequestWithHeaders_Headers_Call) RunAndReturn(run func() http.Header) *MockRequestWithHeaders_Headers_Call {
	_c.Call.Return(run)
	return _c
}

// Method provides a mock function with given fields:
func (_m *MockRequestWithHeaders) Method() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Method")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockRequestWithHeaders_Method_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Method'
type MockRequestWithHeaders_Method_Call struct {
	*mock.Call
}

// Method is a helper method to define mock.On call
func (_e *MockRequestWithHeaders_Expecter) Method() *MockRequestWithHeaders_Method_Call {
	return &MockRequestWithHeaders_Method_Call{Call: _e.mock.On("Method")}
}

func (_c *MockRequestWithHeaders_Method_Call) Run(run func()) *MockRequestWithHeaders_Method_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRequestWithHeaders_Method_Call) Return(_a0 string) *MockRequestWithHeaders_Method_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRequestWithHeaders_Method_Call) RunAndReturn(run func() string) *MockRequestWithHeaders_Method_Call {
	_c.Call.Return(run)
	return _c
}

// URL provides a mock function with given fields:
func (_m *MockRequestWithHeaders) URL() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for URL")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockRequestWithHeaders_URL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'URL'
type MockRequestWithHeaders_URL_Call struct {
	*mock.Call
}

// URL is a helper method to define mock.On call
func (_e *MockRequestWithHeaders_Expecter) URL() *MockRequestWithHeaders_URL_Call {
	return &MockRequestWithHeaders_URL_Call{Call: _e.mock.On("URL")}
}

func (_c *MockRequestWithHeaders_URL_Call) Run(run func()) *MockRequestWithHeaders_URL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRequestWithHeaders_URL_Call) Return(_a0 string) *MockRequestWithHeaders_URL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRequestWithHeaders_URL_Call) RunAndReturn(run func() string) *MockRequestWithHeaders_URL_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRequestWithHeaders creates a new instance of MockRequestWithHeaders. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRequestWithHeaders(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRequestWithHeaders {
	mock := &MockRequestWithHeaders{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
