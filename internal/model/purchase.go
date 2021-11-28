package model

import (
	"fmt"
	"time"

	"github.com/duyquang6/wager-management-be/pkg/validator"
	"gorm.io/gorm"
)

// Purchase data model
type Purchase struct {
	BaseModel
	WagerID     uint      `gorm:"column:wager_id" validate:"required"`
	Wager       Wager     `gorm:"foreignKey:WagerID" validate:"-"`
	BuyingPrice float64   `gorm:"column:buying_price" validate:"required"`
	BoughtAt    time.Time `gorm:"column:bought_at" validate:"required"`
}

// BeforeSave validate Purchase model
func (t *Purchase) BeforeSave(tx *gorm.DB) error {
	err := t.validateModel()
	if err != nil {
		return fmt.Errorf("can't save invalid purchase: %w", err)
	}
	return nil
}

func (t *Purchase) validateModel() error {
	err := validator.GetValidate().Struct(t)
	if err != nil {
		return err
	}

	return nil
}
