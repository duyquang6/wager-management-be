package purchase

import (
	"bytes"
	"encoding/json"
	"github.com/duyquang6/wager-management-be/internal/service/mocks"
	"github.com/duyquang6/wager-management-be/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreatePurchase(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockPurchaseService := new(mocks.PurchaseService)
	c := NewController(mockPurchaseService)
	r.POST("/buy/:wager_id", c.HandleCreatePurchase())

	t.Run("Success", func(t *testing.T) {
		mockResp := dto.CreatePurchaseResponse{
			ID:          0,
			BuyingPrice: 0,
			WagerID:     0,
			BoughtAt:    0,
		}
		mockReq := dto.CreatePurchaseRequest{
			BuyingPrice: 1,
			WagerID: 1,
		}

		mockPurchaseService.On("CreatePurchase", mock.Anything, mockReq).Return(mockResp, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/buy/1", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockPurchaseService.AssertExpectations(t)
	})

	t.Run("Invalid wager ID", func(t *testing.T) {
		mockReq := dto.CreatePurchaseRequest{
			BuyingPrice: 1,
			WagerID: 0,
		}

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()


		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/buy/0", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		jsonActual := rr.Body.String()
		assert.Equal(t, strings.Contains(jsonActual, "WagerID"), true)
	})

	t.Run("Invalid buying price", func(t *testing.T) {
		mockReq := dto.CreatePurchaseRequest{
			BuyingPrice: 0,
			WagerID: 1,
		}

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()


		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/buy/1", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		jsonActual := rr.Body.String()
		assert.Equal(t, strings.Contains(jsonActual, "BuyingPrice"), true)
	})
}
