package middleware

import (
	"github.com/duyquang6/wager-management-be/internal/controller"
	"github.com/duyquang6/wager-management-be/pkg/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PopulateLogger add logger context
func PopulateLogger(originalLogger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := originalLogger

		// Only override the logger if it's the default logger. This is only used
		// for testing and is intentionally a strict object equality check because
		// the default logger is a global default in the logger package.
		if existing := logging.FromContext(ctx); existing == logging.DefaultLogger() {
			logger = existing
		}

		// If there's a request ID, set that on the logger.
		if id := controller.RequestIDFromContext(ctx); id != "" {
			logger = logger.With("request_id", id)
		}

		ctx = logging.WithLogger(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		c.Next()
	}
}
