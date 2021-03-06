// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/duyquang6/wager-management-be/pkg/dto"
	mock "github.com/stretchr/testify/mock"
)

// WagerService is an autogenerated mock type for the WagerService type
type WagerService struct {
	mock.Mock
}

// CreateWager provides a mock function with given fields: ctx, request
func (_m *WagerService) CreateWager(ctx context.Context, request dto.CreateWagerRequest) (dto.CreateWagerResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 dto.CreateWagerResponse
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreateWagerRequest) dto.CreateWagerResponse); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(dto.CreateWagerResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.CreateWagerRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListWagers provides a mock function with given fields: ctx, page, limit
func (_m *WagerService) ListWagers(ctx context.Context, page uint, limit uint) (dto.ListWagersResponse, error) {
	ret := _m.Called(ctx, page, limit)

	var r0 dto.ListWagersResponse
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) dto.ListWagersResponse); ok {
		r0 = rf(ctx, page, limit)
	} else {
		r0 = ret.Get(0).(dto.ListWagersResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) error); ok {
		r1 = rf(ctx, page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
