// Package repository repos
package repository

import (
	"sync"

	"Price-Service/internal/model"

	"github.com/google/uuid"
)

// Subscribers storing subscribers
type Subscribers map[uuid.UUID]chan *model.Price

// StreamPool cache of prices
type StreamPool struct {
	MU     sync.RWMutex
	prices map[string]Subscribers
}

// NewStreamPool constructor
func NewStreamPool() *StreamPool {
	return &StreamPool{prices: make(map[string]Subscribers)}
}

// Update add new pairs: price-stream
func (s *StreamPool) Update(streamID uuid.UUID, streamChan chan *model.Price, prices []string) {
	s.MU.Lock()
	for _, p := range prices {
		cp, ok := s.prices[p]
		if ok {
			cp[streamID] = streamChan
			continue
		}
		s.prices[p] = make(Subscribers)
		s.prices[p][streamID] = streamChan
	}
	s.MU.Unlock()
}

// Delete remove pairs: price-stream
func (s *StreamPool) Delete(streamID uuid.UUID) {
	s.MU.Lock()
	for _, p := range s.prices {
		delete(p, streamID)
	}
	s.MU.Unlock()
}

// Send prices to streams
func (s *StreamPool) Send(prices []*model.Price) {
	s.MU.RLock()
	for _, p := range prices {
		for _, c := range s.prices[p.Name] {
			select {
			case c <- p:
			default:
				go func(goP *model.Price, goC chan *model.Price) {
					goC <- goP
				}(p, c)
			}
		}
	}
	s.MU.RUnlock()
}
