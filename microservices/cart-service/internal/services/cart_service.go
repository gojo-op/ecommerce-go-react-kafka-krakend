package services

import (
    "context"
    "errors"
    "cart-service/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "cart-service/internal/events"
)

type Service struct {
    db    *gorm.DB
    kafka *events.Publisher
}

func New(db *gorm.DB, kafka *events.Publisher) *Service { return &Service{db: db, kafka: kafka} }

func (s *Service) Get(ctx context.Context, userID string) (*models.Cart, error) {
    var cartE models.CartEntity
    if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&cartE).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &models.Cart{UserID: userID, Items: []models.CartItem{}}, nil
        }
        return nil, err
    }
    var itemsE []models.CartItemEntity
    if err := s.db.WithContext(ctx).Where("cart_id = ?", cartE.ID).Find(&itemsE).Error; err != nil { return nil, err }
    items := make([]models.CartItem, 0, len(itemsE))
    for _, it := range itemsE { items = append(items, models.CartItem{SKU: it.SKU, Name: it.Name, Price: it.Price, Quantity: it.Quantity}) }
    return &models.Cart{UserID: userID, Items: items}, nil
}

func (s *Service) AddItem(ctx context.Context, userID string, item models.CartItem) (*models.Cart, error) {
    var cartE models.CartEntity
    if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&cartE).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            cartE = models.CartEntity{ID: uuid.New(), UserID: userID}
            if err := s.db.WithContext(ctx).Create(&cartE).Error; err != nil { return nil, err }
        } else { return nil, err }
    }
    var itemE models.CartItemEntity
    if err := s.db.WithContext(ctx).Where("cart_id = ? AND sku = ?", cartE.ID, item.SKU).First(&itemE).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            itemE = models.CartItemEntity{ID: uuid.New(), CartID: cartE.ID, SKU: item.SKU, Name: item.Name, Price: item.Price, Quantity: item.Quantity}
            if err := s.db.WithContext(ctx).Create(&itemE).Error; err != nil { return nil, err }
        } else { return nil, err }
    } else {
        itemE.Quantity += item.Quantity
        if err := s.db.WithContext(ctx).Save(&itemE).Error; err != nil { return nil, err }
    }
    _ = s.kafka.Publish(ctx, "cart.item_added", map[string]interface{}{"user_id": userID, "sku": item.SKU, "qty": item.Quantity})
    return s.Get(ctx, userID)
}

func (s *Service) RemoveItem(ctx context.Context, userID, sku string) (*models.Cart, error) {
    var cartE models.CartEntity
    if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&cartE).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return &models.Cart{UserID: userID, Items: []models.CartItem{}}, nil }
        return nil, err
    }
    _ = s.db.WithContext(ctx).Where("cart_id = ? AND sku = ?", cartE.ID, sku).Delete(&models.CartItemEntity{})
    _ = s.kafka.Publish(ctx, "cart.item_removed", map[string]interface{}{"user_id": userID, "sku": sku})
    return s.Get(ctx, userID)
}

func (s *Service) Clear(ctx context.Context, userID string) error {
    var cartE models.CartEntity
    if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&cartE).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil }
        return err
    }
    _ = s.db.WithContext(ctx).Where("cart_id = ?", cartE.ID).Delete(&models.CartItemEntity{})
    _ = s.kafka.Publish(ctx, "cart.cleared", map[string]interface{}{"user_id": userID})
    return nil
}