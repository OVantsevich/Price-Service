package repository

import (
	"github.com/stretchr/testify/require"
	"testing"
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
