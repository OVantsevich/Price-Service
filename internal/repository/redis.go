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
func (c *Redis) GetPrices(ctx context.Context, offset string) ([]*model.Price, string, error) {
	data, err := c.Client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{c.StreamName, offset},
		Count:   1,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, "", fmt.Errorf("redis - GetPrices - XRead: %w", err)
	}

	message := data[0].Messages[0]
	dataFromStream := []byte(message.Values["data"].(string))
	var prices []*model.Price
	err = json.Unmarshal(dataFromStream, &prices)
	if err != nil {
		return nil, "", fmt.Errorf("redis - GetPrices - Unmarshal: %w", err)
	}
	return prices, message.ID, nil
}
