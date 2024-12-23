package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int       `gorm:"column:id;primaryKey;type:INTEGER;not null" json:"id"`
	Username string    `gorm:"column:username;type:VARCHAR(100);not null" json:"username" validate:"required"`
	Password string    `gorm:"column:password;type:VARCHAR(100);not null" json:"password" validate:"required"`
	Email    string    `gorm:"column:email;type:VARCHAR(100);not null" json:"email" validate:"required"`
}

type UserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}


func (u *User) TableName() string {
	return "Users"
}

func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}