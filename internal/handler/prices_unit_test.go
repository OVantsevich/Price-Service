package handler

import (
	"Price-Service/internal/handler/mocks"
	"Price-Service/internal/model"
	pr "Price-Service/proto"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

func server(ctx context.Context) (pr.PriceServiceClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pr.RegisterPriceServiceServer(baseServer, prices)
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			logrus.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			logrus.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pr.NewPriceServiceClient(conn)

	return client, closer
}

func testGetPriceNames() []string {
	return []string{
		"gold",
		"oil",
		"tesla",
		"google",
	}
}

func TestPrices_GetPrices(t *testing.T) {
	priceService := mocks.NewPriceService(t)
	prices = &Prices{
		pr.UnimplementedPriceServiceServer{},
		priceService,
	}
	client, closer := server(context.Background())
	defer closer()
	channel := make(chan []*model.Price, 1)
	priceService.On("Subscribe", mock.AnythingOfType("uuid.UUID")).Return(channel, nil).Once()
	priceService.On("UpdateSubscription", mock.AnythingOfType("[]string"), mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
	priceService.On("DeleteSubscription", mock.AnythingOfType("uuid.UUID")).Return().Maybe()

	outClient, err := client.GetPrices(context.Background())
	require.NoError(t, err)

	err = outClient.Send(&pr.GetPricesRequest{Names: testGetPriceNames()})

	channel <- []*model.Price{
		{
			Name:          "gold",
			SellingPrice:  50,
			PurchasePrice: 50,
		},
		{
			Name:          "oil",
			SellingPrice:  50,
			PurchasePrice: 50,
		},
		{
			Name:          "tesla",
			SellingPrice:  50,
			PurchasePrice: 50,
		},
		{
			Name:          "google",
			SellingPrice:  50,
			PurchasePrice: 50,
		},
	}

	//time.Sleep(time.Second * 2)

	resp, err := outClient.Recv()

	fmt.Print(resp.Prices)

	require.NoError(t, err)
}
