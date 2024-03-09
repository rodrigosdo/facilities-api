package domain

import (
	"context"

	"github.com/rodrigosdo/facilities-api/internal/cursor"

	"cloud.google.com/go/civil"
)

//go:generate mockgen -destination=internal/domain/worker_mock.go -package=domain -source=./internal/domain/worker.go WorkerRepository
type WorkerRepository interface {
	GetAvailableShifts(ctx context.Context, queryCursor *cursor.Cursor, limit int, workerID int64, start civil.Date, end civil.Date) (Shifts, error)
}
