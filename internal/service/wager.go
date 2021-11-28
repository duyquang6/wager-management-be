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
	"time"
)

const (
	wagerServiceLoggingNameFmt = "WagerService.%s"
)

// WagerService provide wager service functionality
type WagerService interface {
	ListWagers(ctx context.Context, page, limit uint) (dto.ListWagersResponse, error)
	CreateWager(ctx context.Context, request dto.CreateWagerRequest) (dto.CreateWagerResponse, error)
}

type wagerSvc struct {
	dbFactory database.DBFactory
	wagerRepo repository.WagerRepository
}

// NewWagerService create concrete object which implement UserService
func NewWagerService(dbFactory database.DBFactory,
	wagerRepo repository.WagerRepository) WagerService {
	return &wagerSvc{
		dbFactory: dbFactory,
		wagerRepo: wagerRepo,
	}
}

// ListWagers list out wagers with pagination
func (s *wagerSvc) ListWagers(ctx context.Context, page, limit uint) (dto.ListWagersResponse, error) {
	var (
		db       = s.dbFactory.GetDB()
		function = "ListWagers"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(wagerServiceLoggingNameFmt, function))
	)

	if page == 0 || limit == 0 {
		return dto.ListWagersResponse{}, exception.New(exception.ErrPaginationParamInvalid, "wager not found")
	}
	offset := (page - 1) * limit

	wagers, err := s.wagerRepo.List(ctx, db, offset, limit)
	if err != nil {
		if database.IsNotFound(err) {
			logger.Info("wager data not found")
			return dto.ListWagersResponse{}, exception.New(exception.ErrWagerNotFound, "wager not found")
		}
		logger.Error("get wager failed:", err)
		return dto.ListWagersResponse{}, exception.Wrap(exception.ErrInternalServer, err, "list out wager failed")
	}

	return buildListWagerDtoFromModel(wagers), nil
}

// CreateWager place new wager
func (s *wagerSvc) CreateWager(ctx context.Context, request dto.CreateWagerRequest) (dto.CreateWagerResponse, error) {
	var (
		db       = s.dbFactory.GetDB()
		function = "CreateWager"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(wagerServiceLoggingNameFmt, function))
	)

	wager := model.Wager{
		TotalWagerValue:     request.TotalWagerValue,
		Odds:                request.Odds,
		SellingPercentage:   request.SellingPercentage,
		SellingPrice:        request.SellingPrice,
		CurrentSellingPrice: request.SellingPrice,
		PlacedAt:            time.Now().UTC(),
	}

	wager, err := s.wagerRepo.Create(ctx, db, wager)
	if err != nil {
		logger.Error("create wager failed:", err)
		return dto.CreateWagerResponse{}, exception.Wrap(exception.ErrInternalServer, err, "create wager failed")
	}

	return dto.CreateWagerResponse{Wager: buildWagerDtoFromModel(wager)}, nil
}

func buildListWagerDtoFromModel(wagers []model.Wager) dto.ListWagersResponse {
	var (
		data = make([]dto.Wager, 0, len(wagers))
	)

	for _, wager := range wagers {
		data = append(data, buildWagerDtoFromModel(wager))
	}

	return dto.ListWagersResponse{Data: data}
}

func buildWagerDtoFromModel(wager model.Wager) dto.Wager {
	return dto.Wager{
		ID:                  wager.ID,
		TotalWagerValue:     wager.TotalWagerValue,
		Odds:                wager.Odds,
		SellingPercentage:   wager.SellingPercentage,
		SellingPrice:        wager.SellingPrice,
		CurrentSellingPrice: wager.CurrentSellingPrice,
		PercentageSold:      wager.PercentageSold,
		AmountSold:          wager.AmountSold,
		PlacedAt:            uint(wager.PlacedAt.Unix()),
	}
}
