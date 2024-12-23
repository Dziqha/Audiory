package models

import (
	"time"

	"gorm.io/gorm"
)

type ListeningHistory struct {
	gorm.Model
	ID             int       `gorm:"column:id;type:INTEGER;autoIncrement; primaryKey; not null" json:"id"`
	PlayedAt       time.Time `gorm:"column:played_at;type:TIMESTAMP;not null" json:"played_at"`
	DurationPlayed int       `gorm:"column:duration_played;type:INTEGER;not null" json:"duration_played"`
	UserID         int       `gorm:"column:user_id;type:INTEGER;not null" json:"user_id"`
	SongID         int       `gorm:"column:song_id;type:INTEGER;not null" json:"song_id"`
	User           User      `gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
}

func (l *ListeningHistory) TableName() string {
	return "listening_history"
}

func MigrateListeningHistory(db *gorm.DB) error {
	return db.AutoMigrate(&ListeningHistory{})
}
