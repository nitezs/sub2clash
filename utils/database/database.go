package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sub2clash/model"
)

var DB *gorm.DB

func ConnectDB() error {
	db, err := gorm.Open(sqlite.Open("sub2clash.db"), &gorm.Config{})
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
