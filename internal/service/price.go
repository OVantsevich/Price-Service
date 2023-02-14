// Package service serv
package service

import (
	"Price-Provider/internal/model"
	"context"
	"fmt"
)

// MQ stream interface for user service
//
//go:generate mockery --name=MQ --case=underscore --output=./mocks
type MQ interface {
	GetPrices(ctx context.Context) (model.Prices, error)
}

// Prices price service
type Prices struct {
	messageQueue MQ
}

// NewPrices constructor
func NewPrices(mq MQ) *Prices {
	return &Prices{messageQueue: mq}
}

// GetPrices get prices service
func (s *Prices) GetPrices(ctx context.Context) (model.Prices, error) {
	p, err := s.messageQueue.GetPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("prices - GetPrices - GetPrices: %w", err)
	}
	return p, nil
}
