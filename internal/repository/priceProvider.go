// Package repository repos
package repository

import (
	"context"
	"fmt"

	"github.com/OVantsevich/Price-Service/internal/model"
	pr "github.com/OVantsevich/PriceProvider/proto"
)

// PriceProvider entity
type PriceProvider struct {
	client pr.PriceServiceClient
}

// NewPriceProvider constructor
func NewPriceProvider(client pr.PriceServiceClient) *PriceProvider {
	return &PriceProvider{client}
}

// GetCurrentPrices get current prices
func (p *PriceProvider) GetCurrentPrices(ctx context.Context, names []string) (map[string]*model.Price, error) {
	res, err := p.client.GetCurrentPrices(ctx, &pr.GetPricesRequest{Names: names})
	if err != nil {
		return nil, fmt.Errorf("priceProvider - GetCurrentPrices - GetCurrentPrices: %e", err)
	}
	prices := fromGRPC(res.Prices)
	return prices, nil
}

func fromGRPC(req map[string]*pr.Price) map[string]*model.Price {
	result := make(map[string]*model.Price, len(req))
	for _, r := range result {
		result[r.Name] = &model.Price{
			Name:          r.Name,
			SellingPrice:  r.SellingPrice,
			PurchasePrice: r.PurchasePrice,
		}
	}
	return result
}
