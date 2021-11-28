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

func TestNewWagerService(t *testing.T) {
	t.Parallel()

	type args struct {
		dbFactory database.DBFactory
		wagerRepo repository.WagerRepository
	}
	wagerRepoMock := &_mockRepo.WagerRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}
	tests := []struct {
		name string
		args args
		want WagerService
	}{
		{
			name: "TC1_NewUserServiceSuccess",
			args: args{
				dbFactory: dbFactoryMock,
				wagerRepo: wagerRepoMock,
			},
			want: &wagerSvc{dbFactoryMock, wagerRepoMock},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWagerService(tt.args.dbFactory, tt.args.wagerRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wagerSvc_ListWagers(t *testing.T) {
	t.Parallel()

	type fields struct {
		db        database.DBFactory
		wagerRepo repository.WagerRepository
	}

	type args struct {
		ctx   context.Context
		page  uint
		limit uint
	}

	timeNow := time.Now().UTC()
	wagers := []model.Wager{
		{BaseModel: model.BaseModel{ID: 1}, PlacedAt: timeNow},
		{BaseModel: model.BaseModel{ID: 2}, PlacedAt: timeNow},
	}
	dbFactoryMock := &_mockDB.DBFactory{}
	dbFactoryMock.On("GetDB").Return(&gorm.DB{})

	wagerRepoMockNotFound := &_mockRepo.WagerRepository{}
	wagerRepoMockNotFound.On("List", mock.Anything, mock.Anything,
		uint(0), uint(1)).Return(nil, gorm.ErrRecordNotFound)

	wagerRepoMockFailed := &_mockRepo.WagerRepository{}
	wagerRepoMockFailed.On("List", mock.Anything, mock.Anything,
		uint(0), uint(1)).Return(nil, errors.New("unexpected error"))

	wagerRepoMockSuccess := &_mockRepo.WagerRepository{}
	wagerRepoMockSuccess.On("List", mock.Anything, mock.Anything, uint(0), uint(2)).Return(wagers, nil)

	expectWagers := []dto.Wager{{ID: 1, PlacedAt: uint(timeNow.Unix())}, {ID: 2, PlacedAt: uint(timeNow.Unix())}}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.ListWagersResponse
		wantErr bool
	}{
		{
			name: "TC1_NotFound",
			fields: fields{
				db:        dbFactoryMock,
				wagerRepo: wagerRepoMockNotFound,
			},
			args: args{
				ctx:  context.TODO(),
				page: 1, limit: 1,
			},
			wantErr: true,
		},
		{
			name: "TC2_DBError",
			fields: fields{
				db:        dbFactoryMock,
				wagerRepo: wagerRepoMockFailed,
			},
			args: args{
				ctx:  context.TODO(),
				page: 1, limit: 1,
			},
			wantErr: true,
		},
		{
			name: "TC3_Success",
			fields: fields{
				db:        dbFactoryMock,
				wagerRepo: wagerRepoMockSuccess,
			},
			args: args{
				ctx:  context.TODO(),
				page: 1, limit: 2,
			},
			wantErr: false,
			want: dto.ListWagersResponse{
				Data: expectWagers,
			},
		},
		{
			name: "TC4_InvalidLimitParam",
			fields: fields{
				db:        dbFactoryMock,
				wagerRepo: wagerRepoMockSuccess,
			},
			args: args{
				ctx:  context.TODO(),
				page: 1, limit: 0,
			},
			wantErr: true,
		},
		{
			name: "TC5_InvalidPageParam",
			fields: fields{
				db:        dbFactoryMock,
				wagerRepo: wagerRepoMockSuccess,
			},
			args: args{
				ctx:  context.TODO(),
				page: 0, limit: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &wagerSvc{
				dbFactory: tt.fields.db,
				wagerRepo: tt.fields.wagerRepo,
			}
			got, err := s.ListWagers(tt.args.ctx, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("wagerSvc.ListWagers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wagerSvc.ListWagers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wagerSvc_CreateWager(t *testing.T) {
	t.Parallel()

	type fields struct {
		db        database.DBFactory
		wagerRepo repository.WagerRepository
	}

	type args struct {
		ctx context.Context
		req dto.CreateWagerRequest
	}

	req := dto.CreateWagerRequest{
		TotalWagerValue:   1,
		Odds:              1,
		SellingPercentage: 1,
		SellingPrice:      1,
	}

	timeNow := time.Now().UTC()
	wager := model.Wager{
		BaseModel:           model.BaseModel{ID: 1},
		TotalWagerValue:     0,
		Odds:                0,
		SellingPercentage:   0,
		SellingPrice:        0,
		CurrentSellingPrice: 0,
		PercentageSold:      null.Uint{},
		AmountSold:          null.Float64{},
		PlacedAt:            timeNow,
	}

	dbFactoryMock := &_mockDB.DBFactory{}
	dbFactoryMock.On("GetDB").Return(&gorm.DB{})

	wagerRepoMockFailed := &_mockRepo.WagerRepository{}
	wagerRepoMockFailed.On("Create", mock.Anything, mock.Anything,
		mock.Anything).Return(model.Wager{}, errors.New("unexpected error"))

	wagerRepoMockSuccess := &_mockRepo.WagerRepository{}
	wagerRepoMockSuccess.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(wager, nil)

	expectWager := dto.Wager{ID: 1, PlacedAt: uint(timeNow.Unix())}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.CreateWagerResponse
		wantErr bool
	}{
		{
			name: "TC1_DBError",
			fields: fields{
				db:        dbFactoryMock,
				wagerRepo: wagerRepoMockFailed,
			},
			args: args{
				ctx: context.TODO(),
				req: req,
			},
			wantErr: true,
		},
		{
			name: "TC2_Success",
			fields: fields{
				db:        dbFactoryMock,
				wagerRepo: wagerRepoMockSuccess,
			},
			args: args{
				ctx: context.TODO(),
				req: req,
			},
			wantErr: false,
			want:    dto.CreateWagerResponse{Wager: expectWager},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &wagerSvc{
				dbFactory: tt.fields.db,
				wagerRepo: tt.fields.wagerRepo,
			}
			got, err := s.CreateWager(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("wagerSvc.ListWagers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wagerSvc.ListWagers() = %v, want %v", got, tt.want)
			}
		})
	}
}
