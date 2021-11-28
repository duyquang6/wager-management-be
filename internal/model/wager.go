package model

import (
	"fmt"
	"math"
	"time"

	"github.com/duyquang6/wager-management-be/pkg/null"
	"github.com/duyquang6/wager-management-be/pkg/validator"
	"gorm.io/gorm"
)

// Wager data model
type Wager struct {
	BaseModel
	TotalWagerValue     uint         `gorm:"column:total_wager_value" validate:"required"`
	Odds                uint         `gorm:"column:odds" validate:"required"`
	SellingPercentage   uint         `gorm:"column:selling_percentage" validate:"gte=1,lte=100"`
	SellingPrice        float64      `gorm:"column:selling_price" validate:"required,monetary-format"`
	CurrentSellingPrice float64      `gorm:"column:current_selling_price"`
	PercentageSold      null.Uint    `gorm:"column:percentage_sold"`
	AmountSold          null.Float64 `gorm:"column:amount_sold"`
	PlacedAt            time.Time    `gorm:"column:placed_at"`
}

// BeforeSave validate Wager model
func (t *Wager) BeforeSave(tx *gorm.DB) error {
	err := t.validateModel()
	if err != nil {
		return fmt.Errorf("can't save invalid wager: %w", err)
	}
	return nil
}

func (t *Wager) validateModel() error {
	const epsilon = 1e-8
	err := validator.GetValidate().Struct(t)
	if err != nil {
		return fmt.Errorf("db field validation failed: %w", err)
	}

	if t.AmountSold.Valid && t.AmountSold.Float64 <= 0 {
		return fmt.Errorf("field AmountSold is invalid")
	}

	if t.SellingPrice <= float64(t.TotalWagerValue*t.SellingPercentage)/100 {
		return fmt.Errorf("field SellingPrice must be larger than TotalWagerValue * SellingPercentage")
	}

	// Validate consistency AmountSold + CurrentSellingPrice equal SellingPrice
	if (!t.AmountSold.Valid && math.Abs(t.CurrentSellingPrice-t.SellingPrice) > epsilon) ||
		(t.AmountSold.Valid && math.Abs(t.AmountSold.Float64+t.CurrentSellingPrice-t.SellingPrice) > epsilon) {
		return fmt.Errorf("violate consistency, AmountSold + CurrentSellingPrice not equal SellingPrice")
	}

	return nil
}
