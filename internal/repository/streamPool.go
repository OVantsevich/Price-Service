// Package repository repos
package repository

import (
	"sync"

	"Price-Service/internal/model"
)

// Subscribers storing subscribers
type Subscribers map[string]chan *model.Price

// StreamPool cache of prices
type StreamPool struct {
	sync.RWMutex
	prices map[string]Subscribers
}

// NewStreamPool constructor
func NewStreamPool() *StreamPool {
	return &StreamPool{prices: make(map[string]Subscribers)}
}

// Update add new pairs: price-stream
func (s *StreamPool) Update(streamID string, streamChan chan *model.Price, prices []string) error {
	s.Lock()
	for _, p := range prices {
		cp, ok := s.prices[p]
		if ok {
			cp[streamID] = streamChan
			continue
		}
		s.prices[p] = make(Subscribers)
		s.prices[p][streamID] = streamChan
	}
	s.Unlock()
	return nil
}

// Delete remove pairs: price-stream
func (s *StreamPool) Delete(streamID string, prices []string) error {
	s.Lock()
	for _, p := range prices {
		delete(s.prices[p], streamID)
	}
	s.Unlock()
	return nil
}

// Send prices to streams
func (s *StreamPool) Send(prices []*model.Price) {
	s.RLock()
	for _, p := range prices {
		for _, c := range s.prices[p.Name] {
			select {
			case c <- p:
				break
			default:
				go func(goP *model.Price, goC chan *model.Price) {
					goC <- goP
				}(p, c)
			}
		}
	}
	s.RUnlock()
}
