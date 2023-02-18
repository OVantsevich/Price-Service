// Package model stream
package model

// Stream grpc stream channel and list of required prices
type Stream struct {
	Channel    chan []*Price
	PricesList []string
}
