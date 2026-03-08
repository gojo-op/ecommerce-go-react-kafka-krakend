package repositories

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "product-service/internal/models"
)

type ProductRepository struct {
    db *gorm.DB
}

func New(db *gorm.DB) *ProductRepository { return &ProductRepository{db: db} }

func (r *ProductRepository) Create(p *models.Product) error { return r.db.Create(p).Error }
func (r *ProductRepository) Update(p *models.Product) error { return r.db.Save(p).Error }
func (r *ProductRepository) Delete(id uuid.UUID) error { return r.db.Delete(&models.Product{}, "id = ?", id).Error }
func (r *ProductRepository) FindByID(id uuid.UUID) (*models.Product, error) {
    var p models.Product
    if err := r.db.First(&p, "id = ?", id).Error; err != nil { return nil, err }
    return &p, nil
}
func (r *ProductRepository) FindBySKU(sku string) (*models.Product, error) {
    var p models.Product
    if err := r.db.First(&p, "sku = ?", sku).Error; err != nil { return nil, err }
    return &p, nil
}
func (r *ProductRepository) List(limit, offset int) ([]models.Product, int64, error) {
    var items []models.Product
    var total int64
    if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil { return nil, 0, err }
    if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&items).Error; err != nil { return nil, 0, err }
    return items, total, nil
}