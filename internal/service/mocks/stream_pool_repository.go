// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	model "github.com/OVantsevich/Price-Service/internal/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// StreamPoolRepository is an autogenerated mock type for the StreamPoolRepository type
type StreamPoolRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: streamID
func (_m *StreamPoolRepository) Delete(streamID uuid.UUID) {
	_m.Called(streamID)
}

// Send provides a mock function with given fields: prices
func (_m *StreamPoolRepository) Send(prices []*model.Price) {
	_m.Called(prices)
}

// Update provides a mock function with given fields: streamID, streamChan, prices
func (_m *StreamPoolRepository) Update(streamID uuid.UUID, streamChan chan *model.Price, prices []string) {
	_m.Called(streamID, streamChan, prices)
}

type mockConstructorTestingTNewStreamPoolRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewStreamPoolRepository creates a new instance of StreamPoolRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStreamPoolRepository(t mockConstructorTestingTNewStreamPoolRepository) *StreamPoolRepository {
	mock := &StreamPoolRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
