package main

import (
	"Shop-api/user_web/initialize"
	"go.uber.org/zap"
)

func main() {
	//初始化全局zap
	initialize.InitLogger()

	initialize.InitConfig()

	r := initialize.InitRouters()
	zap.S().Info("server start,port :", 9090)

	err := r.Run(":9090")
	if err != nil {
		zap.S().Panic("start failed", err.Error())
	}
}
