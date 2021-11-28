package repository

import (
	"context"

	"github.com/duyquang6/wager-management-be/internal/model"
	"gorm.io/gorm"
)

type purchaseRepo struct{}

// PurchaseRepository provide interface interact with Purchase model
type PurchaseRepository interface {
	Create(ctx context.Context, tx *gorm.DB, purchase model.Purchase) (model.Purchase, error)
}

// NewPurchaseRepository create PurchaseRepository concrete object
func NewPurchaseRepository() PurchaseRepository {
	return &purchaseRepo{}
}

// Create new purchase with given Purchase model
func (s *purchaseRepo) Create(ctx context.Context, tx *gorm.DB, purchase model.Purchase) (model.Purchase, error) {
	res := tx.Select("BuyingPrice", "WagerID", "BoughtAt").Create(&purchase)
	return purchase, res.Error
}
