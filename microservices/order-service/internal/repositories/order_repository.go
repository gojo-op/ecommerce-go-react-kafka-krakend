package repositories

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "order-service/internal/models"
)

type Repository struct { db *gorm.DB }
func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(o *models.Order) error { return r.db.Create(o).Error }
func (r *Repository) Update(o *models.Order) error { return r.db.Save(o).Error }
func (r *Repository) FindByID(id uuid.UUID) (*models.Order, error) {
    var o models.Order
    if err := r.db.Preload("Items").First(&o, "id = ?", id).Error; err != nil { return nil, err }
    return &o, nil
}
func (r *Repository) ListByUser(userID uuid.UUID, limit, offset int) ([]models.Order, int64, error) {
    var items []models.Order
    var total int64
    q := r.db.Where("user_id = ?", userID).Model(&models.Order{})
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    if err := r.db.Where("user_id = ?", userID).Preload("Items").Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil { return nil, 0, err }
    return items, total, nil
}