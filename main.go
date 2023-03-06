// Package main
package main

import (
	"context"
	"fmt"
	"net"

	"github.com/OVantsevich/Price-Service/internal/config"
	"github.com/OVantsevich/Price-Service/internal/handler"
	"github.com/OVantsevich/Price-Service/internal/repository"
	"github.com/OVantsevich/Price-Service/internal/service"
	pr "github.com/OVantsevich/Price-Service/proto"

	pppr "github.com/OVantsevich/PriceProvider/proto"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", cfg.Port))
	if err != nil {
		defer logrus.Fatalf("error while listening port: %e", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})
	defer client.Close()

	redisRep := repository.NewRedis(client, cfg.StreamName)
	streamPool := repository.NewStreamPool()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.PriceProviderHost, cfg.PriceProviderHost), opts...)
	if err != nil {
		logrus.Fatal("Fatal Dial: ", err)
	}
	ppClient := pppr.NewPriceServiceClient(conn)
	priceProvider := repository.NewPriceProvider(ppClient)
	cls := make(chan struct{})
	priceService := service.NewPrices(context.Background(), streamPool, redisRep, priceProvider, "0", cls)
	defer close(cls)
	priceHandler := handler.NewPrice(priceService)

	ns := grpc.NewServer()
	pr.RegisterPriceServiceServer(ns, priceHandler)

	if err = ns.Serve(listen); err != nil {
		defer logrus.Fatalf("error while listening server: %e", err)
	}
}
