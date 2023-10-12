package main

import (
	"Shop-api/goods_web/proto"
	"context"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial(
		"consul://192.168.1.103:8500/Goods_srv?wait=14s&tag=C",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	defer conn.Close()

	client := proto.NewGoodsClient(conn)
	list, err := client.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		PriceMin:    0,
		PriceMax:    0,
		IsHot:       false,
		IsNew:       false,
		IsTab:       false,
		TopCategory: 0,
		Pages:       0,
		PagePerNums: 0,
		KeyWords:    "",
		Brand:       0,
	})
	if err != nil {
		log.Println("......", err)
	}

	zap.S().Info("Get Goods List Successfully ", list.Data)
}
