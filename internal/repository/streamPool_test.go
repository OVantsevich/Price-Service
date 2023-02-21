package repository

import (
	"fmt"
	"testing"

	"Price-Service/internal/model"

	"github.com/google/uuid"
)

func TestCache_Get(t *testing.T) {
	prices := testGetStartPrices()
	names := make([]string, len(prices))
	for i := range prices {
		names[i] = prices[i].Name
	}

	streamID2 := uuid.New()
	streamID1 := uuid.New()
	streamChan1 := make(chan *model.Price, 100)
	streamChan2 := make(chan *model.Price, 100)

	streamPool.Update(streamID1, streamChan1, names)
	streamPool.Update(streamID2, streamChan2, names)

	streamPool.Send(prices)

	for range names {
		select {
		case i := <-streamChan1:
			fmt.Println(i)
		default:
		}
		select {
		case j := <-streamChan2:
			fmt.Println(j)
		default:
		}
	}

	streamPool.Delete(streamID1)

	streamPool.Send(prices)

	for range names {
		select {
		case i := <-streamChan1:
			fmt.Println(i)
		default:
		}
		select {
		case j := <-streamChan2:
			fmt.Println(j)
		default:
		}
	}

	streamPool.Delete(streamID2)
	streamPool.Delete(streamID2)
	streamPool.Delete(streamID2)

	streamPool.Send(prices)

	for range names {
		select {
		case i := <-streamChan1:
			fmt.Println(i)
		default:
		}
		select {
		case j := <-streamChan2:
			fmt.Println(j)
		default:
		}
	}
}
