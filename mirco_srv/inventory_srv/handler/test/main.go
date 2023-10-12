package main

import (
	"MircoServer/inventory_srv/config"
	"MircoServer/inventory_srv/global"
	"MircoServer/inventory_srv/handler"
	"MircoServer/inventory_srv/model"
	"MircoServer/inventory_srv/proto"
	"context"
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
	"sync"
	"time"
)

var (
	wg              sync.WaitGroup
	inventoryServer = &handler.InventoryServer{}
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

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
func InitConfig() {

	InitLogger()

	v := viper.New()
	if isDebug = GetEnv("SHOP_ENV", v); isDebug {
		groupId = "dev"
	} else {
		groupId = "pro"
	}

	nacosConfigFileName = "E:\\Project\\micoservice_shop\\mirco_srv\\inventory_srv\\config_nacos.yaml"
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
}

func TestSell(Wg *sync.WaitGroup) {
	_, _ = inventoryServer.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{{
			GoodsId: 19,
			Num:     10,
		}},
	})
	Wg.Done()
}
func main() {
	InitConfig()
	var err error
	dsn := "root:qq31415926535--@tcp(127.0.0.1:3306)/shop_inventory_srv?charset=utf8mb4&parseTime=True&loc=Local"
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
	globalDb, err := global.MysqlDB.DB()
	if err != nil {
		zap.Error(err)
	}
	err = global.MysqlDB.AutoMigrate(&model.Inventory{}, &model.InventoryNew{}, &model.Delivery{}, &model.StockSellDetail{})
	globalDb.SetMaxIdleConns(300)
	globalDb.SetMaxOpenConns(600)
	globalDb.SetConnMaxLifetime(time.Minute)
	wg.Add(250)
	for i := 0; i < 250; i++ {
		go TestSell(&wg)
	}
	wg.Wait()
}
