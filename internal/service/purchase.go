package service

import (
	"context"
	"fmt"
	"github.com/duyquang6/wager-management-be/internal/database"
	"github.com/duyquang6/wager-management-be/internal/model"
	"github.com/duyquang6/wager-management-be/internal/repository"
	"github.com/duyquang6/wager-management-be/pkg/dto"
	"github.com/duyquang6/wager-management-be/pkg/exception"
	"github.com/duyquang6/wager-management-be/pkg/logging"
	"github.com/duyquang6/wager-management-be/pkg/null"
	"math"
	"time"
)

const (
	purchaseServiceLoggingFmt = "PurchaseService.%s"
)

// PurchaseService provide purchase service functionality
type PurchaseService interface {
	CreatePurchase(ctx context.Context, request dto.CreatePurchaseRequest) (dto.CreatePurchaseResponse, error)
}

type purchaseSvc struct {
	dbFactory    database.DBFactory
	purchaseRepo repository.PurchaseRepository
	wagerRepo    repository.WagerRepository
}

// NewPurchaseService create concrete object which implement UserService
func NewPurchaseService(dbFactory database.DBFactory,
	purchaseRepo repository.PurchaseRepository,
	wagerRepo repository.WagerRepository) PurchaseService {
	return &purchaseSvc{
		dbFactory:    dbFactory,
		purchaseRepo: purchaseRepo,
		wagerRepo:    wagerRepo,
	}
}

// CreatePurchase place new wager
func (s *purchaseSvc) CreatePurchase(ctx context.Context, request dto.CreatePurchaseRequest) (dto.CreatePurchaseResponse, error) {
	var (
		// Begin transaction, CreatePurchase need manipulate to modify Wager and Purchase at same time
		tx       = s.dbFactory.GetDBWithTx()
		function = "CreateWager"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(purchaseServiceLoggingFmt, function))
		purchase model.Purchase
	)
	defer s.dbFactory.Rollback(tx)

	wager, err := s.wagerRepo.GetByIDAndLockForUpdate(ctx, tx, request.WagerID)
	if err != nil {
		if database.IsNotFound(err) {
			logger.Infof("Related wager id %d not found", wager.ID)
			return dto.CreatePurchaseResponse{}, exception.Newf(exception.ErrRelatedWagerNotFound,
				"related wager id %d not found", request.WagerID)
		}
		logger.Error("get wager by id failed:", err)
		return dto.CreatePurchaseResponse{},
			exception.Wrap(exception.ErrInternalServer, err, "get wager failed")
	}

	if wager.CurrentSellingPrice < request.BuyingPrice {
		logger.Infof("CurrentSellingPrice %f < BuyingPrice %f", wager.CurrentSellingPrice, request.BuyingPrice)
		return dto.CreatePurchaseResponse{}, exception.New(exception.ErrBuyingPriceGreaterThanCurrentSellingPrice,
			"buying price must be smaller or equal current selling price")
	}

	wager.CurrentSellingPrice -= request.BuyingPrice
	wager.AmountSold.Valid = true
	wager.AmountSold.Float64 += request.BuyingPrice
	wager.PercentageSold = null.Uint{
		Uint:  uint(math.Round(wager.AmountSold.Float64 / wager.SellingPrice * 100)),
		Valid: true,
	}

	wager, err = s.wagerRepo.Update(ctx, tx, wager)
	if err != nil {
		logger.Error("update wager failed:", err)
		return dto.CreatePurchaseResponse{}, exception.Wrap(exception.ErrInternalServer, err, "update wager failed")
	}

	purchase, err = s.purchaseRepo.Create(ctx, tx, model.Purchase{
		WagerID:     wager.ID,
		Wager:       wager,
		BuyingPrice: request.BuyingPrice,
		BoughtAt:    time.Now().UTC(),
	})

	if err != nil {
		logger.Error("create purchase failed:", err)
		return dto.CreatePurchaseResponse{}, exception.Wrap(exception.ErrInternalServer, err, "create purchase failed")
	}

	s.dbFactory.Commit(tx)

	return dto.CreatePurchaseResponse{
		ID:          purchase.ID,
		BuyingPrice: purchase.BuyingPrice,
		WagerID:     purchase.WagerID,
		BoughtAt:    uint(purchase.BoughtAt.Unix()),
	}, nil
}
