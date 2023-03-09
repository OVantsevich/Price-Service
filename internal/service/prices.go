// Package service serv
package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/OVantsevich/Price-Service/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// batchSize number of messages read from mq
const batchSize = 100

// bufferSize number of messages stored for every grpc stream
const bufferSize = 1000

// MQ message queue
//
//go:generate mockery --name=MQ --case=underscore --output=./mocks
type MQ interface {
	GetPrices(ctx context.Context, count int64, start string) ([][]*model.Price, string, int, error)
}

// PriceProvider grpc price provider
//
//go:generate mockery --name=PriceProvider --case=underscore --output=./mocks
type PriceProvider interface {
	GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error)
}

// StreamPoolRepository pool of channels from prices to streams
//
//go:generate mockery --name=StreamPoolRepository --case=underscore --output=./mocks
type StreamPoolRepository interface {
	Send(prices []*model.Price)
	Update(streamID uuid.UUID, streamChan chan *model.Price, prices []string)
	Delete(streamID uuid.UUID)
}

// Prices prices service
type Prices struct {
	messageQueue  MQ
	streamPool    StreamPoolRepository
	priceProvider PriceProvider
	sMap          sync.Map
}

// NewPrices constructor
func NewPrices(ctx context.Context, spr StreamPoolRepository, mq MQ, pp PriceProvider, startPosition string, end <-chan struct{}) *Prices {
	prc := &Prices{messageQueue: mq, priceProvider: pp, streamPool: spr}
	go prc.cycle(ctx, end, startPosition)
	return prc
}

// Subscribe allocating new channel for grpc stream with id and returning it
func (p *Prices) Subscribe(streamID uuid.UUID) chan *model.Price {
	streamChan := make(chan *model.Price, bufferSize)
	p.sMap.Store(streamID, streamChan)
	return streamChan
}

// UpdateSubscription clear previous stream subscription and creat new
func (p *Prices) UpdateSubscription(names []string, streamID uuid.UUID) error {
	streamChan, ok := p.sMap.Load(streamID)
	if !ok {
		return fmt.Errorf("not found")
	}
	p.streamPool.Delete(streamID)
	p.streamPool.Update(streamID, streamChan.(chan *model.Price), names)
	return nil
}

// DeleteSubscription delete grpc stream subscription and close it's chan
func (p *Prices) DeleteSubscription(streamID uuid.UUID) error {
	streamChan, ok := p.sMap.Load(streamID)
	if !ok {
		return fmt.Errorf("not found")
	}
	p.streamPool.Delete(streamID)
	close(streamChan.(chan *model.Price))
	p.sMap.Delete(streamID)
	return nil
}

// GetCurrentPrices get current price from price provider
func (p *Prices) GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error) {
	return p.priceProvider.GetCurrentPrices(ctx, names)
}

// cycle getting data from redis and adding it to grpc streams
func (p *Prices) cycle(ctx context.Context, end <-chan struct{}, sp string) {
	var prices [][]*model.Price
	var startPosition = sp
	var err error

	for {
		select {
		case <-end:
			return
		default:
			prices, startPosition, _, err = p.messageQueue.GetPrices(ctx, batchSize, startPosition)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"batchSize":     batchSize,
					"startPosition": startPosition,
				}).Fatalf("prices - cycle - GetPrices: %v", err)
				return
			}
			for _, price := range prices {
				p.streamPool.Send(price)
			}
		}
	}
}
