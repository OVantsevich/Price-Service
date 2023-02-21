package service

import (
	"context"
	"testing"

	"Price-Service/internal/model"
	"Price-Service/internal/service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPrices_All(t *testing.T) {
	end := make(chan struct{}, 1)
	mq := mocks.NewMQ(t)
	sp := mocks.NewStreamPoolRepository(t)
	sp.On("Send", mock.AnythingOfType("[]*model.Price")).Maybe()
	sp.On("Update", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("chan *model.Price"),
		mock.AnythingOfType("[]string")).Maybe()
	sp.On("Delete", mock.AnythingOfType("uuid.UUID")).Maybe()

	mq.On("GetPrices", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return([][]*model.Price{}, "", 0, nil).Maybe()

	prices = NewPrices(context.Background(), sp, mq, "0", end)

	testStreamID := uuid.New()
	prices.Subscribe(testStreamID)

	err := prices.DeleteSubscription(testStreamID)
	require.NoError(t, err)

	err = prices.UpdateSubscription([]string{}, testStreamID)
	require.Error(t, err)

	err = prices.DeleteSubscription(testStreamID)
	require.Error(t, err)
}
