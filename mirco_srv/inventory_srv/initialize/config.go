package initialize

import (
	"MircoServer/inventory_srv/config"
	"MircoServer/inventory_srv/global"
	"MircoServer/inventory_srv/model"
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	isDebug             bool
	nacosConfigFileName string
	groupId             string
)

func GetEnv(env string, v *viper.Viper) bool {
	v.AutomaticEnv()
	return v.GetBool(env)
}

func InitInventorySrvMysqlConfig() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.TheServerConfig.MysqlConfig.InventoryName,
		config.TheServerConfig.MysqlConfig.Password,
		config.TheServerConfig.MysqlConfig.Host,
		config.TheServerConfig.MysqlConfig.Port,
		config.TheServerConfig.MysqlConfig.DBName)
	zap.S().Info("mysql_dsn: ", dsn)
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
		LogLevel:                  logger.Info,
	})

	global.MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Println(">>>>>>", err)
		return
	}

	err = global.MysqlDB.AutoMigrate(&model.Inventory{}, &model.InventoryNew{}, &model.Delivery{}, &model.StockSellDetail{})
	if err != nil {
		log.Println(">>>>>>", err)
		return
	}
}

func InitConfig() {

	InitLogger()

	//v := viper.New()
	//if isDebug = GetEnv("SHOP_ENV", v); isDebug {
	//	configFileName = "config_debug.yaml"
	//} else {
	//	configFileName = "config_pro.yaml"
	//}
	//
	//log.Println(isDebug, configFileName)
	//v.SetConfigFile(configFileName)
	//if err := v.ReadInConfig(); err != nil {
	//	zap.S().Panicf(" 配置文件读取错误,err:%s", err.Error())
	//}
	//
	//zap.S().Debug("name", viper.GetString("name"))
	//
	//err := v.Unmarshal(&config.TheServerConfig)
	//if err != nil {
	//	zap.S().Panicf("配置文件反序列化错误，err:%s", err.Error())
	//}
	//
	//fmt.Println(config.TheServerConfig)

	v := viper.New()
	if isDebug = GetEnv("SHOP_ENV", v); isDebug {
		groupId = "dev"
	} else {
		groupId = "pro"
	}

	nacosConfigFileName = "config_nacos.yaml"
	v.SetConfigFile(nacosConfigFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panic(err)
	}

	err2 := v.Unmarshal(&config.TheNacosConfig)
	if err2 != nil {
		zap.S().Panic(err2)
	}

	log.Printf("%+v", config.TheNacosConfig)
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: config.TheNacosConfig.NacosServer.Ip,
			Port:   config.TheNacosConfig.NacosServer.Port,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         config.TheNacosConfig.NacosClient.NamespaceId, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           config.TheNacosConfig.NacosClient.TimeoutMs,
		NotLoadCacheAtStart: config.TheNacosConfig.NacosClient.NotLoadCacheAtStart,
		LogDir:              config.TheNacosConfig.NacosClient.LogDir,
		CacheDir:            config.TheNacosConfig.NacosClient.CacheDir,
	}

	// 创建动态配置客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	InventorySrvConfig, err := client.GetConfig(vo.ConfigParam{
		DataId: config.TheNacosConfig.NacosServer.DataId,
		Group:  groupId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(InventorySrvConfig)
	err = json.Unmarshal([]byte(InventorySrvConfig), &config.TheServerConfig)
	if err != nil {
		zap.S().Panic(err)
	}

	zap.S().Info("TheServerConfig: ", config.TheServerConfig)
	InitInventorySrvMysqlConfig()
}
