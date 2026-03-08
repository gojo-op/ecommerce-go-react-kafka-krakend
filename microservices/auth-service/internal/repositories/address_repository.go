package repositories

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "auth-service/internal/models"
)

type AddressRepository struct { db *gorm.DB }
func NewAddressRepository(db *gorm.DB) *AddressRepository { return &AddressRepository{ db: db } }

func (r *AddressRepository) ListByUser(userID uuid.UUID) ([]models.Address, error) {
    var items []models.Address
    if err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&items).Error; err != nil { return nil, err }
    return items, nil
}

func (r *AddressRepository) Create(a *models.Address) error {
    if a.ID == uuid.Nil { a.ID = uuid.New() }
    if a.IsDefault {
        _ = r.db.Model(&models.Address{}).Where("user_id = ?", a.UserID).Update("is_default", false).Error
    }
    return r.db.Create(a).Error
}

func (r *AddressRepository) Update(a *models.Address) error {
    if a.IsDefault {
        _ = r.db.Model(&models.Address{}).Where("user_id = ?", a.UserID).Update("is_default", false).Error
    }
    return r.db.Save(a).Error
}

func (r *AddressRepository) Delete(userID, id uuid.UUID) error {
    return r.db.Where("user_id = ? AND id = ?", userID, id).Delete(&models.Address{}).Error
}

func (r *AddressRepository) FindByID(userID, id uuid.UUID) (*models.Address, error) {
    var a models.Address
    if err := r.db.Where("user_id = ? AND id = ?", userID, id).First(&a).Error; err != nil { return nil, err }
    return &a, nil
}