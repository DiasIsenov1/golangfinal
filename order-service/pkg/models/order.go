package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID        int64          `gorm:"primary_key"`
	ProductID int64          `gorm:"not null"`
	UserID    int64          `gorm:"not null"`
	Quantity  int64          `gorm:"not null"`
	Price     int64          `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
