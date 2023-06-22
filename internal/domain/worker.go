package domain

import (
	"context"

	"cloud.google.com/go/civil"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/cursor"
)

//go:generate mockgen -destination=internal/domain/worker_mock.go -package=domain -source=./internal/domain/worker.go WorkerRepository
type WorkerRepository interface {
	GetAvailableShifts(ctx context.Context, queryCursor *cursor.Cursor, limit int, workerID int64, start civil.Date, end civil.Date) (Shifts, error)
}
