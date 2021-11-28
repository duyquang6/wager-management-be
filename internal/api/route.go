package api

import (
	"context"
	purchaseControllerPkg "github.com/duyquang6/wager-management-be/internal/controller/purchase"
	wagerControllerPkg "github.com/duyquang6/wager-management-be/internal/controller/wager"
	"github.com/duyquang6/wager-management-be/internal/database"
	"github.com/duyquang6/wager-management-be/internal/middleware"
	"github.com/duyquang6/wager-management-be/internal/repository"
	"github.com/duyquang6/wager-management-be/internal/service"
	"github.com/duyquang6/wager-management-be/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *httpServer) setupDependencyAndRouter(ctx context.Context, r *gin.Engine, db *database.DB) {
	wagerRepo := repository.NewWagerRepository()
	purchaseRepo := repository.NewPurchaseRepository()

	wagerService := service.NewWagerService(db, wagerRepo)
	purchaseService := service.NewPurchaseService(db, purchaseRepo, wagerRepo)

	wagerController := wagerControllerPkg.NewController(wagerService)
	purchaseController := purchaseControllerPkg.NewController(purchaseService)

	initRoute(ctx, r, wagerController, purchaseController)
}


func initRoute(ctx context.Context, r *gin.Engine,
	wagerController *wagerControllerPkg.Controller, purchaseController *purchaseControllerPkg.Controller) {
	r.Use(middleware.PopulateRequestID())
	r.Use(middleware.PopulateLogger(logging.FromContext(ctx)))
	{
		sub := r.Group("/wagers")
		{
			sub.POST("", wagerController.HandleCreateWager())
			sub.GET("", wagerController.HandleListWager())
		}

	}

	r.POST("/buy/:wager_id", purchaseController.HandleCreatePurchase())

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
