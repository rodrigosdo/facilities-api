package main

import (
	"context"
	"fmt"

	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/config"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/logger"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/postgres"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/server"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/usecase/worker"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	logger, err := logger.New(cfg.Log.Level)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	database, err := postgres.NewDatabase(ctx, cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to postgres",
			zap.Any("error", err),
		)
	}

	workerAvailableShiftsUseCase := worker.NewAvailableShifts(database)

	srv, err := server.New(
		fmt.Sprintf(":%s", cfg.Server.Port),
		cfg,
		logger,
		database,
		workerAvailableShiftsUseCase,
	)
	if err != nil {
		logger.Fatal("failed to build server",
			zap.Any("error", err),
		)
	}

	if err := srv.Run(ctx); err != nil {
		logger.Info("server exited",
			zap.Any("error", err),
		)
	}
}
