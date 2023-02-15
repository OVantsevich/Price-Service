// Package handler handle
package handler

import (
	"Price-Provider/internal/model"
	pr "Price-Provider/proto"
	"context"

	"github.com/google/uuid"
)

// PriceService service interface for user handler
//
//go:generate mockery --name=PriceService --case=underscore --output=./mocks
type PriceService interface {
	GetPrices(ctx context.Context) ([]*model.Price, error)

	AddChannel(chan []*model.Price, uuid.UUID)
	DelChannel(uuid.UUID)
}

// Price handler
type Price struct {
	pr.UnimplementedPriceServiceServer
	s PriceService
}

// NewPrice constructor
func NewPrice(s PriceService) *Price {
	return &Price{s: s}
}

// GetPrices add new grpc stream to stream slice
func (h *Price) GetPrices(_ *pr.GetPricesRequest, server pr.PriceService_GetPricesServer) error {
	var streamChan = make(chan []*model.Price, 1)
	var id = uuid.New()
	h.s.AddChannel(streamChan, id)

	for {
		p, open := <-streamChan
		if !open {
			return nil
		}
		protoPrices := make([]*pr.Price, len(p))
		for i := range protoPrices {
			protoPrices[i] = &pr.Price{
				Name:          p[i].Name,
				SellingPrice:  p[i].SellingPrice,
				PurchasePrice: p[i].PurchasePrice,
			}
		}
		err := server.Send(&pr.GetPricesResponse{
			Prices: protoPrices,
		})
		if err != nil {
			h.s.DelChannel(id)
			return err
		}
	}
}
