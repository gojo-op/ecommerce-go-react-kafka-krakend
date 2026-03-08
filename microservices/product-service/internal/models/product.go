package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Product struct {
    ID          uuid.UUID      `gorm:"type:text;primaryKey"`
    Name        string         `gorm:"not null"`
    SKU         string         `gorm:"uniqueIndex;not null"`
    Description string         `gorm:"type:text"`
    Price       int64          `gorm:"not null"`
    Currency    string         `gorm:"type:varchar(10);default:'USD'"`
    Stock       int            `gorm:"not null;default:0"`
    Category    string         `gorm:"type:varchar(50)"`
    ImageURL    string         
    CreatedAt   time.Time      
    UpdatedAt   time.Time      
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type CreateProductRequest struct {
    Name        string `json:"name" binding:"required"`
    SKU         string `json:"sku" binding:"required"`
    Description string `json:"description"`
    Price       int64  `json:"price" binding:"required,min=0"`
    Currency    string `json:"currency"`
    Stock       int    `json:"stock"`
    Category    string `json:"category"`
    ImageURL    string `json:"image_url"`
}

type UpdateProductRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Price       int64  `json:"price"`
    Currency    string `json:"currency"`
    Stock       int    `json:"stock"`
    Category    string `json:"category"`
    ImageURL    string `json:"image_url"`
}