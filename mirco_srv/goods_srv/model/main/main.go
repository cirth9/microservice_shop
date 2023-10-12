package main

import (
	"MircoServer/goods_srv/global"
	"MircoServer/goods_srv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var err error

func main() {
	dsn := "root:qq31415926535--@tcp(127.0.0.1:3306)/shop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
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

	err = global.MysqlDB.AutoMigrate(&model.Goods{}, &model.Brands{}, &model.GoodsCategoryBrand{}, &model.Category{}, &model.Banner{})
	if err != nil {
		log.Println(">>>>>>", err)
		return
	}
}
