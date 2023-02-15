// Package service serv
package service

import (
	"context"
	"fmt"
	"sync"

	"Price-Provider/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// MQ stream interface for user service
//
//go:generate mockery --name=MQ --case=underscore --output=./mocks
type MQ interface {
	GetPrices(ctx context.Context) ([]*model.Price, error)
}

// grpc streams sync slice
type streams struct {
	MU   sync.Mutex
	data map[string]chan []*model.Price
}

// Price prices service
type Price struct {
	messageQueue MQ
	str          streams
}

// NewPrices constructor
func NewPrices(mq MQ, end <-chan struct{}) *Price {
	prc := &Price{messageQueue: mq, str: streams{
		MU:   sync.Mutex{},
		data: make(map[string]chan []*model.Price, 0),
	}}
	go prc.Cycle(end)
	return prc
}

// GetPrices get prices service
func (s *Price) GetPrices(ctx context.Context) ([]*model.Price, error) {
	p, err := s.messageQueue.GetPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("prices - GetPrices - GetPrices: %w", err)
	}
	return p, nil
}

// AddChannel add new channel
func (s *Price) AddChannel(ch chan []*model.Price, id uuid.UUID) {
	s.str.MU.Lock()
	s.str.data[id.String()] = ch
	s.str.MU.Unlock()
}

// DelChannel delete channel
func (s *Price) DelChannel(id uuid.UUID) {
	s.str.MU.Lock()
	delete(s.str.data, id.String())
	s.str.MU.Unlock()
}

// Cycle getting data from redis and adding it to grpc streams
func (s *Price) Cycle(end <-chan struct{}) {
	var price []*model.Price
	var err error

	for {
		select {
		case <-end:
			return
		default:
			price, err = s.GetPrices(context.Background())
			if err != nil {
				logrus.Fatalf("Prices - Cycle - GetPrices:%e", err)
				return
			}

			s.str.MU.Lock()
			for _, ch := range s.str.data {
				ch <- price
			}
			s.str.MU.Unlock()
		}
	}
}
