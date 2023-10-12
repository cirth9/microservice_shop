package initialize

import (
	"Shop-api/user_web/config"
	"Shop-api/user_web/global"
	"Shop-api/user_web/proto"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	isDebug             bool
	groupId             string
	nacosConfigFileName string
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
		config.TheServerConfig.UserSrvInfo.Name,
		"15s", config.TheServerConfig.UserSrvInfo.Tags[0])
	zap.S().Info(target)
	zap.S().Info("[GRPC_Target]  ", target)

	//负载均衡，轮询
	global.UserConn, err = grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Errorw("GRPC服务连接失败",
			"error", err.Error())
		return
	}

	global.UserClient = proto.NewUserClient(global.UserConn)

}

func initRedisConfig() {
	global.Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			config.TheServerConfig.UserRedisInfo.Host,
			config.TheServerConfig.UserRedisInfo.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func InitConfig() {
	//初始化基本配置信息
	v := viper.New()
	if isDebug = GetEnv("SHOP_ENV", v); isDebug {
		groupId = "dev"
	} else {
		groupId = "pro"
	}

	log.Println(isDebug, groupId)
	nacosConfigFileName = "config_nacos.yaml"
	v.SetConfigFile(nacosConfigFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicf(" 配置文件读取错误,err:%s", err.Error())
	}

	zap.S().Debug("name ", viper.GetString("name"))

	err := v.Unmarshal(&config.TheNacosConfig)
	if err != nil {
		zap.S().Panicf("配置文件反序列化错误，err:%s", err.Error())
	}

	log.Println(config.TheNacosConfig)

	serverConfig := []constant.ServerConfig{
		{
			IpAddr: config.TheNacosConfig.NacosServer.Ip,
			Port:   config.TheNacosConfig.NacosServer.Port,
		},
	}

	clientConfig := constant.ClientConfig{
		TimeoutMs:            config.TheNacosConfig.NacosClient.TimeoutMs,
		NamespaceId:          config.TheNacosConfig.NacosClient.NamespaceId,
		CacheDir:             config.TheNacosConfig.NacosClient.CacheDir,
		NotLoadCacheAtStart:  config.TheNacosConfig.NacosClient.NotLoadCacheAtStart,
		UpdateCacheWhenEmpty: false,
		LogDir:               config.TheNacosConfig.NacosClient.LogDir,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfig,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		zap.S().Fatal(err)
	}

	getConfig, err := configClient.GetConfig(vo.ConfigParam{
		DataId: config.TheNacosConfig.NacosServer.DataId,
		Group:  groupId,
	})
	log.Println("getConfig:", getConfig)
	err = json.Unmarshal([]byte(getConfig), &config.TheServerConfig)
	if err != nil {
		zap.S().Fatal(err)
	}

	log.Println(config.TheServerConfig)
	//redis config
	initRedisConfig()

	initGrpcConfig()
}

func WatchConfigChange(v *viper.Viper) {
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件被更改，Name:%s", in.Name)
	})
}
