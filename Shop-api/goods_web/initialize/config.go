package initialize

import (
	"Shop-api/goods_web/config"
	"Shop-api/goods_web/global"
	"Shop-api/goods_web/proto"
	"fmt"
	"github.com/fsnotify/fsnotify"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	isDebug        bool
	GroupId        string
	configFileName string
)

func GetEnv(env string, v *viper.Viper) bool {
	v.AutomaticEnv()
	return v.GetBool(env)
}

func initGrpcConfig() {
	var err error
	target := fmt.Sprintf("consul://%s:%d/%s?wait=%s&tag=%s",
		config.TheServerConfig.ConsulConfigInfo.Host,
		config.TheServerConfig.ConsulConfigInfo.Port,
		config.TheServerConfig.GoodsSrvInfo.Name,
		"15s", config.TheServerConfig.GoodsSrvInfo.Tags[0])
	zap.S().Info(target)
	zap.S().Info("[GRPC_Target]  ", target)

	//负载均衡，轮询
	global.GoodsConn, err = grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Errorw("GRPC服务连接失败",
			"error", err.Error())
		return
	}

	global.GoodsClient = proto.NewGoodsClient(global.GoodsConn)

}

//func initRedisConfig() {
//	global.Rdb = redis.NewClient(&redis.Options{
//		Addr: fmt.Sprintf("%s:%d",
//			config.TheServerConfig.GoodsRedisInfo.Host,
//			config.TheServerConfig.GoodsRedisInfo.Port),
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//}

func InitConfig() {
	//初始化基本配置信息
	v := viper.New()
	if isDebug = GetEnv("SHOP_ENV", v); isDebug {
		configFileName = "config-debug.yaml"
	} else {
		configFileName = "config-pro.yaml"
	}

	log.Println(isDebug, configFileName)
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicf(" 配置文件读取错误,err:%s", err.Error())
	}

	zap.S().Debug("name ", viper.GetString("name"))

	err := v.Unmarshal(&config.TheServerConfig)
	if err != nil {
		zap.S().Panicf("配置文件反序列化错误，err:%s", err.Error())
	}

	fmt.Println(config.TheServerConfig)

	//grpc config
	initGrpcConfig()
}

func WatchConfigChange(v *viper.Viper) {
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件被更改，Name:%s", in.Name)
	})
}
