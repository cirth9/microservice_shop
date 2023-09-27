package main

import (
	"MircoServer/user_srv/config"
	"MircoServer/user_srv/handler"
	"MircoServer/user_srv/initialize"
	"MircoServer/user_srv/proto"

	"flag"
	"fmt"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

func main() {
	initialize.InitConfig()
	IP := flag.String("ip", "192.168.1.103", "ip地址")
	Port := flag.Int("port", 8002, "port地址")
	flag.Parse()
	userServer := grpc.NewServer()
	proto.RegisterUserServer(userServer, &handler.UserServer{})

	//注册健康检查
	grpc_health_v1.RegisterHealthServer(userServer, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", config.TheServerConfig.ConsulConfig.Host,
		config.TheServerConfig.ConsulConfig.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("192.168.1.103:%d", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = config.TheServerConfig.Name
	serviceID := config.TheServerConfig.Name
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = []string{"C", "I", "R", "N"}
	registration.Address = "192.168.1.103"
	registration.Check = check
	//1. 如何启动两个服务
	//2. 即使我能够通过终端启动两个服务，但是注册到consul中的时候也会被覆盖
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	listenAddress := fmt.Sprintf("%s:%d", *IP, *Port)
	fmt.Println(listenAddress)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Println(">>>>>", err)
		return
	}
	err = userServer.Serve(lis)
	if err != nil {
		log.Println(">>>>>", err)
		return
	}
}
