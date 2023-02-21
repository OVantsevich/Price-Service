// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/OVantsevich/Price-Service/internal/model"

	uuid "github.com/google/uuid"
)

// PriceService is an autogenerated mock type for the PriceService type
type PriceService struct {
	mock.Mock
}

// DeleteSubscription provides a mock function with given fields: streamID
func (_m *PriceService) DeleteSubscription(streamID uuid.UUID) error {
	ret := _m.Called(streamID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(streamID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCurrentPrices provides a mock function with given fields: ctx, names
func (_m *PriceService) GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error) {
	ret := _m.Called(ctx, names)

	var r0 map[string]*model.Price
	if rf, ok := ret.Get(0).(func(context.Context, []string) map[string]*model.Price); ok {
		r0 = rf(ctx, names)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*model.Price)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, names)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subscribe provides a mock function with given fields: streamID
func (_m *PriceService) Subscribe(streamID uuid.UUID) chan *model.Price {
	ret := _m.Called(streamID)

	var r0 chan *model.Price
	if rf, ok := ret.Get(0).(func(uuid.UUID) chan *model.Price); ok {
		r0 = rf(streamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *model.Price)
		}
	}

	return r0
}

// UpdateSubscription provides a mock function with given fields: names, streamID
func (_m *PriceService) UpdateSubscription(names []string, streamID uuid.UUID) error {
	ret := _m.Called(names, streamID)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, uuid.UUID) error); ok {
		r0 = rf(names, streamID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPriceService interface {
	mock.TestingT
	Cleanup(func())
}

// NewPriceService creates a new instance of PriceService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPriceService(t mockConstructorTestingTNewPriceService) *PriceService {
	mock := &PriceService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
