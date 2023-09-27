package initialize

import (
	"Shop-api/user_web/config"
	"Shop-api/user_web/global"
	"Shop-api/user_web/proto"
	"fmt"
	"github.com/fsnotify/fsnotify"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	isDebug        bool
	configFileName string
)

func GetEnv(env string, v *viper.Viper) bool {
	v.AutomaticEnv()
	return v.GetBool(env)
}

func initGrpcConfig() {
	//初始化grpc设置
	cfg := consulApi.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",
		config.TheServerConfig.ConsulConfigInfo.Host,
		config.TheServerConfig.ConsulConfigInfo.Port)
	userSrvHost := ""
	userSrvPort := 0
	client, err := consulApi.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,
		config.TheServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("[InitGrpcConfig] 连接 【用户服务失败】")
		return
	}

	global.GrpcAddress = fmt.Sprintf("%s:%d",
		userSrvHost,
		userSrvPort)
	zap.S().Info("GrpcAddress:", global.GrpcAddress)
	global.UserConn, err = grpc.Dial(global.GrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		configFileName = "user_web/config-debug.yaml"
	} else {
		configFileName = "user_web/config-pro.yaml"
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

	//redis config
	initRedisConfig()
}

func WatchConfigChange(v *viper.Viper) {
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件被更改，Name:%s", in.Name)
	})
}
