package main

//
//import (
//	"MircoServer/goods_srv/proto"
//	"context"
//	"fmt"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//	"log"
//)
//
//var goodsClient proto.goodsClient
//var conn *grpc.ClientConn
//var err error
//
//func Init() {
//	conn, err = grpc.Dial(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		log.Println(">>>>>")
//		log.Println(err)
//		return
//	}
//	goodsClient = proto.NewoodsClient(conn)
//
//}
//
//func TestGetgoodsList() {
//	list, err1 := goodsClient.GetgoodsList(context.Background(), &proto.PageInfo{
//		Page:  1,
//		PSize: 10,
//	})
//	fmt.Println(list.Total)
//	if err1 != nil {
//		log.Println("GetgoodsList Error\n", err1)
//		return
//	}
//	for _, v := range list.Data {
//		fmt.Println(v)
//	}
//	for _, v := range list.Data {
//		rsp, err2 := goodsClient.CheckPassword(context.Background(), &proto.PasswordInfo{
//			Password:        "Test",
//			EncryptPassword: v.Password,
//		})
//		if err2 != nil {
//			fmt.Println("......", err1)
//			return
//		}
//		fmt.Printf("%+v \n", rsp.IsOK)
//	}
//}
//
//func main() {
//	Init()
//	TestGetgoodsList()
//
//	err = conn.Close()
//	if err != nil {
//		return
//	}
//}
