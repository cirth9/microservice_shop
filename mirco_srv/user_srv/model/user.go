package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID          int32          `gorm:"primaryKey"`
	CreatedTime time.Time      `gorm:"column:add_time"`
	UpdatedTime time.Time      `gorm:"column:update_time"`
	DeletedTime gorm.DeletedAt `gorm:"column:delete_time"`
	IsDeleted   bool           `gorm:"column:is_deleted"`
}

type User struct {
	BaseModel
	MobilePhoneNumber string     `gorm:"index:idx_mobile;unique;type:varchar(20);not null"`
	Password          string     `gorm:"column:password;type:varchar(200);not null"`
	NickName          string     `gorm:"column:nick_name;type:varchar(20);not null"`
	BirthDay          *time.Time `gorm:"column:birthday;type:datetime"`
	Gender            string     `gorm:"column:gender;default:male;type:varchar(20)"` //male or female
	Role              int32      `gorm:"column:role;default:1;type:int"`              //1 normal 2 admin
}
