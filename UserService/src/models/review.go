package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	ID int `gorm:"column:id;primaryKey;type:INTEGER;autoIncrement;not null" json:"id"`
	UserID int `gorm:"column:user_id;type:INTEGER;not null" json:"user_id"`
	User   User `gorm:"-" json:"-"`
	AlbumID int `gorm:"column:album_id;type:INTEGER;not null" json:"album_id"`
	Rating string `gorm:"column:rating;type:VARCHAR(256);not null" json:"rating"`
	Comment string `gorm:"column:comment;type:VARCHAR(256);not null" json:"comment"`
	CreatedAt time.Time`gorm:"column:created_at;type:TIMESTAMP;not null" json:"created_at"`
}

type ReviewRequest struct {
	UserID int `json:"user_id" validate:"required"`
	AlbumID int `json:"album_id" validate:"required"`
	Rating string `json:"rating" validate:"required"`
	Comment string `json:"comment" validate:"required"`
}


func (s *Review) TabelName() string {
	return "Review"
}


func MigrateReview(db *gorm.DB) error {
	return db.AutoMigrate(&Review{})
}