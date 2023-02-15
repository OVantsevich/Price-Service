// Package main
package main

import (
	"fmt"
	"net"

	"Price-Provider/internal/config"
	"Price-Provider/internal/handler"
	"Price-Provider/internal/repository"
	"Price-Provider/internal/service"
	pr "Price-Provider/proto"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

	redisRep := repository.NewRedis(client, cfg.StreamName, cfg.GroupName)
	cls := make(chan struct{})
	priceService := service.NewPrices(redisRep, cls)
	defer close(cls)
	priceHandler := handler.NewPrice(priceService)

	ns := grpc.NewServer()
	pr.RegisterPriceServiceServer(ns, priceHandler)

	if err = ns.Serve(listen); err != nil {
		defer logrus.Fatalf("error while listening server: %e", err)
	}
}
