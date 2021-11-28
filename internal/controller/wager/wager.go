package wager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/duyquang6/wager-management-be/pkg/dto"
	"github.com/duyquang6/wager-management-be/pkg/exception"
	"github.com/duyquang6/wager-management-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

func (s *Controller) HandleListWager() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err         error
			page, limit uint64
		)
		ctx := c.Request.Context()

		pageStr := c.Query("page")
		if len(pageStr) == 0 {
			// Default 1st page
			page = 1
		} else {
			page, err = strconv.ParseUint(pageStr, 10, 32)
			if err != nil {
				appErr := exception.Wrap(http.StatusBadRequest, err, "parse page field failed").(exception.AppError)
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
		}

		limitStr := c.Query("limit")
		if len(limitStr) == 0 {
			// Default limit is 10 items
			limit = 10
		} else {
			limit, err = strconv.ParseUint(limitStr, 10, 32)
			if err != nil {
				appErr := exception.Wrap(http.StatusBadRequest, err, "parse limit field failed").(exception.AppError)
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
		}

		req := dto.ListWagersRequest{
			Page:  uint(page),
			Limit: uint(limit),
		}

		if err := validator.GetValidate().Struct(req); err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "validation error").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		resp, err := s.service.ListWagers(ctx, req.Page, req.Limit)
		if err != nil {
			if appErr, ok := err.(exception.AppError); ok {
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, resp.Data)
	}
}

func (s *Controller) HandleCreateWager() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "read data failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		req := dto.CreateWagerRequest{}
		err = json.Unmarshal(data, &req)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "unmarshal failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		if err := validator.GetValidate().Struct(req); err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "validation error").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		if req.SellingPrice <= float64(req.TotalWagerValue*req.SellingPercentage)/100 {
			appErr := exception.New(http.StatusBadRequest,
				"field SellingPrice must be larger than TotalWagerValue * SellingPercentage").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		resp, err := s.service.CreateWager(ctx, req)
		if err != nil {
			if appErr, ok := err.(exception.AppError); ok {
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}
