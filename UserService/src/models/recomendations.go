package models

import (
	"time"

	"gorm.io/gorm"
)

type Recomendation struct {
	gorm.Model
	ID     int 			 `gorm:"column:id;primaryKey;type:INTEGER,autoIncrement;not null" json:"id"`
	SongID int 			 `gorm:"column:song_id;type:INTEGER;foreignKey:SongID;not null;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"song_id"`
	UserID int 		     `gorm:"column:user_id;type:INTEGER;not null" json:"user_id"`
	User   User 		 `gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
	GenreID int 		 `gorm:"column:genre_id;type:INTEGER;foreignKey:GenreID;not null;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"genre_id"`
	Created_at time.Time `gorm:"column:created_at;type:TIMESTAMP;not null" json:"created_at"`
	Reason string 		 `gorm:"column:reason;type:TEXT;not null" json:"reason"`
}

func (r *Recomendation) TableName() string {
	return "recomendations"
}

func MigrateRecomendation(db *gorm.DB) error {
	return db.AutoMigrate(&Recomendation{})
}