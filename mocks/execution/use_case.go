// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	execution "github.com/octopipe/cloudx/internal/execution"
	mock "github.com/stretchr/testify/mock"

	pagination "github.com/octopipe/cloudx/internal/pagination"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *UseCase) Create(ctx context.Context, _a1 execution.Execution) (execution.Execution, error) {
	ret := _m.Called(ctx, _a1)

	var r0 execution.Execution
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, execution.Execution) (execution.Execution, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, execution.Execution) execution.Execution); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(execution.Execution)
	}

	if rf, ok := ret.Get(1).(func(context.Context, execution.Execution) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, name, namespace
func (_m *UseCase) Delete(ctx context.Context, name string, namespace string) error {
	ret := _m.Called(ctx, name, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, name, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, name, namespace
func (_m *UseCase) Get(ctx context.Context, name string, namespace string) (execution.Execution, error) {
	ret := _m.Called(ctx, name, namespace)

	var r0 execution.Execution
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (execution.Execution, error)); ok {
		return rf(ctx, name, namespace)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) execution.Execution); ok {
		r0 = rf(ctx, name, namespace)
	} else {
		r0 = ret.Get(0).(execution.Execution)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, name, namespace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, namespace, chunkPagination
func (_m *UseCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[execution.Execution], error) {
	ret := _m.Called(ctx, namespace, chunkPagination)

	var r0 pagination.ChunkingPaginationResponse[execution.Execution]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[execution.Execution], error)); ok {
		return rf(ctx, namespace, chunkPagination)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, pagination.ChunkingPaginationRequest) pagination.ChunkingPaginationResponse[execution.Execution]); ok {
		r0 = rf(ctx, namespace, chunkPagination)
	} else {
		r0 = ret.Get(0).(pagination.ChunkingPaginationResponse[execution.Execution])
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, pagination.ChunkingPaginationRequest) error); ok {
		r1 = rf(ctx, namespace, chunkPagination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *UseCase) Update(ctx context.Context, _a1 execution.Execution) (execution.Execution, error) {
	ret := _m.Called(ctx, _a1)

	var r0 execution.Execution
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, execution.Execution) (execution.Execution, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, execution.Execution) execution.Execution); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(execution.Execution)
	}

	if rf, ok := ret.Get(1).(func(context.Context, execution.Execution) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}