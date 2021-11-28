// +build integration

package integration

import (
	"context"
	"github.com/duyquang6/wager-management-be/internal/model"
	"github.com/duyquang6/wager-management-be/internal/repository"
	"github.com/duyquang6/wager-management-be/pkg/null"
	"time"
)

func (p *MySqlRepositoryTestSuite) TestMySqlWagerRepository_CreateWager() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	r := repository.NewWagerRepository()

	p.Run("Failed because of invalid SellingPrice", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        0,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().Error(err)
	})

	p.Run("Failed because of SellingPrice must be larger than Total x Percentage", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     2,
			Odds:                1,
			SellingPercentage:   50,
			SellingPrice:        1,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().Error(err)
	})

	p.Run("Failed because of invalid Odds", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     1,
			Odds:                0,
			SellingPercentage:   1,
			SellingPrice:        1,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().Error(err)
	})

	p.Run("Failed because of invalid TotalWagerValue", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     0,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().Error(err)
	})

	p.Run("Failed because of invalid SellingPrice monetary format", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1.222,
			CurrentSellingPrice: 1,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().Error(err)
	})

	p.Run("Failed because of consistency CurrentSellingPrice SellingPrice and AmountSold", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1.22,
			CurrentSellingPrice: 1.21,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().Error(err)
	})

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID: 1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1.22,
			CurrentSellingPrice: 1.22,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().NoError(err)

		actualWager, err := r.GetByID(ctx, tx, newWager.ID)
		p.Assert().NoError(err)
		p.Assert().Equal(uint(1), actualWager.ID)
		p.Assert().Equal(1.22, actualWager.SellingPrice)
	})
}

func (p *MySqlRepositoryTestSuite) TestMySqlWagerRepository_List() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	r := repository.NewWagerRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1.22,
			CurrentSellingPrice: 1.22,
			PlacedAt:            time.Now(),
		}
		_, err := r.Create(ctx, tx, newWager)
		p.Assert().NoError(err)
		_, err = r.Create(ctx, tx, newWager)
		p.Assert().NoError(err)
		_, err = r.Create(ctx, tx, newWager)
		p.Assert().NoError(err)

		actualWagers, err := r.List(ctx, tx, 1, 10)
		p.Assert().NoError(err)
		p.Assert().Equal(2, len(actualWagers))
	})
}

func (p *MySqlRepositoryTestSuite) TestMySqlWagerRepository_GetByID() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	r := repository.NewWagerRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID:1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1.22,
			CurrentSellingPrice: 1.22,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().NoError(err)

		actualWager, err := r.GetByID(ctx, tx, newWager.ID)
		p.Assert().NoError(err)
		p.Assert().Equal(uint(1), actualWager.ID)
		p.Assert().Equal(1.22, actualWager.SellingPrice)
	})
}

func (p *MySqlRepositoryTestSuite) TestMySqlWagerRepository_GetByIDAndLockForUpdate() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	r := repository.NewWagerRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID:1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1.22,
			CurrentSellingPrice: 1.22,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().NoError(err)

		actualWager, err := r.GetByIDAndLockForUpdate(ctx, tx, newWager.ID)
		p.Assert().NoError(err)
		p.Assert().Equal(uint(1), actualWager.ID)
		p.Assert().Equal(1.22, actualWager.SellingPrice)
	})
}

func (p *MySqlRepositoryTestSuite) TestMySqlWagerRepository_Update() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	r := repository.NewWagerRepository()

	p.Run("Failed because of consistency CurrentSellingPrice SellingPrice and AmountSold", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID:1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1.22,
			CurrentSellingPrice: 1.22,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().NoError(err)

		newWager.PercentageSold.Uint = 2
		newWager.AmountSold = null.NewFloat64(2.22)
		newWager.PercentageSold.Valid = true
		_, err = r.Update(ctx, tx, newWager)
		p.Assert().Error(err)
	})

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		newWager := model.Wager{
			BaseModel:           model.BaseModel{ID:1},
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1.22,
			CurrentSellingPrice: 1.22,
			PlacedAt:            time.Now(),
		}
		newWager, err := r.Create(ctx, tx, newWager)
		p.Assert().NoError(err)

		newWager.PercentageSold.Uint = 2
		newWager.PercentageSold.Valid = true
		_, err = r.Update(ctx, tx, newWager)
		p.Assert().NoError(err)

		actualWager, err := r.GetByID(ctx, tx, newWager.ID)
		p.Assert().NoError(err)
		p.Assert().Equal(uint(2), actualWager.PercentageSold.Uint)
	})
}
