package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

type LogRecord struct {
	http.ResponseWriter
	status int
}

func (r *LogRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			record := &LogRecord{
				ResponseWriter: w,
			}

			next.ServeHTTP(record, r)

			if record.status >= http.StatusBadRequest {
				logger.Sugar().With("status", record.status).Warn("request failed")
			}
		})
	}
}
