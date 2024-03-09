package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rodrigosdo/facilities-api/internal/config"
	"github.com/rodrigosdo/facilities-api/internal/cursor"
	"github.com/rodrigosdo/facilities-api/internal/postgres"
	"github.com/rodrigosdo/facilities-api/internal/server/middleware"
	"github.com/rodrigosdo/facilities-api/internal/usecase/worker"

	"cloud.google.com/go/civil"
	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"go.uber.org/zap"
)

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
	logger     *zap.Logger
}

type ErrorResponse struct {
	Message string `json:"error"`
}

type GetAvailableShiftsFromUserRequest struct {
	End   civil.Date `in:"query=end"`
	Start civil.Date `in:"query=start"`
}

type GetAvailableShiftsFromUserReponse struct {
	NextCursor *cursor.Cursor   `json:"next_cursor"`
	Shifts     []AvailableShift `json:"data"`
}

type Facility struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type AvailableShift struct {
	End      time.Time `json:"end"`
	Facility Facility  `json:"facility"`
	ID       int64     `json:"id"`
	Start    time.Time `json:"start"`
}

func New(
	addr string,
	cfg *config.Config,
	logger *zap.Logger,
	conn postgres.Conn,
	workerAvailableShiftsUseCase worker.AvailableShifts,
) (*Server, error) {
	rt := httprouter.New()
	rt.Handler(http.MethodGet, "/healthcheck", NewHandler(Healthcheck(conn), logger))
	rt.Handler(http.MethodGet, "/v1/workers/:id/available_shifts", NewHandler(GetAvailableShiftsFromWorker(workerAvailableShiftsUseCase), logger))

	handler := alice.New(
		middleware.PanicRecovery(logger),
		middleware.RequestLogger(logger),
		gziphandler.GzipHandler,
	).Then(rt)

	srv := &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		},
		logger: logger,
	}

	return srv, nil
}

func (s Server) Run(ctx context.Context) error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			s.logger.Error("server: fail to listen and server", zap.Error(err))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Block until we receive an exit signal
	<-signalChan

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Server.GracefulShutdownPeriod)
	defer cancel()

	return s.Shutdown(ctx)
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
