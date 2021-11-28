// +build integration

package integration

import (
	"context"
	"github.com/duyquang6/wager-management-be/internal/model"
	"github.com/duyquang6/wager-management-be/internal/repository"
	"time"
)

func (p *MySqlRepositoryTestSuite) TestMySqlPurchaseRepository_Create() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	purchaseRepo := repository.NewPurchaseRepository()
	wagerRepo := repository.NewWagerRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		wager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		purchase := model.Purchase{
			WagerID:     1,
			Wager:       wager,
			BuyingPrice: 1,
			BoughtAt:    time.Now().UTC(),
		}
		_, err := wagerRepo.Create(ctx, tx, wager)
		p.Assert().NoError(err)
		purchase, err = purchaseRepo.Create(ctx, tx, purchase)
		p.Assert().NoError(err)
		actualPurchase := model.Purchase{}
		actualPurchase.ID = purchase.ID
		err = tx.First(&actualPurchase).Error
		p.Assert().NoError(err)
		p.Assert().Equal(float64(1), actualPurchase.BuyingPrice)
	})

	p.Run("Failed because of missing wager id", func() {
		tx := db.Begin()
		defer tx.Rollback()
		wager := model.Wager{
			BaseModel:           model.BaseModel{ID: 0},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		purchase := model.Purchase{
			WagerID:     0,
			Wager:       wager,
			BuyingPrice: 0.25,
			BoughtAt:    time.Now().UTC(),
		}
		_, err := wagerRepo.Create(ctx, tx, wager)
		p.Assert().NoError(err)
		purchase, err = purchaseRepo.Create(ctx, tx, purchase)
		p.Assert().Error(err)
	})

	p.Run("Failed because of missing buying price", func() {
		tx := db.Begin()
		defer tx.Rollback()
		wager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		purchase := model.Purchase{
			WagerID:     1,
			Wager:       wager,
			BuyingPrice: 0,
		}
		_, err := wagerRepo.Create(ctx, tx, wager)
		p.Assert().NoError(err)
		purchase, err = purchaseRepo.Create(ctx, tx, purchase)
		p.Assert().Error(err)
	})

	p.Run("Failed because of missing bought at", func() {
		tx := db.Begin()
		defer tx.Rollback()
		wager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		purchase := model.Purchase{
			WagerID:     1,
			Wager:       wager,
			BuyingPrice: 0.25,
		}
		_, err := wagerRepo.Create(ctx, tx, wager)
		p.Assert().NoError(err)
		purchase, err = purchaseRepo.Create(ctx, tx, purchase)
		p.Assert().Error(err)
	})

	p.Run("Failed because of related wager not found", func() {
		tx := db.Begin()
		defer tx.Rollback()
		wager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		purchase := model.Purchase{
			WagerID:     1,
			Wager:       wager,
			BuyingPrice: 0.25,
			BoughtAt:    time.Now().UTC(),
		}
		purchase, err := purchaseRepo.Create(ctx, tx, purchase)
		p.Assert().Error(err)
	})
}

