package database

import (
	"path/filepath"
	"sub2clash/logger"
	"sub2clash/model"
	"sub2clash/utils"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	// 用上面的数据库连接初始化 gorm
	err := utils.MKDir("data")
	if err != nil {
		return err
	}
	db, err := gorm.Open(
		sqlite.Open(filepath.Join("data", "sub2clash.db")), &gorm.Config{
			Logger: nil,
		},
	)
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

func FindShortLinkByUrl(url string, shortLink *model.ShortLink) *gorm.DB {
	logger.Logger.Debug("find short link by url", zap.String("url", url))
	return DB.Where("url = ?", url).First(&shortLink)
}

func FindShortLinkByHash(hash string, shortLink *model.ShortLink) *gorm.DB {
	logger.Logger.Debug("find short link by hash", zap.String("hash", hash))
	return DB.Where("hash = ?", hash).First(&shortLink)
}

func SaveShortLink(shortLink *model.ShortLink) {
	logger.Logger.Debug("save short link", zap.String("hash", shortLink.Hash))
	DB.Save(shortLink)
}

func FirstOrCreateShortLink(shortLink *model.ShortLink) {
	logger.Logger.Debug("first or create short link", zap.String("hash", shortLink.Hash))
	DB.FirstOrCreate(shortLink)
}
