package models

import (
	"time"

	"gorm.io/gorm"
)

type FavoriteSong struct {
	gorm.Model
	ID      int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID  int       `gorm:"column:user_id;not null" json:"user_id"` 
	User    User      `gorm:"-" json:"-"`
	SongID  int       `gorm:"column:song_id;not null" json:"song_id"`
	GenreID int       `gorm:"column:genre_id;not null" json:"genre_id"`
	AddedAt time.Time `gorm:"column:add_at;not null" json:"add_at"`
}

type FavoriteSongRequest struct {
	UserID  int `json:"user_id"`
	SongID  int `json:"song_id"`
}

func (f *FavoriteSong) TableName() string{
	return "favorite_songs"
}


func MigrateFavoriteSong(db *gorm.DB) error {
	return db.AutoMigrate(&FavoriteSong{})
}