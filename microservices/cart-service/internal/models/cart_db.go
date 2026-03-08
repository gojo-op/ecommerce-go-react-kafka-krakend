package models

import (
    "github.com/google/uuid"
    "time"
    "gorm.io/gorm"
)

type CartEntity struct {
    ID        uuid.UUID      `gorm:"type:text;primaryKey"`
    UserID    string         `gorm:"index;unique;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Items     []CartItemEntity `gorm:"foreignKey:CartID"`
}

type CartItemEntity struct {
    ID        uuid.UUID `gorm:"type:text;primaryKey"`
    CartID    uuid.UUID `gorm:"index;not null"`
    SKU       string    `gorm:"index;not null"`
    Name      string    `gorm:"not null"`
    Price     int64     `gorm:"not null"`
    Quantity  int       `gorm:"not null;default:1"`
}