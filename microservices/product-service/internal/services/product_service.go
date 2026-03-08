package services

import (
    "context"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "product-service/internal/models"
    "product-service/internal/repositories"
    "product-service/internal/events"
)

type Service struct {
    repo    *repositories.ProductRepository
    kafka   *events.Publisher
}

func New(db *gorm.DB, kafka *events.Publisher) *Service {
    return &Service{repo: repositories.New(db), kafka: kafka}
}

func (s *Service) Create(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
    p := &models.Product{
        ID: uuid.New(),
        Name: req.Name,
        SKU: req.SKU,
        Description: req.Description,
        Price: req.Price,
        Currency: req.Currency,
        Stock: req.Stock,
        Category: req.Category,
        ImageURL: req.ImageURL,
    }
    if err := s.repo.Create(p); err != nil { return nil, err }
    _ = s.kafka.Publish(ctx, "product.created", map[string]interface{}{"id": p.ID, "sku": p.SKU, "name": p.Name})
    return p, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error) {
    p, err := s.repo.FindByID(id)
    if err != nil { return nil, err }
    if req.Name != "" { p.Name = req.Name }
    if req.Description != "" { p.Description = req.Description }
    if req.Price != 0 { p.Price = req.Price }
    if req.Currency != "" { p.Currency = req.Currency }
    if req.Stock != 0 { p.Stock = req.Stock }
    if req.Category != "" { p.Category = req.Category }
    if req.ImageURL != "" { p.ImageURL = req.ImageURL }
    if err := s.repo.Update(p); err != nil { return nil, err }
    _ = s.kafka.Publish(ctx, "product.updated", map[string]interface{}{"id": p.ID, "sku": p.SKU, "name": p.Name})
    return p, nil
}

func (s *Service) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) (*models.Product, error) {
    p, err := s.repo.FindByID(id)
    if err != nil { return nil, err }
    p.Stock = quantity
    if err := s.repo.Update(p); err != nil { return nil, err }
    _ = s.kafka.Publish(ctx, "product.updated", map[string]interface{}{"id": p.ID, "sku": p.SKU, "name": p.Name, "stock": p.Stock})
    return p, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
    p, err := s.repo.FindByID(id)
    if err != nil { return err }
    if err := s.repo.Delete(id); err != nil { return err }
    _ = s.kafka.Publish(ctx, "product.deleted", map[string]interface{}{"id": p.ID, "sku": p.SKU})
    return nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.Product, error) {
    return s.repo.FindByID(id)
}

func (s *Service) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
    return s.repo.FindBySKU(sku)
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]models.Product, int64, error) {
    items, total, err := s.repo.List(limit, offset)
    if err != nil { return nil, 0, err }
    return items, total, nil
}