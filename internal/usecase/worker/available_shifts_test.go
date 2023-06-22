package worker_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/golang/mock/gomock"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/cursor"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/domain"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/usecase/worker"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailableShifts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fakeCursor := cursor.Cursor{}
	fakeLimit := 10
	fakeWorkerID := int64(123123)
	fakeStartDate := civil.Date{Year: 2023, Month: 06, Day: 04}
	fakeEndDate := civil.Date{Year: 2023, Month: 06, Day: 10}
	fakeShifts := domain.Shifts{
		{
			End:   time.Time{},
			ID:    123,
			Start: time.Time{},
		},
		{
			End:   time.Time{},
			ID:    321,
			Start: time.Time{},
		},
	}

	t.Run("should return error when no worker is provided", func(t *testing.T) {
		t.Parallel()

		us := worker.NewAvailableShifts(nil)
		availableShifts, err := us.GetAvailableShifts(
			ctx,
			worker.GetAvailableShiftsRequest{
				End:      civil.Date{},
				Start:    civil.Date{},
				WorkerID: 0,
			},
		)
		assert.Error(t, err)
		assert.Nil(t, availableShifts)
	})

	t.Run("should return error when start date is provided without end date", func(t *testing.T) {
		t.Parallel()

		us := worker.NewAvailableShifts(nil)
		availableShifts, err := us.GetAvailableShifts(
			ctx,
			worker.GetAvailableShiftsRequest{
				End:      civil.Date{},
				Start:    fakeStartDate,
				WorkerID: fakeWorkerID,
			},
		)
		assert.Error(t, err)
		assert.Nil(t, availableShifts)
	})

	t.Run("should return error when end date is provided without start date", func(t *testing.T) {
		t.Parallel()

		us := worker.NewAvailableShifts(nil)
		availableShifts, err := us.GetAvailableShifts(
			ctx,
			worker.GetAvailableShiftsRequest{
				End:      fakeEndDate,
				Start:    civil.Date{},
				WorkerID: fakeWorkerID,
			},
		)
		assert.Error(t, err)
		assert.Nil(t, availableShifts)
	})

	t.Run("should return error if fails to get available shifts", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWorkerRepository := domain.NewMockWorkerRepository(ctrl)
		mockWorkerRepository.
			EXPECT().
			GetAvailableShifts(ctx, &fakeCursor, fakeLimit, fakeWorkerID, fakeStartDate, fakeEndDate).
			Times(1).
			Return(nil, errors.New("fake error"))

		us := worker.NewAvailableShifts(mockWorkerRepository)
		availableShifts, err := us.GetAvailableShifts(
			ctx,
			worker.GetAvailableShiftsRequest{
				Cursor:   &fakeCursor,
				End:      fakeEndDate,
				Limit:    fakeLimit,
				Start:    fakeStartDate,
				WorkerID: fakeWorkerID,
			},
		)
		assert.Error(t, err)
		assert.Nil(t, availableShifts)
	})

	t.Run("should successfully get available shifts", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWorkerRepository := domain.NewMockWorkerRepository(ctrl)
		mockWorkerRepository.
			EXPECT().
			GetAvailableShifts(ctx, &fakeCursor, fakeLimit, fakeWorkerID, fakeStartDate, fakeEndDate).
			Times(1).
			Return(fakeShifts, nil)

		us := worker.NewAvailableShifts(mockWorkerRepository)

		availableShifts, err := us.GetAvailableShifts(
			ctx,
			worker.GetAvailableShiftsRequest{
				Cursor:   &fakeCursor,
				End:      fakeEndDate,
				Limit:    fakeLimit,
				Start:    fakeStartDate,
				WorkerID: fakeWorkerID,
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, &worker.GetAvailableShiftsResponse{
			NextCursor: cursor.New(cursor.DirectionAfter, "321"),
			Shifts:     fakeShifts,
		}, availableShifts)
	})

	t.Run("should successfully get available shifts with default limit", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWorkerRepository := domain.NewMockWorkerRepository(ctrl)
		mockWorkerRepository.
			EXPECT().
			GetAvailableShifts(ctx, &fakeCursor, worker.DefaultLimit, fakeWorkerID, fakeStartDate, fakeEndDate).
			Times(1).
			Return(fakeShifts, nil)

		us := worker.NewAvailableShifts(mockWorkerRepository)

		availableShifts, err := us.GetAvailableShifts(
			ctx,
			worker.GetAvailableShiftsRequest{
				Cursor:   &fakeCursor,
				End:      fakeEndDate,
				Start:    fakeStartDate,
				WorkerID: fakeWorkerID,
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, &worker.GetAvailableShiftsResponse{
			NextCursor: cursor.New(cursor.DirectionAfter, "321"),
			Shifts:     fakeShifts,
		}, availableShifts)
	})

	t.Run("should successfully get available shifts without a cursor", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWorkerRepository := domain.NewMockWorkerRepository(ctrl)
		mockWorkerRepository.
			EXPECT().
			GetAvailableShifts(ctx, nil, worker.DefaultLimit, fakeWorkerID, fakeStartDate, fakeEndDate).
			Times(1).
			Return(fakeShifts, nil)

		us := worker.NewAvailableShifts(mockWorkerRepository)

		availableShifts, err := us.GetAvailableShifts(
			ctx,
			worker.GetAvailableShiftsRequest{
				End:      fakeEndDate,
				Start:    fakeStartDate,
				WorkerID: fakeWorkerID,
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, &worker.GetAvailableShiftsResponse{
			NextCursor: cursor.New(cursor.DirectionAfter, "321"),
			Shifts:     fakeShifts,
		}, availableShifts)
	})

	t.Run("should return no error when there's no shifts available", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWorkerRepository := domain.NewMockWorkerRepository(ctrl)
		mockWorkerRepository.
			EXPECT().
			GetAvailableShifts(ctx, nil, worker.DefaultLimit, fakeWorkerID, fakeStartDate, fakeEndDate).
			Times(1).
			Return(domain.Shifts{}, nil)

		us := worker.NewAvailableShifts(mockWorkerRepository)

		availableShifts, err := us.GetAvailableShifts(
			ctx,
			worker.GetAvailableShiftsRequest{
				End:      fakeEndDate,
				Start:    fakeStartDate,
				WorkerID: fakeWorkerID,
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, &worker.GetAvailableShiftsResponse{}, availableShifts)
	})
}
