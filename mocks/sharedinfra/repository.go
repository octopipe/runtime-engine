// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	pagination "github.com/octopipe/cloudx/internal/pagination"
	mock "github.com/stretchr/testify/mock"

	v1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Apply provides a mock function with given fields: ctx, s
func (_m *Repository) Apply(ctx context.Context, s v1alpha1.Infra) (v1alpha1.Infra, error) {
	ret := _m.Called(ctx, s)

	var r0 v1alpha1.Infra
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, v1alpha1.Infra) (v1alpha1.Infra, error)); ok {
		return rf(ctx, s)
	}
	if rf, ok := ret.Get(0).(func(context.Context, v1alpha1.Infra) v1alpha1.Infra); ok {
		r0 = rf(ctx, s)
	} else {
		r0 = ret.Get(0).(v1alpha1.Infra)
	}

	if rf, ok := ret.Get(1).(func(context.Context, v1alpha1.Infra) error); ok {
		r1 = rf(ctx, s)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, name, namespace
func (_m *Repository) Delete(ctx context.Context, name string, namespace string) error {
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
func (_m *Repository) Get(ctx context.Context, name string, namespace string) (v1alpha1.Infra, error) {
	ret := _m.Called(ctx, name, namespace)

	var r0 v1alpha1.Infra
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (v1alpha1.Infra, error)); ok {
		return rf(ctx, name, namespace)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) v1alpha1.Infra); ok {
		r0 = rf(ctx, name, namespace)
	} else {
		r0 = ret.Get(0).(v1alpha1.Infra)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, name, namespace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, namespace, chunkPagination
func (_m *Repository) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (v1alpha1.InfraList, error) {
	ret := _m.Called(ctx, namespace, chunkPagination)

	var r0 v1alpha1.InfraList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, pagination.ChunkingPaginationRequest) (v1alpha1.InfraList, error)); ok {
		return rf(ctx, namespace, chunkPagination)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, pagination.ChunkingPaginationRequest) v1alpha1.InfraList); ok {
		r0 = rf(ctx, namespace, chunkPagination)
	} else {
		r0 = ret.Get(0).(v1alpha1.InfraList)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, pagination.ChunkingPaginationRequest) error); ok {
		r1 = rf(ctx, namespace, chunkPagination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Reconcile provides a mock function with given fields: ctx, name, namespace
func (_m *Repository) Reconcile(ctx context.Context, name string, namespace string) error {
	ret := _m.Called(ctx, name, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, name, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
