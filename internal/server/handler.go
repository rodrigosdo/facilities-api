package server

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/cursor"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/postgres"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/usecase/worker"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Error struct {
	Code    int    `json:"code,omitempty"`
	Error   error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type Handler struct {
	handlerFunc HandlerFunc
	logger      *zap.Logger
}

func NewHandler(handler HandlerFunc, logger *zap.Logger) Handler {
	return Handler{
		handlerFunc: handler,
		logger:      logger,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.handlerFunc(w, r); err != nil {
		httpRequest, _ := httputil.DumpRequest(r, false)

		h.logger.Warn(err.Message,
			zap.Error(err.Error),
			zap.String("request", string(httpRequest)),
			zap.Time("time", time.Now()),
			zap.Int("status_code", err.Code),
		)

		type errorResponse struct {
			Error Error `json:"error"`
		}

		w.WriteHeader(err.Code)

		if encodeErr := json.NewEncoder(w).Encode(errorResponse{Error: Error{Message: err.Message}}); encodeErr != nil {
			http.Error(w, err.Message, err.Code)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

type HandlerFunc func(http.ResponseWriter, *http.Request) *Error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		type errorResponse struct {
			Error Error `json:"error"`
		}

		w.WriteHeader(err.Code)

		if encodeErr := json.NewEncoder(w).Encode(errorResponse{Error: Error{Message: err.Message}}); encodeErr != nil {
			http.Error(w, err.Message, err.Code)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

func GetAvailableShiftsFromWorker(uc worker.AvailableShifts) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		queryCursor, err := cursor.Parse(r.URL.Query().Get("cursor"))
		if err != nil {
			return &Error{
				Code:    http.StatusBadRequest,
				Error:   err,
				Message: "cursor query param is invalid",
			}
		}

		start, err := parseDate(r.URL.Query().Get("start"))
		if err != nil {
			return &Error{
				Code:    http.StatusBadRequest,
				Error:   err,
				Message: "start date query param is invalid",
			}
		}

		end, err := parseDate(r.URL.Query().Get("end"))
		if err != nil {
			return &Error{
				Code:    http.StatusBadRequest,
				Error:   err,
				Message: "end date query param is invalid",
			}
		}

		limit, err := parseLimit(r.URL.Query().Get("limit"))
		if err != nil {
			return &Error{
				Code:    http.StatusBadRequest,
				Error:   err,
				Message: "limit query param is invalid",
			}
		}

		workerID, err := parseWorkerID(httprouter.ParamsFromContext(r.Context()).ByName("id"))
		if err != nil {
			return &Error{
				Code:    http.StatusBadRequest,
				Error:   err,
				Message: "worker id is invalid",
			}
		}

		availableShifts, err := uc.GetAvailableShifts(r.Context(), worker.GetAvailableShiftsRequest{
			Cursor:   queryCursor,
			End:      *end,
			Limit:    limit,
			Start:    *start,
			WorkerID: workerID,
		})
		if err != nil {
			return &Error{
				Code:    http.StatusInternalServerError,
				Error:   err,
				Message: "failed to get available shifts",
			}
		}

		w.Header().Set("Content-Type", "application/json")

		resp := GetAvailableShiftsFromUserReponse{
			NextCursor: availableShifts.NextCursor,
		}

		for _, s := range availableShifts.Shifts {
			resp.Shifts = append(resp.Shifts, AvailableShift{
				End: s.End,
				Facility: Facility{
					ID:   s.Facility.ID,
					Name: s.Facility.Name,
				},
				ID:    s.ID,
				Start: s.Start,
			})
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return &Error{
				Code:    http.StatusInternalServerError,
				Error:   err,
				Message: "failed to encode response",
			}
		}

		return nil
	}
}

func parseDate(dateStr string) (*civil.Date, error) {
	var date civil.Date

	if dateStr != "" {
		parsedDate, err := civil.ParseDate(dateStr)
		if err != nil {
			return nil, err
		}

		date = parsedDate
	}

	return &date, nil
}

func parseLimit(limitStr string) (int, error) {
	var limit int

	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, err
		}

		if limitInt > 0 {
			limit = limitInt
		}
	}

	return limit, nil
}

func parseWorkerID(workerIDStr string) (int64, error) {
	workerID, err := strconv.ParseInt(workerIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return workerID, nil
}

func Healthcheck(conn postgres.Conn) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		if err := conn.Ping(r.Context()); err != nil {
			return &Error{
				Code:    http.StatusServiceUnavailable,
				Error:   err,
				Message: "failed to query postgres",
			}
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}
