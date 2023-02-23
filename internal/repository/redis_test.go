package repository

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/OVantsevich/Price-Service/internal/model"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// 24 bits “mantissa”, otherwise known as a coefficient or significand.
const maxInt int64 = 1 << 24

// starting prices
const startPrice float64 = 50

// maxPrice max price
const maxChange float64 = 5

// numberOfPublishers number of publishers
const numberOfPublishers int = 5

// numberOfMessages number of publishers
const numberOfMessages int = 10000

// Float64 random float64 using crypto/rand
func Float64() float64 {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(maxInt))
	return float64(nBig.Int64()) / float64(maxInt)
}

func testGetStartPrices() []*model.Price {
	return []*model.Price{
		{
			Name:          "gold",
			SellingPrice:  startPrice,
			PurchasePrice: startPrice + Float64(),
		},
		{
			Name:          "oil",
			SellingPrice:  startPrice,
			PurchasePrice: startPrice + Float64(),
		},
		{
			Name:          "tesla",
			SellingPrice:  startPrice,
			PurchasePrice: startPrice + Float64(),
		},
		{
			Name:          "google",
			SellingPrice:  startPrice,
			PurchasePrice: startPrice + Float64(),
		},
	}
}

func publishPrices(ctx context.Context, client *redis.Client, streamName string, messageNumber int) {
	var prices = testGetStartPrices()

	for i := 0; i < messageNumber; i++ {
		for _, pr := range prices {
			chg := -maxChange + (2*maxChange)*Float64()
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
	var id = "0"
	var err error
	var cnorm int
	for ; numberOfReceivedMessages < numberOfMessages*numberOfPublishers; numberOfReceivedMessages += cnorm {
		_, id, cnorm, err = redisRps.GetPrices(ctx, 100, id)
		require.NoError(t, err)
	}
	logrus.Infof("end time: %s", time.Now())
	require.Equal(t, numberOfReceivedMessages, numberOfMessages*numberOfPublishers)
}
