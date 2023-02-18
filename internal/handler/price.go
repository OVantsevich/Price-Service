// Package handler handle
package handler

import (
	"Price-Service/internal/model"
	pr "Price-Service/proto"

	"github.com/google/uuid"
)

// PriceService service interface for user handler
//
//go:generate mockery --name=PriceService --case=underscore --output=./mocks
type PriceService interface {
	Sub(c chan struct{}, id uuid.UUID)
	UnSub(id uuid.UUID)
	Lock()
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
func (h *Price) GetPrices(server pr.PriceService_GetPricesServer) error {
	var streamChan = make(chan map[string]*model.Price, 1)
	var id = uuid.New()
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

func (h *Price) change(server pr.PriceService_GetPricesServer) {
	var response *pr.GetPricesRequest
	var err error
	for {
		response, err = server.Recv()
		if err != nil {

		}
	}
}
