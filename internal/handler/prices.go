// Package handler handle
package handler

import (
	"Price-Service/internal/model"
	pr "Price-Service/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PriceService service interface for prices handler
//
//go:generate mockery --name=PriceService --case=underscore --output=./mocks
type PriceService interface {
	Subscribe(names []string, streamID uuid.UUID) (chan []*model.Price, error)
	UpdateSubscription(names []string, streamID uuid.UUID) error
	DeleteSubscription(streamID uuid.UUID)
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

// GetPrices add new grpc stream to stream slice
func (h *Prices) GetPrices(server pr.PriceService_GetPricesServer) error {
	response, err := server.Recv()
	if err != nil {
		logrus.Errorf("prices - GetPrices - Recv:%e", err)
		return status.Errorf(codes.DataLoss, "no init data")
	}

	var streamID = uuid.New()
	streamChan, err := h.service.Subscribe(response.Names, streamID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"names":    response.Names,
			"streamID": streamID,
		}).Errorf("prices - GetPrices - Subscribe:%e", err)
		return status.Error(codes.Unknown, err.Error())
	}

	changeError := make(chan error, 1)
	changeStop := make(chan struct{}, 1)
	go h.change(server, streamID, changeStop, changeError)

	var open bool
	var currentPrices []*model.Price
	for {
		select {
		case err = <-changeError:
			logrus.Errorf("prices - GetPrices - change: %e", err)
			return err
		default:
			currentPrices, open = <-streamChan
			if !open {
				logrus.Fatal("prices - GetPrices - <-streamChan")
				return nil
			}
			err = server.Send(&pr.GetPricesResponse{
				Prices: toGRPC(currentPrices),
			})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"grpcPrices": toGRPC(currentPrices),
				}).Errorf("prices - GetPrices - Send:%e", err)
				h.service.DeleteSubscription(streamID)
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
			break
		default:
			response, err = server.Recv()
			if err != nil {
				logrus.Errorf("prices - change - Recv:%e", err)
				end <- status.Error(codes.DataLoss, err.Error())
				break
			}
			err = h.service.UpdateSubscription(response.Names, streamID)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"response.Names": response.Names,
					"streamID":       streamID,
				}).Errorf("prices - change - UpdateSubscription:%e", err)
				end <- status.Error(codes.Unknown, err.Error())
				break
			}
		}
	}
}

func toGRPC(req []*model.Price) []*pr.Price {
	result := make([]*pr.Price, len(req))
	for i := range result {
		result[i] = &pr.Price{
			Name:          result[i].Name,
			SellingPrice:  result[i].SellingPrice,
			PurchasePrice: result[i].PurchasePrice,
		}
	}
	return result
}
