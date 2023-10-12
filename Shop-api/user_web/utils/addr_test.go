package utils

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
)

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort()
	if err != nil {
		zap.S().Errorw("error ", err.Error())
		return
	}
	fmt.Println(port)
}
