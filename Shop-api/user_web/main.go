package main

import (
	"Shop-api/user_web/initialize"
	"Shop-api/user_web/utils"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	var err error
	//初始化全局zap
	initialize.InitLogger()

	initialize.InitConfig()

	r := initialize.InitRouters()

	port, err := utils.GetFreePort()
	zap.S().Info("server start,port :", port)
	addr := fmt.Sprintf(":%d", port)
	err = r.Run(addr)
	if err != nil {
		zap.S().Panic("start failed", err.Error())
	}
}
