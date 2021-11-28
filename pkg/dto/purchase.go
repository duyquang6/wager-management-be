package dto

// CreatePurchaseRequest ...
type CreatePurchaseRequest struct {
	BuyingPrice float64 `json:"buying_price" validate:"gt=0"`
	WagerID     uint    `json:"wager_id" validate:"required"`
}

// CreatePurchaseResponse ...
type CreatePurchaseResponse struct {
	ID          uint    `json:"id"`
	BuyingPrice float64 `json:"buying_price"`
	WagerID     uint    `json:"wager_id"`
	BoughtAt    uint    `json:"bought_at"`
}
