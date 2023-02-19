// Package repository repos
package repository

import (
	"fmt"
	"sync"

	"Price-Service/internal/model"
)

// CacheElement storing prices and subscribers
type CacheElement struct {
	subs map[string]struct{}
	*model.Price
}

// Cache cache of prices
type Cache struct {
	sync.RWMutex
	data map[string]*CacheElement
}

// NewCache constructor
func NewCache() *Cache {
	return &Cache{data: make(map[string]*CacheElement)}
}

// Get get price from map
func (c *Cache) Get(key string) (*model.Price, error) {
	c.RLock()
	v, ok := c.data[key]
	if ok {
		c.RUnlock()
		return v, nil
	}
	c.RUnlock()
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
