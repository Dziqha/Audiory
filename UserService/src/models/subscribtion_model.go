package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	SubscriptionStart time.Time `gorm:"column:subscription_start;type:TIMESTAMP;not null" json:"subscription_start"`
	SubscriptionEnd   time.Time `gorm:"column:subscription_end;type:TIMESTAMP;not null" json:"subscription_end"`
	SubscriptionToken string    `gorm:"column:subscription_token;type:VARCHAR(100);not null" json:"subscription_token"`
	SubscriptionType  string    `gorm:"column:subscription_type;type:subscription_type;not null" json:"subscription_type"` // Pastikan ini adalah ENUM yang sudah dibuat
	IsActive          bool      `gorm:"column:is_active;not null" json:"is_active"` // Tidak perlu mendefinisikan type:BOOLEAN di sini
	UserID 			  int       `gorm:"column:user_id;not null" json:"user_id"`  // Foreign key
	User    		  User      `gorm:"-" json:"-"` // Relasi ke model User
}

type SubscriptionRequest struct {
	SubscriptionType  string    `gorm:"column:subscription_type;type:subscription_type;not null" json:"subscription_type"`
	UserID 			  int       `gorm:"column:user_id;not null" json:"user_id"`
}


func (s *Subscription) TableName() string {
	return "subscriptions" // Nama tabel dalam huruf kecil
}

func MigrateSubscription(db *gorm.DB) error {
	return db.AutoMigrate(&Subscription{})
}
