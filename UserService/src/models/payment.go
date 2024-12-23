package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	ID       int              `gorm:"column:id;primaryKey;type:INTEGER,autoIncrement;not null" json:"id"`
	UserID   User             `gorm:"column:user_id;type:INTEGER;foreignKey:UserID;not null;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user_id"`
	Amount   int              `gorm:"column:amount;type:INTEGER;not null" json:"amount"`
	Payment_date time.Time    `gorm:"column:add_at;type:TIMESTAMP;not null" json:"payment_date"`
	Payment_method string     `gorm:"column:payment_method;type:ENUM('credit_card','debit_card','paypal');not null" json:"payment_method"`
	Transaction_status string `gorm:"column:transaction_status;type:ENUM('pending','success','failed');not null" json:"transaction_status"`
}

func (*Payment) TableName() string {
	return "payments"
}

func MigratePayment(db *gorm.DB) error {
	return db.AutoMigrate(&Payment{})
}