package middleware

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

// ResponseLogger wraps http.ResponseWriter to capture status code
type ResponseLogger struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader captures the status code
func (rl *ResponseLogger) WriteHeader(code int) {
	rl.StatusCode = code
	rl.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			formattedTime := start.Format("02-01-2006 15:50:15")

			// Wrap response writer to capture status code
			rl := &ResponseLogger{ResponseWriter: w, StatusCode: http.StatusOK}

			var requestBody string
			if r.Body != nil {
				bodyBites, _ := io.ReadAll(r.Body)
				requestBody = string(bodyBites)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBites))
			}

			next.ServeHTTP(rl, r)
			duration := time.Since(start)

			logger.Log(
				"time", formattedTime,
				"method", r.Method,
				"path", r.URL.Path,
				"duration", duration,
				"statusCode", rl.StatusCode,
				"requestBody", requestBody,
			)
		})
	}
}

func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
				logger.Log("Error: ", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
