package dto

import (
	"github.com/duyquang6/wager-management-be/pkg/null"
)

type Wager struct {
	ID                  uint         `json:"id"`
	TotalWagerValue     uint         `json:"total_wager_value"`
	Odds                uint         `json:"odds"`
	SellingPercentage   uint         `json:"selling_percentage"`
	SellingPrice        float64      `json:"selling_price"`
	CurrentSellingPrice float64      `json:"current_selling_price"`
	PercentageSold      null.Uint    `json:"percentage_sold"`
	AmountSold          null.Float64 `json:"amount_sold"`
	PlacedAt            uint         `json:"placed_at"`
}

// ListWagersRequest ...
type ListWagersRequest struct {
	Page  uint `validate:"gte=1"`
	Limit uint `validate:"gte=1"`
}

// ListWagersResponse ...
type ListWagersResponse struct {
	Data []Wager
}

// CreateWagerRequest ...
type CreateWagerRequest struct {
	TotalWagerValue   uint    `json:"total_wager_value" validate:"gt=0"`
	Odds              uint    `json:"odds" validate:"gt=0"`
	SellingPercentage uint    `json:"selling_percentage" validate:"gte=1,lte=100"`
	SellingPrice      float64 `json:"selling_price" validate:"gt=0,monetary-format"`
}

// CreateWagerResponse ...
type CreateWagerResponse struct {
	Wager
}
