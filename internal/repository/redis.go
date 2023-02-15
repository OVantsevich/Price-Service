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
func NewRedis(client *redis.Client, streamName, groupName string) *Redis {
	rds := &Redis{Client: client, StreamName: streamName}
	_, _ = client.XGroupCreateMkStream(context.Background(), streamName, groupName, "$").Result()
	return rds
}

// GetPrices get last prices from stream
func (c *Redis) GetPrices(ctx context.Context) ([]*model.Price, error) {
	data, err := c.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:   "gr",
		Streams: []string{c.StreamName, ">"},
		Count:   1,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("redis - GetPrices - XRead: %w", err)
	}

	message := data[0].Messages[0]
	dataFromStream := []byte(message.Values["prices"].(string))
	var prices []*model.Price
	err = json.Unmarshal(dataFromStream, &prices)
	if err != nil {
		return nil, fmt.Errorf("redis - GetPrices - Unmarshal: %w", err)
	}
	return prices, nil
}
