// Package model models
package model

// Prices name and price storage map
type Prices map[string]float32

// GetPricesKeys get map keys
func (p *Prices) GetPricesKeys() []string {
	return []string{"gold", "oil", "tesla", "google"}
}
