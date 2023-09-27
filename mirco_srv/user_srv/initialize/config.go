package initialize

import (
	"MircoServer/user_srv/config"
	"MircoServer/user_srv/global"
	"MircoServer/user_srv/model"
	"fmt"
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
	isDebug        bool
	configFileName string
)

func GetEnv(env string, v *viper.Viper) bool {
	v.AutomaticEnv()
	return v.GetBool(env)
}

func InitUserSrvMysqlConfig() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.TheServerConfig.MysqlConfig.UserName,
		config.TheServerConfig.MysqlConfig.Password,
		config.TheServerConfig.MysqlConfig.Host,
		config.TheServerConfig.MysqlConfig.Port,
		config.TheServerConfig.MysqlConfig.DBName)

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

	err = global.MysqlDB.AutoMigrate(&model.User{})
	if err != nil {
		log.Println(">>>>>>", err)
		return
	}
}

func InitConfig() {

	InitLogger()

	v := viper.New()
	if isDebug = GetEnv("SHOP_ENV", v); isDebug {
		configFileName = "user_srv/config_debug.yaml"
	} else {
		configFileName = "user_srv/config_pro.yaml"
	}

	log.Println(isDebug, configFileName)
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicf(" 配置文件读取错误,err:%s", err.Error())
	}

	zap.S().Debug("name", viper.GetString("name"))

	err := v.Unmarshal(&config.TheServerConfig)
	if err != nil {
		zap.S().Panicf("配置文件反序列化错误，err:%s", err.Error())
	}

	fmt.Println(config.TheServerConfig)

	InitUserSrvMysqlConfig()
}
