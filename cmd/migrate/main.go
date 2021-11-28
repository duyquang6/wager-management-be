// This package serve for db migration
package main

import (
	"context"
	"log"

	"github.com/duyquang6/wager-management-be/internal/database"
	"github.com/duyquang6/wager-management-be/internal/setup"
	"github.com/duyquang6/wager-management-be/pkg/logging"
	"github.com/sethvargo/go-signalcontext"
)

func main() {
	ctx, done := signalcontext.OnInterrupt()

	logger := logging.NewLoggerFromEnv()
	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Fatalw("application panic", "panic", r)
		}
	}()

	err := realMain(ctx)
	done()

	if err != nil {
		log.Fatal(err)
	}
	logger.Info("successful shutdown")
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	var config database.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		logger.Fatal(err)
	}
	if err := env.Database().Migrate(ctx); err != nil {
		logger.Fatal("cannot migrate: %v", err.Error())
	}
	return nil
}
