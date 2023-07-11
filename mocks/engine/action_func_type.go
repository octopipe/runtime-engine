// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	engine "github.com/octopipe/cloudx/internal/engine"
	mock "github.com/stretchr/testify/mock"

	v1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
)

// ActionFuncType is an autogenerated mock type for the ActionFuncType type
type ActionFuncType struct {
	mock.Mock
}

// Execute provides a mock function with given fields: taskName, executionOutputs
func (_m *ActionFuncType) Execute(taskName string, executionOutputs engine.ExecutionContext) (v1alpha1.TaskExecutionStatus, map[string]engine.ExecutionOutputItem) {
	ret := _m.Called(taskName, executionOutputs)

	var r0 v1alpha1.TaskExecutionStatus
	var r1 map[string]engine.ExecutionOutputItem
	if rf, ok := ret.Get(0).(func(string, engine.ExecutionContext) (v1alpha1.TaskExecutionStatus, map[string]engine.ExecutionOutputItem)); ok {
		return rf(taskName, executionOutputs)
	}
	if rf, ok := ret.Get(0).(func(string, engine.ExecutionContext) v1alpha1.TaskExecutionStatus); ok {
		r0 = rf(taskName, executionOutputs)
	} else {
		r0 = ret.Get(0).(v1alpha1.TaskExecutionStatus)
	}

	if rf, ok := ret.Get(1).(func(string, engine.ExecutionContext) map[string]engine.ExecutionOutputItem); ok {
		r1 = rf(taskName, executionOutputs)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]engine.ExecutionOutputItem)
		}
	}

	return r0, r1
}

// NewActionFuncType creates a new instance of ActionFuncType. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewActionFuncType(t interface {
	mock.TestingT
	Cleanup(func())
}) *ActionFuncType {
	mock := &ActionFuncType{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
