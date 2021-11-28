package purchase

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/duyquang6/wager-management-be/pkg/dto"
	"github.com/duyquang6/wager-management-be/pkg/exception"
	_validator "github.com/duyquang6/wager-management-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

func (s *Controller) HandleCreatePurchase() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "read data failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		req := dto.CreatePurchaseRequest{}
		err = json.Unmarshal(data, &req)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "unmarshal failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		wagerIDStr := c.Param("wager_id")
		wagerID, err := strconv.ParseUint(wagerIDStr, 10, 32)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse user_id failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		req.WagerID = uint(wagerID)

		if err := _validator.GetValidate().Struct(req); err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "validation error").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		res, err := s.service.CreatePurchase(ctx, req)
		if err != nil {
			if appErr, ok := err.(exception.AppError); ok {
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusCreated, res)
	}
}
