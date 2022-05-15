package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(next http.Handler, logger logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logText := fmt.Sprintf("%s [%s] %s %s %s %d %s \"%s\"",
			r.RemoteAddr,
			time.Now().String(),
			r.Method,
			r.RequestURI,
			r.Proto,
			200,
			time.Since(start),
			r.UserAgent(),
		)
		logger.Info(logText)
	})
}
