// Package service streams
package service

import (
	"sync"

	"Price-Service/internal/model"

	"github.com/google/uuid"
)

// Streams streams service
type Streams struct {
	mutex sync.Mutex
	data  map[string]*model.Stream
	cache Cache
}

// AddStream add new stream
func (s *Streams) AddStream(stream *model.Stream, id uuid.UUID) {
	s.mutex.Lock()
	s.data[id.String()] = stream
	s.mutex.Unlock()
}

// UpdateStream update stream price list
func (s *Streams) UpdateStream(list []string, id uuid.UUID) {
	s.mutex.Lock()
	s.data[id.String()].PricesList = list
	s.mutex.Unlock()
}

// DelStream delete stream
func (s *Streams) DelStream(id uuid.UUID) {
	s.mutex.Lock()
	delete(s.data, id.String())
	s.mutex.Unlock()
}

// SendPrices send prices
func (s *Streams) SendPrices() {
	s.mutex.Lock()
	for _, str := range s.data {
		var prices = make([]*model.Price, len(str.PricesList))
		for i, name := range str.PricesList {
			prices[i], _ = s.cache.Get(name)
		}
		str.Channel <- prices
	}
	s.mutex.Unlock()
}
