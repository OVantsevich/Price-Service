// Package handler handle
package handler

import (
	"context"
	"sync"

	"Price-Provider/internal/model"
	pr "Price-Provider/proto"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PriceService service interface for user handler
//
//go:generate mockery --name=PriceService --case=underscore --output=./mocks
type PriceService interface {
	GetPrices(ctx context.Context) (model.Prices, error)
}

// StreamGRPC service interface for grpc stream
//
//go:generate mockery --name=StreamGRPC --case=underscore --output=./mocks
type StreamGRPC interface {
	Send(*pr.GetPricesResponse) error
}

// grpc streams sync slice
type streams struct {
	MU   sync.Mutex
	data map[string]chan model.Prices
}

// Price handler
type Price struct {
	pr.UnimplementedPriceServiceServer
	s   PriceService
	str streams
}

// NewPrice constructor
func NewPrice(s PriceService) *Price {
	return &Price{s: s, str: streams{
		MU:   sync.Mutex{},
		data: make(map[string]chan model.Prices, 0),
	}}
}

// GetPrices add new grpc stream to stream slice
func (h *Price) GetPrices(_ *pr.GetPricesRequest, server pr.PriceService_GetPricesServer) error {
	h.str.MU.Lock()
	var streamChan = make(chan model.Prices, 1)
	var id = uuid.New()
	h.str.data[id.String()] = streamChan
	h.str.MU.Unlock()

	for {
		p, open := <-streamChan
		if !open {
			return nil
		}
		err := server.Send(&pr.GetPricesResponse{
			Prices: p,
		})
		if err != nil {
			h.str.MU.Lock()
			delete(h.str.data, id.String())
			h.str.MU.Unlock()
			return err
		}
	}
}

// Cycle getting data from redis and adding it to grpc streams
func (h *Price) Cycle(end <-chan struct{}) {
	var price model.Prices
	var err error

	for {
		select {
		case <-end:
			return
		default:
			price, err = h.s.GetPrices(context.Background())
			if err != nil {
				logrus.Fatalf("price - Cycle - GetPrices:%e", err)
				return
			}

			h.str.MU.Lock()
			for _, s := range h.str.data {
				s <- price
			}
			h.str.MU.Unlock()
		}
	}
}
