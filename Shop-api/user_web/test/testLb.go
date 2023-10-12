package main

import (
	"Shop-api/user_web/proto"
	"context"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial(
		"consul://192.168.1.103:8500/user_srv?wait=14s&tag=C",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	defer conn.Close()

	client := proto.NewUserClient(conn)
	list, err := client.GetUserList(context.Background(), &proto.PageInfo{
		Page:  1,
		PSize: 2,
	})
	if err != nil {
		log.Println("......", err)
	}

	zap.S().Info("Get User List Successfully ", list.Data)
}
