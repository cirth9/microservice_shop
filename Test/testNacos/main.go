package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "192.168.1.109",
			Port:   8848,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         "0484a732-86ac-4e39-a6c0-f478b7a76549", // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
	}

	// 创建动态配置客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	config, err := client.GetConfig(vo.ConfigParam{
		DataId: "test-nacos.yaml",
		Group:  "test",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
