// Package repository repos
package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"Price-Service/internal/model"

	"github.com/go-redis/redis/v8"
)

// Redis entity
type Redis struct {
	Client     *redis.Client
	StreamName string
}

// NewRedis constructor
func NewRedis(client *redis.Client, streamName string) *Redis {
	rds := &Redis{Client: client, StreamName: streamName}
	return rds
}

// GetPrices get last prices from stream
func (c *Redis) GetPrices(ctx context.Context, count int64, start string) (map[string]*model.Price, string, int, error) {
	data, err := c.Client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{c.StreamName, start},
		Count:   count,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, "", 0, fmt.Errorf("redis - GetPrices - XRead: %w", err)
	}

	prices := make(map[string]*model.Price)
	var curPrices []*model.Price

	for _, message := range data[0].Messages {
		dataFromStream := []byte(message.Values["data"].(string))

		err = json.Unmarshal(dataFromStream, &curPrices)
		if err != nil {
			return nil, "", 0, fmt.Errorf("redis - GetPrices - Unmarshal: %w", err)
		}
		for _, p := range curPrices {
			val, ok := prices[p.Name]
			if ok {
				val.PurchasePrice = p.PurchasePrice
				val.SellingPrice = p.SellingPrice
			}
			prices[p.Name] = p
		}
	}

	return prices, data[0].Messages[len(data[0].Messages)-1].ID, len(data[0].Messages), nil
}
