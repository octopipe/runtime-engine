// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Call provides a mock function with given fields: method, args, reply
func (_m *Client) Call(method string, args interface{}, reply interface{}) error {
	ret := _m.Called(method, args, reply)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) error); ok {
		r0 = rf(method, args, reply)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
