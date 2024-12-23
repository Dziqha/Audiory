package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	ID int `gorm:"column:id;primaryKey;type:INTEGER;autoIncrement;not null" json:"id"`
	Username string `gorm:"column:username;type:VARCHAR(100);not null"  json:"username" validate:"required"`
	Email string `gorm:"column:email;type:VARCHAR(100);not null"  json:"email" validate:"required"`
	Password string `gorm:"column:password;type:VARCHAR(100);NOT NULL"  json:"password" validate:"required"`
}


type AdminRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (a *Admin) TableName() string {
	return "Admin"
}

func MigrateAdmin(db *gorm.DB) error {
	return db.AutoMigrate(&Admin{})
}