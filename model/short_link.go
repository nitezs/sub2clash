package model

type ShortLink struct {
	Hash            string `gorm:"primary_key"`
	Url             string
	LastRequestTime int64
}
