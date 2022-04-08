package app

import (
	"context"
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage"
)

/*type Application interface {
	CreateEvent(ctx context.Context, title string) error
	ChangeEvent(ctx context.Context, id int, title string) error
	RemoveEvent(ctx context.Context, id int) error
	DayEvents(ctx context.Context, date time.Time) []storage.Event
	WeekEvents(ctx context.Context, date time.Time) []storage.Event
	MonthEvents(ctx context.Context, date time.Time) []storage.Event
}*/

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
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

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	_, err := a.storage.Append(event)

	return err
}

func (a *App) ChangeEvent(ctx context.Context, id int, event storage.Event) error {
	return a.storage.Edit(id, event)
}

func (a *App) RemoveEvent(ctx context.Context, id int) error {
	return a.storage.Delete(id)
}

func (a *App) DayEvents(ctx context.Context, date time.Time) []storage.Event {
	return a.storage.DayEvents(date)
}

func (a *App) WeekEvents(ctx context.Context, date time.Time) []storage.Event {
	return a.storage.WeekEvents(date)
}

func (a *App) MonthEvents(ctx context.Context, date time.Time) []storage.Event {
	return a.storage.MonthEvents(date)
}
