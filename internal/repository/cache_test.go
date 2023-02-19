package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache_Get(t *testing.T) {
	prices := testGetStartPrices()
	cache.Put(prices)

	for _, p := range prices {
		m, err := cache.Get(p.Name)
		require.NoError(t, err)
		require.Equal(t, p, m)
	}
}
