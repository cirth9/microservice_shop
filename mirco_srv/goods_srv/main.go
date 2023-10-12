package main

import (
	"MircoServer/goods_srv/config"
	"MircoServer/goods_srv/handler"
	"MircoServer/goods_srv/initialize"
	"MircoServer/goods_srv/proto"
	"MircoServer/goods_srv/utils"
	"MircoServer/goods_srv/utils/registery/consul"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"os"
	signal2 "os/signal"
	"syscall"

	"flag"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"log"
	"net"
)

func main() {
	initialize.InitConfig()
	IP := flag.String("ip", config.TheServerConfig.Host, "ip地址")
	Port := flag.Int("port", 0, "port地址")
	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	goodsServer := grpc.NewServer()
	proto.RegisterGoodsServer(goodsServer, &handler.GoodsServer{})

	//注册健康检查
	grpc_health_v1.RegisterHealthServer(goodsServer, health.NewServer())

	//服务注册
	consulRegistry := consul.NewRegistry(config.TheServerConfig.ConsulConfig.Host, config.TheServerConfig.ConsulConfig.Port)
	serviceId := uuid.NewV4().String()
	err2 := consulRegistry.Register(*IP, *Port, config.TheServerConfig.Tags, config.TheServerConfig.Name, serviceId)
	if err2 != nil {
		zap.S().Errorw("consul registry error: ", err2.Error())
	}

	listenAddress := fmt.Sprintf("%s:%d", *IP, *Port)
	zap.S().Info(listenAddress)

	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Println(">>>>>", err)
		return
	}

	//异步执行
	go func() {
		err = goodsServer.Serve(lis)
		if err != nil {
			log.Println(">>>>>", err)
			return
		}
	}()

	signal := make(chan os.Signal)
	signal2.Notify(signal, syscall.SIGINT, syscall.SIGTERM)
	<-signal
	if err = consulRegistry.DeRegister(serviceId); err != nil {
		zap.S().Info("deregister failed")
	}
	zap.S().Info("deregister successfully")
}
