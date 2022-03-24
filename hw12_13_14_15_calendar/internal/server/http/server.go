package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Server struct {
	logger Logger
	app    Application
	server *http.Server
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

type MyHandler struct{}

type Application interface{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		time.Sleep(500 * time.Millisecond) // имитация работы
		statusCode := http.StatusOK
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode("OK")
	}
}

func NewServer(logger Logger, app Application) *Server {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(&MyHandler{}, logger),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return &Server{logger: logger, app: app, server: server}
}

func (s *Server) Start(ctx context.Context) error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Close()
}
