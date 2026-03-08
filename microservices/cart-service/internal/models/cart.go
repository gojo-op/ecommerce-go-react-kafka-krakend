package models

type CartItem struct {
    SKU      string  `json:"sku"`
    Name     string  `json:"name"`
    Price    int64   `json:"price"`
    Quantity int     `json:"quantity"`
}

type Cart struct {
    UserID string     `json:"user_id"`
    Items  []CartItem `json:"items"`
}