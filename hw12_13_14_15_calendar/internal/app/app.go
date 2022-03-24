package app

import (
	"context"
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Storage interface {
	Append(storage.Event) (int, error)
	Edit(id int, newEvent storage.Event) error
	Delete(id int) error
	DayEvents(date time.Time) []storage.Event
	WeekEvents(date time.Time) []storage.Event
	MonthEvents(date time.Time) []storage.Event
}

type DBStorage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
