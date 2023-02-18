package repository

import (
	"Price-Service/internal/model"
	"context"
	"crypto/rand"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"math/big"
	"runtime"
	"testing"
	"time"
)

// 24 bits “mantissa”, otherwise known as a coefficient or significand.
const maxInt int64 = 1 << 24

// starting prices
const startPrice float32 = 50

// maxPrice max price
const maxChange float32 = 5

// numberOfPublishers number of publishers
const numberOfPublishers int = 5

// numberOfMessages number of publishers
const numberOfMessages int = 10000

// Float32 random float32 using crypto/rand
func Float32() float32 {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(maxInt))
	return float32(nBig.Int64()) / float32(maxInt)
}

func testGetStartPrices() []*model.Price {
	return []*model.Price{
		{
			Name:          "gold",
			SellingPrice:  startPrice,
			PurchasePrice: startPrice + Float32(),
		},
		{
			Name:          "oil",
			SellingPrice:  startPrice,
			PurchasePrice: startPrice + Float32(),
		},
		{
			Name:          "tesla",
			SellingPrice:  startPrice,
			PurchasePrice: startPrice + Float32(),
		},
		{
			Name:          "google",
			SellingPrice:  startPrice,
			PurchasePrice: startPrice + Float32(),
		},
	}
}

func publishPrices(ctx context.Context, client *redis.Client, streamName string, messageNumber int) {
	var prices = testGetStartPrices()

	for i := 0; i < messageNumber; i++ {
		for _, pr := range prices {
			chg := -maxChange + (2*maxChange)*Float32()
			pr.SellingPrice += chg
			pr.PurchasePrice += chg
		}
		mp, _ := json.Marshal(prices)
		_, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{
				"data": mp,
			},
		}).Result()
		if err != nil {
			return
		}
	}
}

func TestRedis_GetPrices(t *testing.T) {
	ctx := context.Background()
	logrus.Info(runtime.NumCPU())
	for i := 0; i < numberOfPublishers; i++ {
		go publishPrices(ctx, redisRps.Client, testRedisStream, numberOfMessages)
	}

	var numberOfReceivedMessages = 1
	logrus.Infof("begin time: %s", time.Now())
	_, offset, err := redisRps.GetPrices(ctx, "0")
	require.NoError(t, err)
	for ; numberOfReceivedMessages < numberOfMessages*numberOfPublishers; numberOfReceivedMessages++ {
		_, offset, err = redisRps.GetPrices(ctx, offset)
		require.NoError(t, err)
	}
	logrus.Infof("end time: %s", time.Now())
	require.Equal(t, numberOfReceivedMessages, numberOfMessages*numberOfPublishers)
}
