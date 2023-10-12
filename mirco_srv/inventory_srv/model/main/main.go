package main

import (
	"MircoServer/inventory_srv/global"
	"MircoServer/inventory_srv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var err error

func main() {
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

	err = global.MysqlDB.AutoMigrate(&model.Inventory{}, &model.InventoryNew{}, &model.Delivery{}, &model.StockSellDetail{})

	for i := 0; i < 20; i++ {
		global.MysqlDB.Create(&model.Inventory{
			BaseModel: model.BaseModel{
				CreatedTime: time.Now(),
				UpdatedTime: time.Now(),
				IsDeleted:   false,
			},
			Goods:   int32(i),
			Stocks:  int32(i * 10),
			Version: 0,
		})
	}
	if err != nil {
		log.Println(">>>>>>", err)
		return
	}
}
