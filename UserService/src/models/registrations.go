package models

import (
	"time"
	"gorm.io/gorm"
)

type Registration struct {
	gorm.Model
	ID        int       `gorm:"column:id;primaryKey;type:INTEGER;autoIncrement;not null" json:"id"`
	ArtistID  int       `gorm:"column:artist_id;type:INTEGER;not null" json:"artist_id" validate:"required"`
	AdminID   int       `gorm:"column:admin_id;type:INTEGER;not null" json:"admin_id" validate:"required"`
	Status    string    `gorm:"column:status;type:VARCHAR(100);not null" validate:"required" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;not null" json:"created_at"`
	Admin Admin `gorm:"-" json:"-"`
}

type RegistrationRequest struct {
    ArtistID int    `json:"artist_id" validate:"required"`
    AdminID  int    `json:"admin_id" validate:"required"`
    Status   string `json:"status" validate:"required"`
}

func (r *Registration) TableName() string {
	return "registrations"
}

func MigrateRegistration(db *gorm.DB) error {
	return db.AutoMigrate(&Registration{})
}
