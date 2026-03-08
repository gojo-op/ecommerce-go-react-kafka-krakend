package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Order struct {
    ID        uuid.UUID      `gorm:"type:text;primaryKey"`
    UserID    uuid.UUID      `gorm:"type:text;index;not null"`
    Status    string         `gorm:"type:varchar(20);default:'pending'"`
    Total     int64          `gorm:"not null;default:0"`
    Currency  string         `gorm:"type:varchar(10);default:'USD'"`
    PaymentMethod string     `gorm:"type:varchar(20);default:'online'"`
    // Shipping snapshot
    ShipName   string        `gorm:"type:varchar(100)"`
    ShipPhone  string        `gorm:"type:varchar(30)"`
    ShipAddr1  string        `gorm:"type:text"`
    ShipAddr2  string        `gorm:"type:text"`
    ShipCity   string        `gorm:"type:varchar(50)"`
    ShipState  string        `gorm:"type:varchar(50)"`
    ShipCountry string       `gorm:"type:varchar(50)"`
    ShipPostal string        `gorm:"type:varchar(20)"`
    CreatedAt time.Time      
    UpdatedAt time.Time      
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Items     []OrderItem    `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
    ID        uuid.UUID `gorm:"type:text;primaryKey"`
    OrderID   uuid.UUID `gorm:"type:text;index;not null"`
    SKU       string    `gorm:"not null"`
    Name      string    `gorm:"not null"`
    UnitPrice int64     `gorm:"not null"`
    Quantity  int       `gorm:"not null;default:1"`
    Total     int64     `gorm:"not null;default:0"`
}

type CreateOrderRequest struct {
    UserID string       `json:"user_id" binding:"required"`
    Items  []OrderItemRequest `json:"items" binding:"required,min=1"`
    Currency string     `json:"currency"`
}

type OrderItemRequest struct {
    SKU       string `json:"sku" binding:"required"`
    Name      string `json:"name" binding:"required"`
    UnitPrice int64  `json:"unit_price" binding:"required,min=0"`
    Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type UpdateStatusRequest struct {
    Status string `json:"status" binding:"required"`
}

type CheckoutRequest struct {
    UserID string `json:"user_id" binding:"required"`
    Items  []OrderItemRequest `json:"items" binding:"required,min=1"`
    Currency string `json:"currency"`
    PaymentMethod string `json:"payment_method"`
    Shipping struct {
        Name string `json:"name" binding:"required"`
        Phone string `json:"phone"`
        Address1 string `json:"address1" binding:"required"`
        Address2 string `json:"address2"`
        City string `json:"city" binding:"required"`
        State string `json:"state" binding:"required"`
        Country string `json:"country" binding:"required"`
        Postal string `json:"postal" binding:"required"`
    } `json:"shipping" binding:"required"`
}