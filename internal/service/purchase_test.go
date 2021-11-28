package service

import (
	"context"
	"errors"
	"github.com/duyquang6/wager-management-be/internal/database"
	_mockDB "github.com/duyquang6/wager-management-be/internal/database/mocks"
	"github.com/duyquang6/wager-management-be/internal/model"
	"github.com/duyquang6/wager-management-be/internal/repository"
	_mockRepo "github.com/duyquang6/wager-management-be/internal/repository/mocks"
	"github.com/duyquang6/wager-management-be/pkg/dto"
	"github.com/duyquang6/wager-management-be/pkg/null"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

func TestNewPurchaseService(t *testing.T) {
	t.Parallel()
	type args struct {
		dbFactory    database.DBFactory
		wagerRepo    repository.WagerRepository
		purchaseRepo repository.PurchaseRepository
	}
	wagerRepoMock := &_mockRepo.WagerRepository{}
	purchaseRepoMock := &_mockRepo.PurchaseRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}
	tests := []struct {
		name string
		args args
		want PurchaseService
	}{
		{
			name: "TC1_NewPurchaseServiceSuccess",
			args: args{
				dbFactory:    dbFactoryMock,
				wagerRepo:    wagerRepoMock,
				purchaseRepo: purchaseRepoMock,
			},
			want: &purchaseSvc{dbFactoryMock, purchaseRepoMock, wagerRepoMock},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPurchaseService(tt.args.dbFactory, tt.args.purchaseRepo, tt.args.wagerRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPurchaseService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wagerSvc_CreatePurchase(t *testing.T) {
	t.Parallel()

	type fields struct {
		db           database.DBFactory
		wagerRepo    repository.WagerRepository
		purchaseRepo repository.PurchaseRepository
	}

	type args struct {
		ctx context.Context
		req dto.CreatePurchaseRequest
	}

	timeNow := time.Now().UTC()
	mockWager := model.Wager{
		BaseModel: model.BaseModel{ID: 1}, PlacedAt: timeNow,
		TotalWagerValue:     1000,
		SellingPercentage:   50,
		SellingPrice:        600,
		CurrentSellingPrice: 500,
		AmountSold:          null.NewFloat64(100),
		PercentageSold:      null.NewUint(13),
	}

	updatedMockWager := model.Wager{
		BaseModel: model.BaseModel{ID: 1}, PlacedAt: timeNow,
		TotalWagerValue:     1000,
		SellingPercentage:   50,
		SellingPrice:        600,
		CurrentSellingPrice: 400,
		AmountSold:          null.NewFloat64(200),
		PercentageSold:      null.NewUint(33),
	}
	mockPurchase := model.Purchase{
		BaseModel: model.BaseModel{ID: 1},
		WagerID:     updatedMockWager.ID,
		Wager:       updatedMockWager,
		BuyingPrice: 100,
		BoughtAt:    time.Now().UTC(),
	}
	mockReq := dto.CreatePurchaseRequest{
		BuyingPrice: 100,
		WagerID:     1,
	}

	//db := &gorm.DB{}
	dbFactoryMock := &_mockDB.DBFactory{}
	dbFactoryMock.On("GetDBWithTx").Return(&gorm.DB{})
	dbFactoryMock.On("Rollback", mock.Anything).Return()
	dbFactoryMock.On("Commit", mock.Anything).Return()


	//wagerRepoMockHappyCase := &_mockRepo.WagerRepository{}
	//wagerRepoMockHappyCase.On("GetByID", mock.Anything, mock.Anything,
	//	uint(1)).Return(nil, gorm.ErrRecordNotFound)

	wagerRepoMockUpdateFailed := &_mockRepo.WagerRepository{}
	wagerRepoMockUpdateFailed.On("GetByIDAndLockForUpdate",
		mock.Anything, mock.Anything, uint(1)).Return(mockWager, nil)
	wagerRepoMockUpdateFailed.On("Update",
		mock.Anything, mock.Anything, updatedMockWager).Return(model.Wager{}, errors.New("unexpected error"))

	wagerRepoMockGetByIDAndLockFailed := &_mockRepo.WagerRepository{}
	wagerRepoMockGetByIDAndLockFailed.On("GetByIDAndLockForUpdate", mock.Anything,
		mock.Anything, uint(1)).Return(model.Wager{}, errors.New("unexpected error"))

	wagerRepoMockIDNotFound:= &_mockRepo.WagerRepository{}
	wagerRepoMockIDNotFound.On("GetByIDAndLockForUpdate",
		mock.Anything, mock.Anything, uint(1)).Return(model.Wager{}, gorm.ErrRecordNotFound)

	wagerRepoMockSuccess := &_mockRepo.WagerRepository{}
	wagerRepoMockSuccess.On("GetByIDAndLockForUpdate",
		mock.Anything, mock.Anything, uint(1)).Return(mockWager, nil)
	wagerRepoMockSuccess.On("Update", mock.Anything,
		mock.Anything, updatedMockWager).Return(updatedMockWager, nil)

	purchaseRepoMockSuccess := &_mockRepo.PurchaseRepository{}
	purchaseRepoMockSuccess.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(mockPurchase, nil)

	purchaseRepoMockCreateFailed := &_mockRepo.PurchaseRepository{}
	purchaseRepoMockCreateFailed.On("Create", mock.Anything,
		mock.Anything, mock.Anything).Return(model.Purchase{}, errors.New("unexpected error"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.CreatePurchaseResponse
		wantErr bool
	}{
		{
			name: "TC1_Success",
			fields: fields{
				db:           dbFactoryMock,
				wagerRepo:    wagerRepoMockSuccess,
				purchaseRepo: purchaseRepoMockSuccess,
			},
			args: args{
				ctx: context.TODO(),
				req: mockReq,
			},
			wantErr: false,
			want:    dto.CreatePurchaseResponse{
				ID:          mockPurchase.ID,
				BuyingPrice: mockPurchase.BuyingPrice,
				WagerID:     mockPurchase.WagerID,
				BoughtAt:    uint(mockPurchase.BoughtAt.Unix()),
			},
		},
		{
			name: "TC2_RelatedWagerIDNotFound",
			fields: fields{
				db:           dbFactoryMock,
				wagerRepo:    wagerRepoMockIDNotFound,
			},
			args: args{
				ctx: context.TODO(),
				req: mockReq,
			},
			wantErr: true,
		},
		{
			name: "TC3_GetByIDAndLockForUpdateFailed",
			fields: fields{
				db:        dbFactoryMock,
				wagerRepo: wagerRepoMockGetByIDAndLockFailed,
			},
			args: args{
				ctx: context.TODO(),
				req: mockReq,
			},
			wantErr: true,
		},
		{
			name: "TC4_UpdateWagerFailed",
			fields: fields{
				db:           dbFactoryMock,
				wagerRepo:    wagerRepoMockUpdateFailed,
			},
			args: args{
				ctx: context.TODO(),
				req: mockReq,
			},
			wantErr: true,
		},
		{
			name: "TC5_CreatePurchaseFailed",
			fields: fields{
				db:           dbFactoryMock,
				wagerRepo:    wagerRepoMockSuccess,
				purchaseRepo: purchaseRepoMockCreateFailed,
			},
			args: args{
				ctx: context.TODO(),
				req: mockReq,
			},
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &purchaseSvc{
				dbFactory:    tt.fields.db,
				wagerRepo:    tt.fields.wagerRepo,
				purchaseRepo: tt.fields.purchaseRepo,
			}
			got, err := s.CreatePurchase(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("purchaseSvc.CreatePurchase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("purchaseSvc.CreatePurchase() = %v, want %v", got, tt.want)
			}
		})
	}
}
