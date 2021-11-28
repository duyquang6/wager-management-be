package repository

import (
	"context"
	"gorm.io/gorm/clause"

	"github.com/duyquang6/wager-management-be/internal/model"
	"gorm.io/gorm"
)

type wagerRepo struct{}

// WagerRepository provide interface interact with Wager model
type WagerRepository interface {
	Create(ctx context.Context, tx *gorm.DB, wager model.Wager) (model.Wager, error)
	Update(ctx context.Context, tx *gorm.DB, wager model.Wager) (model.Wager, error)
	List(ctx context.Context, tx *gorm.DB, offset, limit uint) ([]model.Wager, error)
	GetByID(ctx context.Context, tx *gorm.DB, id uint) (model.Wager, error)
	GetByIDAndLockForUpdate(ctx context.Context, tx *gorm.DB, id uint) (model.Wager, error)
}

// NewWagerRepository create WagerRepository concrete object
func NewWagerRepository() WagerRepository {
	return &wagerRepo{}
}

// Create new wager
func (s *wagerRepo) Create(ctx context.Context, tx *gorm.DB, wager model.Wager) (model.Wager, error) {
	res := tx.Create(&wager)
	return wager, res.Error
}

// Update existed wager
func (s *wagerRepo) Update(ctx context.Context, tx *gorm.DB, wager model.Wager) (model.Wager, error) {
	res := tx.Model(&wager).Select("CurrentSellingPrice", "PercentageSold", "AmountSold").Updates(&wager)
	return wager, res.Error
}

// List wager data with pagination support
func (s *wagerRepo) List(ctx context.Context, tx *gorm.DB, offset, limit uint) ([]model.Wager, error) {
	var wagers []model.Wager
	res := tx.Limit(int(limit)).Offset(int(offset)).Find(&wagers)
	return wagers, res.Error
}

// GetByID wager
func (s *wagerRepo) GetByID(ctx context.Context, tx *gorm.DB, id uint) (model.Wager, error) {
	wager := model.Wager{}
	wager.ID = id
	res := tx.First(&wager)
	return wager, res.Error
}

// GetByIDAndLockForUpdate get wager and lock for update
func (s *wagerRepo) GetByIDAndLockForUpdate(ctx context.Context, tx *gorm.DB, id uint) (model.Wager, error) {
	wager := model.Wager{}
	wager.ID = id
	res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wager)
	return wager, res.Error
}
