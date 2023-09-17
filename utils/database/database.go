package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"path/filepath"
	"sub2clash/model"
	"sub2clash/utils"
)

var DB *gorm.DB

func ConnectDB() error {
	// 用上面的数据库连接初始化 gorm
	err := utils.MKDir("data")
	if err != nil {
		return err
	}
	db, err := gorm.Open(sqlite.Open(filepath.Join("data", "sub2clash.db")), &gorm.Config{})
	if err != nil {
		return err
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
