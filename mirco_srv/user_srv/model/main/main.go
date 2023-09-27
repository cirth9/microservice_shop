package main

import (
	"MircoServer/user_srv/global"
	"MircoServer/user_srv/handler"
	"MircoServer/user_srv/model"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

var err error

func main() {
	dsn := "root:qq31415926535--@tcp(127.0.0.1:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
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
	salt, encodedPwd := password.Encode("Test", handler.Options)
	Password := fmt.Sprintf("pbkdf2-sha256$%s$%s", salt, encodedPwd)
	passwordInfo := strings.Split(Password, "$")
	fmt.Printf("%+v", passwordInfo)

	for i := 0; i < 10; i++ {
		nowTime := time.Now()
		user := model.User{
			BaseModel: model.BaseModel{
				CreatedTime: nowTime,
				UpdatedTime: nowTime,
				DeletedTime: gorm.DeletedAt{},
				IsDeleted:   false,
			},
			MobilePhoneNumber: fmt.Sprintf("12345%d", i),
			Password:          Password,
			NickName:          fmt.Sprintf("sadhjnz%d", i),
			BirthDay:          &nowTime,
		}
		global.MysqlDB.Save(&user)
	}
	//for _, v := range passwordInfo {
	//	fmt.Println(v, len(v))
	//}
	//println(password.Verify("tes12sfasdasdafasdz312t", salt, encodedPwd, options))
}
