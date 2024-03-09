package worker

import (
	"context"
	"errors"
	"sort"
	"strconv"

	"github.com/rodrigosdo/facilities-api/internal/cursor"
	"github.com/rodrigosdo/facilities-api/internal/domain"

	"cloud.google.com/go/civil"
)

const (
	DefaultLimit = 20
	MaxLimit     = 200
)

//go:generate mockgen -destination=internal/usecase/worker/available_shifts_mock.go -package=worker -source=internal/usecase/worker/available_shifts.go AvailableShifts
type AvailableShifts interface {
	GetAvailableShifts(ctx context.Context, req GetAvailableShiftsRequest) (*GetAvailableShiftsResponse, error)
}

type availableShifts struct {
	workerRepository domain.WorkerRepository
}

func NewAvailableShifts(wr domain.WorkerRepository) AvailableShifts {
	return &availableShifts{
		workerRepository: wr,
	}
}

type GetAvailableShiftsRequest struct {
	Cursor   *cursor.Cursor
	End      civil.Date
	Limit    int
	Start    civil.Date
	WorkerID int64
}

type GetAvailableShiftsResponse struct {
	NextCursor *cursor.Cursor
	Shifts     domain.Shifts
}

func (as *availableShifts) GetAvailableShifts(ctx context.Context, req GetAvailableShiftsRequest) (*GetAvailableShiftsResponse, error) {
	if req.Limit <= 0 || req.Limit > MaxLimit {
		req.Limit = DefaultLimit
	}

	if req.WorkerID == 0 {
		return nil, errors.New("a worker_id is required to get available shifts from a worker")
	}

	if !req.Start.IsZero() && req.End.IsZero() {
		return nil, errors.New("end is required when start is provided")
	}

	if req.Start.IsZero() && !req.End.IsZero() {
		return nil, errors.New("start is required when end is provided")
	}

	shifts, err := as.workerRepository.GetAvailableShifts(
		ctx,
		req.Cursor,
		req.Limit,
		req.WorkerID,
		req.Start,
		req.End,
	)
	if err != nil {
		return nil, err
	}

	if len(shifts) == 0 {
		return &GetAvailableShiftsResponse{}, nil
	}

	sort.Slice(shifts, func(i, j int) bool {
		return shifts[i].ID < shifts[j].ID
	})

	nextCursor := cursor.New(cursor.DirectionAfter, strconv.FormatInt(shifts[len(shifts)-1].ID, 10))

	return &GetAvailableShiftsResponse{
		NextCursor: nextCursor,
		Shifts:     shifts,
	}, nil
}
