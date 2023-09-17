package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"sub2clash/model"
)

var DB *gorm.DB

func ConnectDB() error {
	// 用上面的数据库连接初始化 gorm
	db, err := gorm.Open(sqlite.Open("sub2clash.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err != nil {
		return err
	}
	DB = db
	err = db.AutoMigrate(&model.ShortLink{})
	if err != nil {
		return err
	}
	return nil
}
