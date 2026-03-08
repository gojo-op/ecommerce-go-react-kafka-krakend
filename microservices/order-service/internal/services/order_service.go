package services

import (
    "context"
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "order-service/internal/models"
    "order-service/internal/repositories"
    "order-service/internal/events"
    "strings"
)

type Service struct {
    repo  *repositories.Repository
    kafka *events.Publisher
}

func New(db *gorm.DB, kafka *events.Publisher) *Service { return &Service{repo: repositories.New(db), kafka: kafka} }

func (s *Service) Create(ctx context.Context, req *models.CreateOrderRequest) (*models.Order, error) {
    userID, err := uuid.Parse(req.UserID)
    if err != nil { return nil, err }
    var total int64
    items := make([]models.OrderItem, 0, len(req.Items))
    for _, it := range req.Items {
        t := it.UnitPrice * int64(it.Quantity)
        items = append(items, models.OrderItem{ID: uuid.New(), SKU: it.SKU, Name: it.Name, UnitPrice: it.UnitPrice, Quantity: it.Quantity, Total: t})
        total += t
    }
    o := &models.Order{ID: uuid.New(), UserID: userID, Status: "pending", Total: total, Currency: req.Currency, Items: items}
    if err := s.repo.Create(o); err != nil { return nil, err }
    _ = s.kafka.Publish(ctx, "order.created", map[string]interface{}{"order_id": o.ID, "user_id": o.UserID, "total": o.Total})
    return o, nil
}

func (s *Service) Checkout(ctx context.Context, req *models.CheckoutRequest) (*models.Order, error) {
    userID, err := uuid.Parse(req.UserID)
    if err != nil { return nil, err }
    var total int64
    items := make([]models.OrderItem, 0, len(req.Items))
    for _, it := range req.Items {
        t := it.UnitPrice * int64(it.Quantity)
        items = append(items, models.OrderItem{ID: uuid.New(), SKU: it.SKU, Name: it.Name, UnitPrice: it.UnitPrice, Quantity: it.Quantity, Total: t})
        total += t
    }
    pm := strings.ToLower(strings.TrimSpace(req.PaymentMethod))
    if pm == "" { pm = "online" }
    status := "pending"
    if pm == "cod" { status = "processing" }
    o := &models.Order{ID: uuid.New(), UserID: userID, Status: status, Total: total, Currency: req.Currency, PaymentMethod: pm, Items: items,
        ShipName: req.Shipping.Name, ShipPhone: req.Shipping.Phone, ShipAddr1: req.Shipping.Address1, ShipAddr2: req.Shipping.Address2,
        ShipCity: req.Shipping.City, ShipState: req.Shipping.State, ShipCountry: req.Shipping.Country, ShipPostal: req.Shipping.Postal,
    }
    if err := s.repo.Create(o); err != nil { return nil, err }
    _ = s.kafka.Publish(ctx, "order.created", map[string]interface{}{"order_id": o.ID, "user_id": o.UserID, "total": o.Total, "payment_method": pm})
    if pm == "cod" { _ = s.kafka.Publish(ctx, "order.status_changed", map[string]interface{}{"order_id": o.ID, "status": "processing"}) }
    return o, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*models.Order, error) {
    o, err := s.repo.FindByID(id)
    if err != nil { return nil, err }
    o.Status = status
    if err := s.repo.Update(o); err != nil { return nil, err }
    _ = s.kafka.Publish(ctx, "order.status_changed", map[string]interface{}{"order_id": o.ID, "status": o.Status, "updated_at": time.Now()})
    return o, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.Order, error) { return s.repo.FindByID(id) }
func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Order, int64, error) { return s.repo.ListByUser(userID, limit, offset) }