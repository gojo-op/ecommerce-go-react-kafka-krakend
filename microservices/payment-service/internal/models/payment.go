package models

type CreateIntentRequest struct {
    OrderID   string `json:"order_id" binding:"required"`
    Amount    int64  `json:"amount" binding:"required,min=1"`
    Currency  string `json:"currency" binding:"required"`
    Provider  string `json:"provider" binding:"required,oneof=stripe razorpay"`
}

type PaymentIntent struct {
    IntentID string `json:"intent_id"`
    OrderID  string `json:"order_id"`
    Amount   int64  `json:"amount"`
    Currency string `json:"currency"`
    Provider string `json:"provider"`
    Status   string `json:"status"`
}

type WebhookEvent struct {
    Provider string                 `json:"provider"`
    Type     string                 `json:"type"`
    Data     map[string]interface{} `json:"data"`
}