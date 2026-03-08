package services

import (
    "context"
    "fmt"
    "sync"
    "notification-service/internal/events"
)

type Notification struct {
    UserID string `json:"user_id"`
    Type   string `json:"type"`
    Data   map[string]interface{} `json:"data"`
}

type Service struct {
    mu sync.RWMutex
    store map[string][]Notification
}

func New() *Service { return &Service{ store: map[string][]Notification{} } }

func (s *Service) Add(n Notification) {
    s.mu.Lock(); defer s.mu.Unlock()
    s.store[n.UserID] = append(s.store[n.UserID], n)
}

func (s *Service) List(userID string) []Notification {
    s.mu.RLock(); defer s.mu.RUnlock()
    return append([]Notification(nil), s.store[userID]...)
}

type Handler struct { svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{ svc: svc } }

func (h *Handler) HandleMessage(msg *events.EventMessage) error {
    var userID string
    if v, ok := msg.Metadata["user_id"].(string); ok {
        userID = v
    }
    if userID == "" {
        if v, ok := msg.Data["user_id"].(string); ok {
            userID = v
        }
    }
    if userID == "" {
        if v, ok := msg.Data["user_id"].(interface{}); ok {
            userID = fmt.Sprintf("%v", v)
        }
    }
    if userID == "" { return nil }
    h.svc.Add(Notification{ UserID: userID, Type: msg.Type, Data: map[string]interface{}{"payload": msg.Data, "meta": msg.Metadata} })
    return nil
}

func StartConsumer(ctx context.Context, topics []string, h events.MessageHandler) (*events.KafkaConsumer, error) { return events.Start(ctx, topics, h) }