// Package repository repos
package repository

import (
	"fmt"
	"sync"

	"Price-Service/internal/model"
)

// Cache cache of prices
type Cache struct {
	sync.Mutex
	data map[string]*model.Price
}

// NewCache constructor
func NewCache() *Cache {
	return &Cache{data: make(map[string]*model.Price)}
}

// Get get price from map
func (c *Cache) Get(key string) (*model.Price, error) {
	c.Lock()
	v, ok := c.data[key]
	if ok {
		c.Unlock()
		return v, nil
	}
	c.Unlock()
	return nil, fmt.Errorf("%s does not exists", key)
}

// Put put price into map
func (c *Cache) Put(prices []*model.Price) {
	c.Lock()
	for _, p := range prices {
		c.data[p.Name] = p
	}
	c.Unlock()
}

// Size get cache size
func (c *Cache) Size() int {
	return c.Size()
}
