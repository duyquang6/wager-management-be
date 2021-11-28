package wager

import (
	"bytes"
	"encoding/json"
	"github.com/duyquang6/wager-management-be/internal/service/mocks"
	"github.com/duyquang6/wager-management-be/pkg/dto"
	"github.com/duyquang6/wager-management-be/pkg/null"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setup() (*mocks.WagerService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockWagerService := new(mocks.WagerService)
	c := NewController(mockWagerService)
	r.POST("/wagers", c.HandleCreateWager())
	r.GET("/wagers", c.HandleListWager())
	return mockWagerService, r
}

func TestCreateWager(t *testing.T) {
	mockWagerService, r := setup()

	t.Run("Success", func(t *testing.T) {
		mockReq := dto.CreateWagerRequest{
			TotalWagerValue:   1,
			Odds:              1,
			SellingPercentage: 1,
			SellingPrice:      1.32,
		}
		mockWager := dto.Wager{
			ID:                  1,
			TotalWagerValue:     1,
			Odds:                1,
			SellingPercentage:   1,
			SellingPrice:        1,
			CurrentSellingPrice: 1,
			PercentageSold:      null.Uint{},
			AmountSold:          null.Float64{},
			PlacedAt:            uint(time.Now().UTC().Unix()),
		}
		mockResp := dto.CreateWagerResponse{Wager: mockWager}
		mockWagerService.On("CreateWager", mock.Anything, mockReq).Return(mockResp, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/wagers", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockWagerService.AssertExpectations(t)
	})

	t.Run("Failed because of SellingPrice larger than TotalWagerValue x SellingPercentage", func(t *testing.T) {
		mockReq := dto.CreateWagerRequest{
			TotalWagerValue:   2,
			Odds:              1,
			SellingPercentage: 50,
			SellingPrice:      1,
		}

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/wagers", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		jsonActual := rr.Body.String()
		assert.Equal(t, strings.Contains(jsonActual, "SellingPrice"), true)
	})

	t.Run("Failed due to invalid totalwagervalue", func(t *testing.T) {
		mockReq := dto.CreateWagerRequest{
			TotalWagerValue:   0,
			Odds:              1,
			SellingPercentage: 1,
			SellingPrice:      1,
		}

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/wagers", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		jsonActual := rr.Body.String()
		assert.Equal(t, strings.Contains(jsonActual, "TotalWagerValue"), true)
	})

	t.Run("Failed due to invalid selling price more than 2 decimal", func(t *testing.T) {
		mockReq := dto.CreateWagerRequest{
			TotalWagerValue:   1,
			Odds:              1,
			SellingPercentage: 1,
			SellingPrice:      1.333,
		}

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/wagers", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		jsonActual := rr.Body.String()
		assert.Equal(t, strings.Contains(jsonActual, "SellingPrice"), true)
	})
	t.Run("Failed due to invalid odd and selling price", func(t *testing.T) {
		mockReq := dto.CreateWagerRequest{
			TotalWagerValue:   1,
			Odds:              0,
			SellingPercentage: 1,
			SellingPrice:      1.333,
		}

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/wagers", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		jsonActual := rr.Body.String()
		assert.Equal(t, strings.Contains(jsonActual, "Odds"), true)
		assert.Equal(t, strings.Contains(jsonActual, "SellingPrice"), true)
	})
	t.Run("Failed due to invalid selling percentage and selling price", func(t *testing.T) {
		mockReq := dto.CreateWagerRequest{
			TotalWagerValue:   1,
			Odds:              1,
			SellingPercentage: 101,
			SellingPrice:      1.333,
		}

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/wagers", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		jsonActual := rr.Body.String()
		assert.Equal(t, strings.Contains(jsonActual, "SellingPercentage"), true)
		assert.Equal(t, strings.Contains(jsonActual, "SellingPrice"), true)
	})
}

func TestListWager(t *testing.T) {
	mockWagerService, r := setup()

	t.Run("Success", func(t *testing.T) {
		mockPage := uint(1)
		mockLimit := uint(1)
		wagers := []dto.Wager{
			{ID: 1},
			{ID: 2},
		}
		mockResp := dto.ListWagersResponse{Data: wagers}
		mockWagerService.On("ListWagers", mock.Anything, mockPage, mockLimit).Return(mockResp, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/wagers?page=1&limit=1", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp.Data)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockWagerService.AssertExpectations(t)
	})

	t.Run("Success with default page & limit", func(t *testing.T) {
		defaultPage := uint(1)
		defaultLimit := uint(10)
		wagers := []dto.Wager{
			{ID: 1},
			{ID: 2},
		}
		mockResp := dto.ListWagersResponse{Data: wagers}
		mockWagerService.On("ListWagers", mock.Anything, defaultPage, defaultLimit).Return(mockResp, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/wagers", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp.Data)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockWagerService.AssertExpectations(t)
	})

	t.Run("Failed due to invalid page & limit", func(t *testing.T) {

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/wagers?page=0&limit=0", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		jsonActual := rr.Body.String()
		assert.Equal(t, strings.Contains(jsonActual, "Page"), true)
		assert.Equal(t, strings.Contains(jsonActual, "Limit"), true)
	})
}