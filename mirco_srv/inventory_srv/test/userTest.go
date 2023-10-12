package main

import (
	"MircoServer/inventory_srv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var InventoryClient proto.InventoryClient
var conn *grpc.ClientConn
var err error

func Init() {
	conn, err = grpc.Dial(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(">>>>>")
		log.Println(err)
		return
	}
	InventoryClient = proto.NewInventoryClient(conn)

}

func TestGetInventoryList() {
	//list, err1 := InventoryClient(context.Background(), &proto.PageInfo{
	//	Page:  1,
	//	PSize: 10,
	//})
	//fmt.Println(list.Total)
	//if err1 != nil {
	//	log.Println("GetInventoryList Error\n", err1)
	//	return
	//}
	//for _, v := range list.Data {
	//	fmt.Println(v)
	//}
	//for _, v := range list.Data {
	//	rsp, err2 := InventoryClient.CheckPassword(context.Background(), &proto.PasswordInfo{
	//		Password:        "Test",
	//		EncryptPassword: v.Password,
	//	})
	//	if err2 != nil {
	//		fmt.Println("......", err1)
	//		return
	//	}
	//	fmt.Printf("%+v \n", rsp.IsOK)
	//}
}

func main() {
	Init()
	TestGetInventoryList()

	err = conn.Close()
	if err != nil {
		return
	}
}
