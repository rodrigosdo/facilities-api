package server_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rodrigosdo/facilities-api/internal/domain"
	"github.com/rodrigosdo/facilities-api/internal/postgres"
	"github.com/rodrigosdo/facilities-api/internal/server"
	"github.com/rodrigosdo/facilities-api/internal/usecase/worker"

	"cloud.google.com/go/civil"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailableShiftsFromUser(t *testing.T) {
	t.Parallel()

	fakeStartDate := civil.Date{Year: 2023, Month: 06, Day: 04}
	fakeEndDate := civil.Date{Year: 2023, Month: 06, Day: 10}
	fakeWorkerID := int64(123123)
	fakeShifts := domain.Shifts{
		{
			End:   time.Time{},
			ID:    123,
			Start: time.Time{},
		},
	}

	t.Run("should successfully return available shifts", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAvailableShiftsUseCase := worker.NewMockAvailableShifts(ctrl)
		mockAvailableShiftsUseCase.
			EXPECT().
			GetAvailableShifts(gomock.Any(), worker.GetAvailableShiftsRequest{
				End:      fakeEndDate,
				Start:    fakeStartDate,
				WorkerID: fakeWorkerID,
			}).
			Times(1).
			Return(&worker.GetAvailableShiftsResponse{Shifts: fakeShifts}, nil)

		handler := server.GetAvailableShiftsFromWorker(mockAvailableShiftsUseCase)

		rt := httprouter.New()
		rt.Handler(http.MethodGet, "/v1/workers/:id/available_shifts", server.HandlerFunc(handler))

		req, err := http.NewRequest("GET", "/v1/workers/123123/available_shifts?end=2023-06-10&start=2023-06-04", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		rt.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should successfully return even when there's no available shifts", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAvailableShiftsUseCase := worker.NewMockAvailableShifts(ctrl)
		mockAvailableShiftsUseCase.
			EXPECT().
			GetAvailableShifts(gomock.Any(), worker.GetAvailableShiftsRequest{
				End:      fakeEndDate,
				Start:    fakeStartDate,
				WorkerID: fakeWorkerID,
			}).
			Times(1).
			Return(&worker.GetAvailableShiftsResponse{}, nil)

		handler := server.GetAvailableShiftsFromWorker(mockAvailableShiftsUseCase)

		rt := httprouter.New()
		rt.Handler(http.MethodGet, "/v1/workers/:id/available_shifts", server.HandlerFunc(handler))

		req, err := http.NewRequest("GET", "/v1/workers/123123/available_shifts?end=2023-06-10&start=2023-06-04", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		rt.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return error when given an invalid cursor", func(t *testing.T) {
		t.Parallel()

		handler := server.GetAvailableShiftsFromWorker(nil)

		req, err := http.NewRequest("GET", "?cursor=invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error when given an invalid start date", func(t *testing.T) {
		t.Parallel()

		handler := server.GetAvailableShiftsFromWorker(nil)

		req, err := http.NewRequest("GET", "?start=invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error when given an invalid end date", func(t *testing.T) {
		t.Parallel()

		handler := server.GetAvailableShiftsFromWorker(nil)

		req, err := http.NewRequest("GET", "?end=invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error when given an invalid limit", func(t *testing.T) {
		t.Parallel()

		handler := server.GetAvailableShiftsFromWorker(nil)

		req, err := http.NewRequest("GET", "?limit=invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error when given an worker id", func(t *testing.T) {
		t.Parallel()

		handler := server.GetAvailableShiftsFromWorker(nil)

		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error if fails to get available shifts", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAvailableShiftsUseCase := worker.NewMockAvailableShifts(ctrl)
		mockAvailableShiftsUseCase.
			EXPECT().
			GetAvailableShifts(gomock.Any(), worker.GetAvailableShiftsRequest{
				End:      fakeEndDate,
				Start:    fakeStartDate,
				WorkerID: fakeWorkerID,
			}).
			Times(1).
			Return(nil, errors.New("error"))

		handler := server.GetAvailableShiftsFromWorker(mockAvailableShiftsUseCase)

		rt := httprouter.New()
		rt.Handler(http.MethodGet, "/v1/workers/:id/available_shifts", server.HandlerFunc(handler))

		req, err := http.NewRequest("GET", "/v1/workers/123123/available_shifts?end=2023-06-10&start=2023-06-04", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		rt.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestHealthcheck(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("should return status code 200", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConn := postgres.NewMockConn(ctrl)
		mockConn.
			EXPECT().
			Ping(ctx).
			Times(1).
			Return(nil)

		handler := server.Healthcheck(mockConn)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return status code 500 if postgres is offline", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConn := postgres.NewMockConn(ctrl)
		mockConn.
			EXPECT().
			Ping(ctx).
			Times(1).
			Return(errors.New("error"))

		handler := server.HandlerFunc(server.Healthcheck(mockConn))

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	})
}
