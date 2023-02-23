// Package handler handle
package handler

import (
	"context"

	pr "github.com/OVantsevich/Price-Service/proto"

	"github.com/OVantsevich/Price-Service/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PriceService service interface for prices handler
//
//go:generate mockery --name=PriceService --case=underscore --output=./mocks
type PriceService interface {
	Subscribe(streamID uuid.UUID) chan *model.Price
	UpdateSubscription(names []string, streamID uuid.UUID) error
	DeleteSubscription(streamID uuid.UUID) error

	GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error)
}

// Prices handler
type Prices struct {
	pr.UnimplementedPriceServiceServer
	service PriceService
}

// NewPrice constructor
func NewPrice(s PriceService) *Prices {
	return &Prices{service: s}
}

// GetCurrentPrices get current prices from price provider
func (h *Prices) GetCurrentPrices(ctx context.Context, request *pr.GetCurrentPricesRequest) (*pr.GetCurrentPricesResponse, error) {
	result, err := h.service.GetCurrentPrices(ctx, request.Names)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request.Names": request.Names,
		}).Errorf("prices - GetCurrentPrices - GetCurrentPrices: %e", err)
		return nil, status.Error(codes.Unknown, err.Error())
	}
	response := &pr.GetCurrentPricesResponse{
		Prices: make(map[string]*pr.Price),
	}
	for _, r := range result {
		response.Prices[r.Name] = &pr.Price{
			Name:          r.Name,
			SellingPrice:  r.SellingPrice,
			PurchasePrice: r.PurchasePrice,
		}
	}
	return response, err
}

// GetPrices add new grpc stream to stream slice
func (h *Prices) GetPrices(server pr.PriceService_GetPricesServer) (err error) {
	var streamID = uuid.New()
	streamChan := h.service.Subscribe(streamID)

	changeError := make(chan error, 1)
	changeStop := make(chan struct{}, 1)
	go h.change(server, streamID, changeStop, changeError)

	var open bool
	var currentPrices *model.Price
	for {
		select {
		case err = <-changeError:
			logrus.Errorf("prices - GetPrices - change: %e", err)
			return err
		case currentPrices, open = <-streamChan:
			if !open {
				logrus.Fatal("prices - GetPrices - <-streamChan")
				return nil
			}
			err = server.Send(&pr.GetPricesResponse{
				Prices: toGRPC([]*model.Price{currentPrices}),
			})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"grpcPrices": toGRPC([]*model.Price{currentPrices}),
				}).Errorf("prices - GetPrices - Send:%e", err)
				if err2 := h.service.DeleteSubscription(streamID); err2 != nil {
					logrus.WithFields(logrus.Fields{
						"streamID": streamID,
					}).Errorf("prices - GetPrices - DeleteSubscription: %e", err2)
				}

				changeStop <- struct{}{}
				return err
			}
		}
	}
}

func (h *Prices) change(server pr.PriceService_GetPricesServer, streamID uuid.UUID, stop chan struct{}, end chan error) {
	var response *pr.GetPricesRequest
	var err error

	for {
		select {
		case <-stop:
			return
		default:
			response, err = server.Recv()
			if err != nil {
				logrus.Errorf("prices - change - Recv: %v", err.Error())
				end <- status.Error(codes.DataLoss, err.Error())
				return
			}
			err = h.service.UpdateSubscription(response.Names, streamID)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"response.Names": response.Names,
					"streamID":       streamID,
				}).Errorf("prices - change - UpdateSubscription:%e", err)
				end <- status.Error(codes.Unknown, err.Error())
				return
			}
		}
	}
}

func toGRPC(req []*model.Price) []*pr.Price {
	result := make([]*pr.Price, len(req))
	for i := range result {
		result[i] = &pr.Price{
			Name:          req[i].Name,
			SellingPrice:  req[i].SellingPrice,
			PurchasePrice: req[i].PurchasePrice,
		}
	}
	return result
}
