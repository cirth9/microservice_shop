package utils

import (
	"go.uber.org/zap"
	"net"
)

func GetFreePort() (int, error) {
	// 动态获取可用端口
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	//如果address参数中的端口号为空或为“0”，则自动选择端口号
	l, err := net.Listen("tcp", addr.String())
	if err != nil {
		return 0, err
	}
	zap.S().Info(l.Addr().String())
	return l.Addr().(*net.TCPAddr).Port, nil
}
