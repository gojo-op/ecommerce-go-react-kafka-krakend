package services

import (
    "context"
    "time"
    "github.com/google/uuid"
    "payment-service/internal/models"
    "payment-service/internal/events"
)

type Service struct {
    kafka *events.Publisher
}

func New(kafka *events.Publisher) *Service { return &Service{kafka: kafka} }

func (s *Service) CreateIntent(ctx context.Context, req *models.CreateIntentRequest) (*models.PaymentIntent, error) {
    id := uuid.New().String()
    intent := &models.PaymentIntent{IntentID: id, OrderID: req.OrderID, Amount: req.Amount, Currency: req.Currency, Provider: req.Provider, Status: "requires_payment"}
    _ = s.kafka.Publish(ctx, "payment.processed", map[string]interface{}{"intent_id": id, "order_id": req.OrderID, "provider": req.Provider, "status": intent.Status})
    return intent, nil
}

func (s *Service) HandleWebhook(ctx context.Context, ev *models.WebhookEvent) error {
    t := time.Now()
    if ev.Type == "payment.succeeded" {
        _ = s.kafka.Publish(ctx, "payment.processed", map[string]interface{}{"provider": ev.Provider, "data": ev.Data, "ts": t})
        return nil
    }
    if ev.Type == "payment.failed" {
        _ = s.kafka.Publish(ctx, "payment.failed", map[string]interface{}{"provider": ev.Provider, "data": ev.Data, "ts": t})
        return nil
    }
    if ev.Type == "payment.refunded" {
        _ = s.kafka.Publish(ctx, "payment.refunded", map[string]interface{}{"provider": ev.Provider, "data": ev.Data, "ts": t})
        return nil
    }
    return nil
}