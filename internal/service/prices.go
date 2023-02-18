// Package service serv
package service

import (
	"context"
	"sync"

	"Price-Service/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// MQ stream interface for prices service
//
//go:generate mockery --name=MQ --case=underscore --output=./mocks
type MQ interface {
	GetPrices(ctx context.Context) ([]*model.Price, error)
}

// Cache cache interface for prices service
//
//go:generate mockery --name=Cache --case=underscore --output=./mocks
type Cache interface {
	Get(key string) (*model.Price, error)
	Put(prices []*model.Price)
	Size() int
}

// subGroup subscribers
type subGroup struct {
	mutex sync.Mutex
	wg    sync.WaitGroup
	data  map[string]chan struct{}
}

// newSub creat new sub
func newSub() *subGroup {
	return &subGroup{data: make(map[string]chan struct{})}
}

// Sub add new sub
func (s *subGroup) Sub(c chan struct{}, id uuid.UUID) {
	s.mutex.Lock()
	s.data[id.String()] = c
	s.mutex.Unlock()
}

// UnSub delete sub
func (s *subGroup) UnSub(id uuid.UUID) {
	s.mutex.Lock()
	delete(s.data, id.String())
	s.mutex.Unlock()
}

// Done sync done
func (s *subGroup) Done() {
	s.wg.Done()
}

// Wait sync wait
func (s *subGroup) Wait() {
	s.mutex.Lock()
	s.wg.Wait()
	s.mutex.Unlock()
}

// pub publish info
func (s *subGroup) pub() {
	s.mutex.Lock()
	for _, c := range s.data {
		s.wg.Add(1)
		c <- struct{}{}
	}
	s.mutex.Unlock()
}

// Prices prices service
type Prices struct {
	*subGroup
	messageQueue MQ
	cache        Cache
}

// NewPrices constructor
func NewPrices(ctx context.Context, ch Cache, mq MQ, end <-chan struct{}) *Prices {
	prc := &Prices{messageQueue: mq, subGroup: newSub(), cache: ch}
	go prc.cycle(ctx, end)
	return prc
}

// cycle getting data from redis and adding it to grpc streams
func (s *Prices) cycle(ctx context.Context, end <-chan struct{}) {
	var prices []*model.Price
	var err error

	for {
		select {
		case <-end:
			return
		default:
			prices, err = s.messageQueue.GetPrices(ctx)
			if err != nil {
				logrus.Fatalf("prices - cycle - GetPrices:%e", err)
				return
			}
			s.cache.Put(prices)
			s.pub()
			s.wg.Wait()
		}
	}
}
