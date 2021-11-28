package wager

import (
	"github.com/duyquang6/wager-management-be/internal/service"
)

type Controller struct {
	service service.WagerService
}

// NewController creates a new controller.
func NewController(wagerService service.WagerService) *Controller {
	return &Controller{
		service: wagerService,
	}
}
