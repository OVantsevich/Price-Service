package service

import (
	"Price-Service/internal/service/mocks"
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"

	"github.com/google/uuid"
)

// testSubscribersNumber number of subscribers in tests
const testSubscribersNumber = 10

func TestSubGroup(t *testing.T) {
	cache := mocks.NewCache(t)
	mq := mocks.NewMQ(t)
	prices = &Prices{
		subGroup:     newSub(),
		messageQueue: mq,
		cache:        cache,
	}
	mq.On("GetPrices", mock.AnythingOfType("*context.emptyCtx")).Return(nil, nil).Maybe()
	cache.On("Put", mock.AnythingOfType("[]*model.Price")).Return().Maybe()

	subscribers := make([]chan struct{}, testSubscribersNumber)
	subID := make([]uuid.UUID, testSubscribersNumber)
	for i, id := range subID {
		id = uuid.New()
		subscribers[i] = make(chan struct{}, 1)
		prices.Sub(subscribers[i], id)
	}

	end := make(chan struct{}, 1)
	go prices.cycle(context.Background(), end)

	var wg sync.WaitGroup
	for _, sub := range subscribers {
		wg.Add(1)
		go func(sub chan struct{}) {
			defer wg.Done()
			defer prices.Done()
			<-sub
		}(sub)
	}
	wg.Wait()

	require.NoError(t, nil)
}
