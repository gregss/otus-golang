package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage"
)

const success = "true"

type Server struct {
	logger logger.Logger
	app    app.App
	server *http.Server
}

type MyHandler struct {
	app app.App
}

type CreateEditRequest struct {
	Title       string        `json:"title"`
	Time        time.Time     `json:"time"`
	Duration    time.Duration `json:"duration"`
	Description string        `json:"description"`
	UserID      int           `json:"userId"`
	NotifyTime  time.Time     `json:"notifyTime"`
}

type SuccessResponse struct {
	Success string `json:"success"`
}

type RemoveRequest struct {
	EventID int `json:"id"`
}

type TimeEventRequest struct {
	Time time.Time `json:"time"`
}

type TimeEventResponse struct {
	Events interface{} `json:"events"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/create":
		h.createEvent(w, r)
	case "/edit":
		h.editEvent(w, r)
	case "/remove":
		h.removeEvent(w, r)
	case "/day":
		h.dayEvents(w, r)
	case "/week":
		h.weekEvents(w, r)
	case "/month":
		h.monthEvents(w, r)
	default:
		time.Sleep(500 * time.Millisecond) // имитация работы
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode("404 not found")
	}
}

func (h *MyHandler) createEvent(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, r.ContentLength)
	_, _ = r.Body.Read(buf)

	req := &CreateEditRequest{}
	_ = json.Unmarshal(buf, req)

	event := storage.Event{
		Title:       req.Title,
		Time:        req.Time,
		Duration:    req.Duration,
		Description: req.Description,
		UserID:      req.UserID,
		NotifyTime:  req.NotifyTime,
	}

	_ = h.app.CreateEvent(context.Background(), event)

	resp := &SuccessResponse{}
	resp.Success = success
	_ = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
}

func (h *MyHandler) editEvent(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id, _ := strconv.Atoi(q.Get("id"))

	buf := make([]byte, r.ContentLength)
	_, _ = r.Body.Read(buf)

	req := &CreateEditRequest{}
	_ = json.Unmarshal(buf, req)

	event := storage.Event{
		Title:       req.Title,
		Time:        req.Time,
		Duration:    req.Duration,
		Description: req.Description,
		UserID:      req.UserID,
		NotifyTime:  req.NotifyTime,
	}

	_ = h.app.ChangeEvent(context.Background(), id, event)

	resp := &SuccessResponse{}
	resp.Success = "true"
	_ = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
}

func (h *MyHandler) removeEvent(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id, _ := strconv.Atoi(q.Get("id"))

	_ = h.app.RemoveEvent(context.Background(), id)

	resp := &SuccessResponse{}
	resp.Success = "true"
	_ = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
}

func (h *MyHandler) dayEvents(w http.ResponseWriter, r *http.Request) {
	resp := &TimeEventResponse{}
	errresp := &ErrorResponse{}

	if r.Method != http.MethodPost {
		errresp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		WriteResponse(w, resp)
		return
	}
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && errors.Is(err, io.EOF) {
		errresp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, resp)
		return
	}

	req := &TimeEventRequest{}
	err = json.Unmarshal(buf, req)
	if err != nil {
		errresp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, resp)
		return
	}

	resp.Events = h.app.DayEvents(context.Background(), req.Time)
	_ = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
}

func (h *MyHandler) weekEvents(w http.ResponseWriter, r *http.Request) {
	resp := &TimeEventResponse{}
	buf := make([]byte, r.ContentLength)
	_, _ = r.Body.Read(buf)

	req := &TimeEventRequest{}
	_ = json.Unmarshal(buf, req)

	resp.Events = h.app.WeekEvents(context.Background(), req.Time)
	_ = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
}

func (h *MyHandler) monthEvents(w http.ResponseWriter, r *http.Request) {
	resp := &TimeEventResponse{}
	buf := make([]byte, r.ContentLength)
	_, _ = r.Body.Read(buf)

	req := &TimeEventRequest{}
	_ = json.Unmarshal(buf, req)

	resp.Events = h.app.WeekEvents(context.Background(), req.Time)
	_ = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
}

func NewServer(logger logger.Logger, app app.App, port string) *Server {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      loggingMiddleware(&MyHandler{app: app}, logger),
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

func WriteResponse(w http.ResponseWriter, resp interface{}) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
	_, err = w.Write(resBuf)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
