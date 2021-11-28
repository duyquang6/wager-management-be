package purchase

import (
	"github.com/duyquang6/wager-management-be/internal/service"
)

type Controller struct {
	service service.PurchaseService
}

// NewController creates a new controller.
func NewController(purchaseService service.PurchaseService) *Controller {
	return &Controller{
		service: purchaseService,
	}
}
