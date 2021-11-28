// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/duyquang6/wager-management-be/pkg/dto"
	mock "github.com/stretchr/testify/mock"
)

// PurchaseService is an autogenerated mock type for the PurchaseService type
type PurchaseService struct {
	mock.Mock
}

// CreatePurchase provides a mock function with given fields: ctx, request
func (_m *PurchaseService) CreatePurchase(ctx context.Context, request dto.CreatePurchaseRequest) (dto.CreatePurchaseResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 dto.CreatePurchaseResponse
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreatePurchaseRequest) dto.CreatePurchaseResponse); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(dto.CreatePurchaseResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.CreatePurchaseRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
