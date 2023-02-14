// Package repository repos
package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"Price-Provider/internal/model"

	"github.com/go-redis/redis/v8"
)

// Redis entity
type Redis struct {
	Client     *redis.Client
	StreamName string
}

// NewRedis constructor
func NewRedis(client *redis.Client, streamName string) *Redis {
	return &Redis{Client: client, StreamName: streamName}
}

// GetPrices get last prices from stream
func (c *Redis) GetPrices(ctx context.Context) (model.Prices, error) {
	data, err := c.Client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{c.StreamName, "0"},
		Count:   1,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("redis - GetPrices - XRead: %w", err)
	}

	message := data[0].Messages[0]
	dataFromStream := []byte(message.Values["prices"].(string))
	var prices = model.Prices{}
	err = json.Unmarshal(dataFromStream, &prices)
	if err != nil {
		return nil, fmt.Errorf("redis - GetPrices - Unmarshal: %w", err)
	}
	return prices, nil
}
